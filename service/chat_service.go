package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ginchat/models"
	"ginchat/redisc"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	log "github.com/pion/ion-log"
	"gopkg.in/fatih/set.v0"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

type ChatService struct {
	//um *models.UserBasicModel
}

func NewChatService() *ChatService {
	return &ChatService{
		//um: models.NewUserBasicModel(),
	}
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSet  set.Interface
}

const (
	PrivateMsg = 1
	GroupMsg = 2
)

// clientMap 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// rwLock 读写锁
var rwLock sync.RWMutex

func (s *ChatService) Chat(c *gin.Context) {
	//query := r.URL.Query()
	//msgType := query.Get("type")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	//userId, _ := strconv.Atoi(query.Get("userId"))
	var req ChatReq
	if err := c.Bind(&req); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "参数错误", nil})
		return
	}

	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf(">>Chat() failed! Err: [%v]", err)
		return
	}

	//	获取conn
	node := &Node{
		Conn:      ws,
		DataQueue: make(chan []byte, 50),
		GroupSet:  set.New(set.ThreadSafe),
	}
	rwLock.Lock()
	clientMap[req.UserId] = node
	rwLock.Unlock()

	go sendProc(node)
	go rcvProc(node)
	//sendPrivateMsg(userId, "欢迎来到聊天系统!")
}

// sendProc 读取消息队列里的数据并发送到websocket
func sendProc(node *Node) {
	for data := range node.DataQueue {
		//	收到私信
		log.Infof(">>收到私信: %s", string(data))
		if err := node.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Errorf(">>sendProc() Send msg to websocket failed! Err: [%v]", err)
			return
		}
	}
}

// rcvProc 读取websocket消息并处理
func rcvProc(node *Node) {
	//  读取websocket消息并处理
	for {
		//  读取消息
		_, msg, err := node.Conn.ReadMessage()
		if err != nil {
			log.Errorf(">>rcvProc() Read msg from websocket failed! Err: [%v]", err)
			return
		}

		//	日志打印收到的消息
		//log.Infof("[ws]: %v", msg)
		//	消息处理,广播,群发,或私信
		msgHandler(msg)
	}
}

// msgHandler 收到ws消息或udp消息后根据消息进行处理
func msgHandler(data []byte) {
	var msg models.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Errorf(">>wsMsgHandler() parse data failed! Err: [%v]", err)
		return
	}

	log.Infof("[ws]: %s", string(data))
	switch msg.Type {
	case PrivateMsg: //	私信
		sendPrivateMsg(msg.TargetID, data, msg)
	case GroupMsg:	//	群聊
		sendGroupMsg(uint(msg.FromID), uint(msg.TargetID), data)

	}
}

// sendPrivateMsg 私信
func sendPrivateMsg(targetId int64, data []byte, msg models.Message) {
	//	用户在线则往websocket里发
	rwLock.RLock()
	node, ok := clientMap[targetId]
	rwLock.RUnlock()
	if ok {
		node.DataQueue <- data
	}

	//	把消息通过sorted set 缓存到redis中
	//	zadd key [msg_大ID_小ID] score [消息创建时间] value [消息本体]
	key := ""
	c := context.Background()
	if msg.FromID > msg.TargetID {
		key = fmt.Sprintf("msg_%d_%d", msg.FromID, msg.TargetID)
	} else {
		key = fmt.Sprintf("msg_%d_%d", msg.TargetID, msg.FromID)
	}
	redisc.Redisc.ZAdd(c, key, &redis.Z{Score: float64(msg.CreateTime), Member: data})
}

//	sendGroupMsg 群聊
func sendGroupMsg(userId uint, groupId uint, data []byte) {
	ids := models.NewContactModel().SearchGroupFriends(groupId)
	for _, id := range ids {
		//	跳过自己
		if id == userId {
			continue
		}

		rwLock.RLock()
		if v, exist := clientMap[int64(id)]; exist {
			v.DataQueue <- data
		}
		rwLock.RUnlock()
	}
	//	todo: 群聊消息也通过redis缓存
}

