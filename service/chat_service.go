package service

import (
	"encoding/json"
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/pion/ion-log"
	"gopkg.in/fatih/set.v0"
	"net"
	"net/http"
	"sync"
)

type ChatService struct {
	//model *models.UserBasicModel
}

func NewChatService() *ChatService {
	return &ChatService{
		//model: models.NewUserBasicModel(),
	}
}

type Node struct {
	Conn *websocket.Conn
	DataQueue chan []byte
	GroupSet set.Interface
}

//	clientMap 映射关系
var clientMap map[uint]*Node = make(map[uint]*Node, 0)
//	rwLock 读写锁
var rwLock sync.RWMutex

func (s *ChatService) Chat(c *gin.Context) {
	//query := r.URL.Query()
	//msgType := query.Get("type")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	//userId, _ := strconv.Atoi(query.Get("userId"))
	userId := c.GetUint("userId")

	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf(">>Chat() failed! Err: [%v]", err)
		return
	}

	//	获取conn
	node := &Node{
		Conn: ws,
		DataQueue: make(chan []byte, 50),
		GroupSet: set.New(set.ThreadSafe),
	}
	rwLock.Lock()
	clientMap[userId] = node
	rwLock.Unlock()

	go sendProc(node)
	go rcvProc(node)
	sendPrivateMsg(userId, "欢迎来到聊天系统!")
}

//	sendProc 读取消息队列里的数据并发送到websocket
func sendProc(node *Node) {
	for data := range node.DataQueue {
		if err := node.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Errorf(">>sendProc() Send msg to websocket failed! Err: [%v]", err)
			return
		}
	}
}

//	rcvProc 读取websocket消息并处理
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
		log.Infof("[ws]: %v", msg)
		//	消息处理,广播,群发,或私信
		wsMsgHandler(msg)
	}
}

//	udpSendChan udp消息管道 里面的消息会发送到udp
var udpSendChan chan []byte = make(chan []byte, 1024)

//	wsMsgHandler websocket消息处理器 收到ws消息后往udp消息管道发消息
func wsMsgHandler(msg []byte) {
	udpSendChan <- msg
}

func init() {
	go udpSendProc()
	go udpRcvProc()
}

//	udpRcvProc udp数据接收协程,收到消息后进行处理
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
		udpMessageHandler(buff[0:n])
	}
}

//	udpMessageHandler udp消息处理器 收到udp消息后根据消息类型进行处理
func udpMessageHandler(data []byte) {
	var msg models.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Errorf(">>udpMessageHandler() parse data failed! Err: [%v]", err)
		return
	}

	switch msg.Type {
	case "1":	//	私信
		sendPrivateMsg(msg.TargetID, msg.Contend)

	}
}

//	sendPrivateMsg 私信
func sendPrivateMsg(id uint, contend string) {
	rwLock.RLock()
	node, ok := clientMap[id]
	rwLock.RUnlock()
	if ok {
		node.DataQueue <- []byte(contend)
	}
}

//	udpSendProc udp数据发送协程
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
