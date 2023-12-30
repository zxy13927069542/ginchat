package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/redisc"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"
	log "github.com/pion/ion-log"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type UserBasicService struct {
	um *models.UserBasicModel
	cm *models.ContactModel
}

func NewUserBasicService() *UserBasicService {
	return &UserBasicService{
		um: models.UserBasicM,
		cm: models.ContactM,
	}
}

// Register
//
//	@Tags		用户模块
//	@Summary	注册用户
//	@Param		name		query		string	true	"账号"
//	@Param		password	query		string	true	"密码"
//	@Param		Identity	query		string	true	"密码确认"
//	@Success	200			{object}	service.JSONResult{data=service.UserRegisterResp}
//	@Router		/user/register [get]
func (s *UserBasicService) Register(c *gin.Context) {
	var user UserRegisterReq
	if err := c.Bind(&user); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "参数错误", nil})
		return
	}

	if user.Password != user.RePassword {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "两次密码不一致", nil})
		return
	}

	if nameByUser := s.um.FindUserByName(user.UserName); nameByUser.Name != "" {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "用户已存在", nil})
		return
	}

	var m models.UserBasic
	m.Name = user.UserName
	m.Salt = utils.GenSalt()
	m.PassWord = utils.MakePassword(user.Password, m.Salt)
	if err := s.um.Create(m).Error; err != nil {
		log.Errorf(">>Create user failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusOK, JSONResult{500, "内部错误", nil})
		return
	}
	c.IndentedJSON(http.StatusOK, JSONResult{200, "创建成功", m})
}

// Delete
//
//	@Tags		用户模块
//	@Summary	删除用户
//	@Param		account	body		service.UserDelReq	true	"ID"
//	@Success	200		{object}	service.JSONResult{}
//	@Failure	400		{object}	service.JSONResult{}
//	@Router		/user/delete [post]
func (s *UserBasicService) Delete(c *gin.Context) {
	var user UserDelReq
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "参数错误", nil})
		return
	}

	var m models.UserBasic
	m.ID = user.ID
	s.um.Delete(m)
	c.IndentedJSON(http.StatusOK, JSONResult{200, "删除成功", nil})
}

// Get
//
//	@Tags		用户模块
//	@Summary	查询用户
//	@Param		account	body		service.UserGetReq	true	"ID"
//	@Success	200		{object}	service.JSONResult{data=models.UserBasic}
//	@Failure	400		{object}	service.JSONResult{}
//	@Router		/user/get [post]
func (s *UserBasicService) Get(c *gin.Context) {
	var user UserGetReq
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "参数错误", nil})
		return
	}

	var m models.UserBasic
	m.ID = user.ID
	if err := s.um.Get(&m).Error; err == gorm.ErrRecordNotFound {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "查无此用户", nil})
		return
	} else if err != nil {
		log.Errorf(">>Get user failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusOK, JSONResult{500, "内部错误", nil})
		return
	}
	c.IndentedJSON(200, JSONResult{200, "成功", m})
}

// Update
//
//	@Tags		用户模块
//	@Summary	更新用户
//	@Param		account	body		service.UserUpdateReq	true	"用户信息"
//	@Success	200		{object}	service.JSONResult{}
//	@Failure	400		{object}	service.JSONResult{}
//	@Router		/user/update [post]
func (s *UserBasicService) Update(c *gin.Context) {
	var user UserUpdateReq
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "参数错误", nil})
		return
	}

	if _, err := govalidator.ValidateStruct(user); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, err.Error(), nil})
		return
	}

	var m models.UserBasic
	m.ID = user.ID
	m.Name = user.Name
	m.PassWord = user.PassWord
	m.Phone = user.Phone
	m.Email = user.Email
	if err := s.um.Update(m).Error; err != nil {
		log.Errorf(">>Get user failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusOK, JSONResult{500, "内部错误", nil})
		return
	}
	c.IndentedJSON(200, JSONResult{200, "成功", nil})
}

// List
//
//	@Tags		用户模块
//	@Summary	用户列表
//	@Success	200	{object}	service.JSONResult{data=[]models.UserBasic}
//	@Failure	400	{object}	service.JSONResult{}
//	@Router		/user/list [get]
func (s *UserBasicService) List(c *gin.Context) {
	list, _ := s.um.List()
	c.IndentedJSON(http.StatusOK, JSONResult{200, "成功", list})
}