// udpSendChan udp消息管道 里面的消息会发送到udp 用于群发
var udpSendChan chan []byte = make(chan []byte, 1024)

func init() {
	go udpSendProc()
	go udpRcvProc()
}

// udpRcvProc udp数据接收协程,收到消息后进行处理
func udpRcvProc() {
	udp, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		log.Errorf(">>udpRcvProc() listen failed! Err: [%v]", err)
		return
	}
	defer udp.Close()
	for {
		var buff [512]byte
		n, err := udp.Read(buff[0:])
		if err != nil {
			log.Errorf(">>udpRcvProc() read message failed! Err: [%v]", err)
			return
		}
		msgHandler(buff[0:n])
	}
}

// udpSendProc udp数据发送协程
func udpSendProc() {
	udp, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 3, 16),
		Port: 3000,
	})
	if err != nil {
		log.Errorf(">>udpSendProc() Dial failed! Err: [%v]", err)
		return
	}
	defer udp.Close()

	for data := range udpSendChan {
		if _, err := udp.Write(data); err != nil {
			log.Errorf(">>udpSendProc() send message failed! Err: [%v]", err)
			return
		}

	}
}

// Upload 上传聊天图片
func (s *ChatService) Upload(c *gin.Context) {
	var req UploadReq
	_ = c.ShouldBindWith(&req, binding.FormMultipart)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Errorf(">>ChatService.Upload() failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusOK, JSONBResult{500, "内部错误", nil})
		return
	}

	content, err := io.ReadAll(file)
	defer file.Close()
	if err != nil {
		log.Errorf(">>ChatService.Upload() failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusOK, JSONBResult{500, "内部错误", nil})
		return
	}
	//	根据文件内容生成校验和,确保文件唯一
	name := utils.Sum(content)
	splitArr := strings.Split(header.Filename, ".")
	if len(splitArr) > 1 {
		name = fmt.Sprintf("%s.%s", name, splitArr[len(splitArr)-1])
	} else {
		//	默认格式为png
		name = fmt.Sprintf("%s%s", name, req.FileType)
	}
	url := fmt.Sprintf("./asset/upload/%s", name)
	dstFile, err := os.Create(url)
	defer dstFile.Close()
	if err != nil {
		log.Errorf(">>ChatService.Upload() failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusOK, JSONBResult{500, "内部错误", nil})
		return
	}
	if _, err = dstFile.Write(content); err != nil {
		log.Errorf(">>ChatService.Upload() failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusOK, JSONBResult{500, "内部错误", nil})
		return
	}

	c.IndentedJSON(http.StatusOK, JSONBResult{200, "上传成功", url})
}

//	LoadMsg 加载历史消息
func (s *ChatService) LoadMsg(c *gin.Context) {
	var req RedisMsgReq
	if err := c.Bind(&req); err != nil {
		c.IndentedJSON(http.StatusOK, JSONBResult{400, "参数错误", nil})
		return
	}

	msgs := loadMsg(utils.Str2Int64(req.UserIdA), utils.Str2Int64(req.UserIdB),
		utils.Str2Bool(req.IsRev), utils.Str2Int64(req.Start), utils.Str2Int64(req.End), c)
	c.IndentedJSON(http.StatusOK, ListResp{Code: 200, Message: "成功", Rows: msgs, Total: len(msgs)})
}

func loadMsg(userA, userB int64, isRev bool, start, end int64, c context.Context) []string {
	key := ""
	if userA > userB {
		key = fmt.Sprintf("msg_%d_%d", userA, userB)
	} else {
		key = fmt.Sprintf("msg_%d_%d", userB, userA)
	}

	switch isRev {
	case true:
		result, _ := redisc.Redisc.ZRange(c, key, start, end).Result()
		return result
	case false:
		result, _ := redisc.Redisc.ZRevRange(c, key, start, end).Result()
		return result
	}
	return nil
}