// Login
//
//	@Tags		用户模块
//	@Summary	用户登陆
//	@Param		account	body		service.UserLoginReq	true	"账号"
//	@Success	200		{object}	service.JSONResult{data=utils.TokenResp}
//	@Failure	400		{object}	service.JSONResult{}
//	@Router		/user/login [post]
func (s *UserBasicService) Login(c *gin.Context) {
	var req UserLoginReq
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "参数错误", nil})
		return
	}

	userByName := s.um.FindUserByName(req.Name)
	if userByName.Name == "" {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "用户不存在", nil})
		return
	}

	if !utils.ValidatePassword(req.PassWord, userByName.Salt, userByName.PassWord) {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "密码错误", nil})
		return
	}
	userByName.UpdatedAt = time.Now()
	userByName.LoginTime.Scan(time.Now())
	userByName.IsLogOut = false
	s.um.Update(userByName)

	token := utils.GenToken(userByName.Name, userByName.PassWord)
	c.IndentedJSON(http.StatusOK, JSONResult{200, "登陆成功", *token})
}

// FindUserByNameAndPwd
//
//	@Tags		用户模块
//	@Summary	用户登陆
//	@Param		name		query		string	true	"账号"
//	@Param		password	query		string	true	"密码"
//	@Success	200			{object}	service.JSONResult{data=models.UserBasic}
//	@Failure	400			{object}	service.JSONResult{}
//	@Router		/user/findUserByNameAndPwd [post]
func (s *UserBasicService) FindUserByNameAndPwd(c *gin.Context) {
	var req FindUserReq
	if err := c.Bind(&req); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "参数错误", nil})
		return
	}

	userByName := s.um.FindUserByName(req.Name)
	if userByName.Name == "" {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "用户不存在", nil})
		return
	}

	if !utils.ValidatePassword(req.PassWord, userByName.Salt, userByName.PassWord) {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "密码错误", nil})
		return
	}
	userByName.UpdatedAt = time.Now()
	userByName.LoginTime.Scan(time.Now())
	userByName.IsLogOut = false
	s.um.Update(userByName)

	userByName.Identity = utils.GenToken(userByName.Name, userByName.PassWord).Token
	c.IndentedJSON(http.StatusOK, JSONResult{200, "登陆成功", FindUserResp{userByName}})
}

// upgrader 跨域检查
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *UserBasicService) SendMsg(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf(">>SendMsg() upgrade failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusOK, JSONResult{500, "内部错误", nil})
		return
	}
	defer conn.Close()

	msgHandler1(conn, c)
}

func (s *UserBasicService) SearchFriends(c *gin.Context) {
	var req SearchFriendReq
	if err := c.Bind(&req); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "参数错误", nil})
		return
	}

	friends, err := s.cm.SearchFriend(req.UserId)
	if err != nil {
		log.Errorf(">>SearchFriends() failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusOK, JSONResult{500, "内部错误", nil})
		return
	}

	c.IndentedJSON(http.StatusOK, ListResp{200, "查询成功", nil, friends, len(friends)})
}

func msgHandler1(ws *websocket.Conn, c *gin.Context) {
	//	监听redis频道消息
	sub := redisc.Subscribe(c, redisc.SubscribeKey)
	defer sub.Close()

	go func() {
		channel := sub.Channel()
		for msg := range channel {
			//  收到redis消息后发送websocket消息
			msgFormat := fmt.Sprintf("[ws][%s]:%s", time.Now().Format("2006-01-02 15:04:05"), msg.Payload)
			if err := ws.WriteMessage(websocket.TextMessage, []byte(msgFormat)); err != nil {
				log.Errorf(">>msgHandler() Sent msg to websocket failed! Err: [%v]", err)
				return
			}
		}
	}()

	//  读取websocket消息并往redis发送消息
	for {
		//  读取消息
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Errorf(">>msgHandler() Read msg from websocket failed! Err: [%v]", err)
			return
		}

		fmt.Printf("收到用户消息: %s\n", string(msg))
		if err = redisc.Publish(c, redisc.PublishKey, string(msg)); err != nil {
			log.Errorf(">>msgHandler() publish msg failed! Err: [%v]", err)
			return
		}

	}
}

//	AddFriend 添加好友
func (s *UserBasicService) AddFriend(c *gin.Context) {
	var req AddFriendReq
	if err := c.Bind(&req); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "参数错误", nil})
		return
	}

	if err := s.cm.AddFriend(req.UserId, req.TargetName); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{500, err.Error(), nil})
		return
	}
	c.IndentedJSON(http.StatusOK, JSONResult{200, "成功", nil})
}
