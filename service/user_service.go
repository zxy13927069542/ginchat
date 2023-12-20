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
	model *models.UserBasicModel
}

func NewUserBasicService() *UserBasicService {
	return &UserBasicService{
		model: models.NewUserBasicModel(),
	}
}

// Register
//
//	@Tags		用户模块
//	@Summary	注册用户
//	@Param		account	body		service.UserRegisterReq	true	"账号"
//	@Success	200		{object}	service.JSONResult{}
//	@Failure	400		{object}	service.JSONResult{}
//	@Router		/user/register [post]
func (s *UserBasicService) Register(c *gin.Context) {
	var user UserRegisterReq
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, "参数错误", nil})
		return
	}

	if user.Password != user.RePassword {
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, "两次密码不一致", nil})
		return
	}

	if nameByUser := s.model.FindUserByName(user.UserName); nameByUser.Name != "" {
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, "用户已存在", nil})
		return
	}


	var m models.UserBasic
	m.Name = user.UserName
	m.Salt = utils.GenSalt()
	m.PassWord = utils.MakePassword(user.Password, m.Salt)
	if err := s.model.Create(m).Error; err != nil {
		log.Errorf(">>Create user failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusInternalServerError, JSONResult{500, "内部错误", nil})
		return
	}
	c.IndentedJSON(http.StatusOK, JSONResult{200, "创建成功", nil})
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
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, "参数错误", nil})
		return
	}

	var m models.UserBasic
	m.ID = user.ID
	s.model.Delete(m)
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
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, "参数错误", nil})
		return
	}

	var m models.UserBasic
	m.ID = user.ID
	if err := s.model.Get(&m).Error; err == gorm.ErrRecordNotFound {
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, "查无此用户", nil})
		return
	} else if err != nil {
		log.Errorf(">>Get user failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusInternalServerError, JSONResult{500, "内部错误", nil})
		return
	}
	c.IndentedJSON(200, JSONResult{200, "成功", m})
}

// Update
//
//	@Tags		用户模块
//	@Summary	更新用户
//	@Param		account	body		service.UserUpdateReq	true "用户信息"
//	@Success	200		{object}	service.JSONResult{}
//	@Failure	400		{object}	service.JSONResult{}
//	@Router		/user/update [post]
func (s *UserBasicService) Update(c *gin.Context) {
	var user UserUpdateReq
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, "参数错误", nil})
		return
	}

	if _, err := govalidator.ValidateStruct(user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, err.Error(), nil})
		return
	}

	var m models.UserBasic
	m.ID = user.ID
	m.Name = user.Name
	m.PassWord = user.PassWord
	m.Phone = user.Phone
	m.Email = user.Email
	if err := s.model.Update(m).Error; err != nil {
		log.Errorf(">>Get user failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusInternalServerError, JSONResult{500, "内部错误", nil})
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
	list, _ := s.model.List()
	c.IndentedJSON(200, JSONResult{200, "成功", list})
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
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, "参数错误", nil})
		return
	}

	userByName := s.model.FindUserByName(req.Name)
	if userByName.Name == "" {
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, "用户不存在", nil})
		return
	}

	if !utils.ValidatePassword(req.PassWord, userByName.Salt, userByName.PassWord) {
		c.IndentedJSON(http.StatusBadRequest, JSONResult{400, "密码错误", nil})
		return
	}
	userByName.UpdatedAt = time.Now()
	userByName.LoginTime.Scan(time.Now())
	userByName.IsLogOut = false
	s.model.Update(userByName)

	token := utils.GenToken(userByName.Name, userByName.PassWord)
	c.IndentedJSON(http.StatusOK, JSONResult{200, "登陆成功", *token})
}

//	upgrader 跨域检查
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *UserBasicService) SendMsg(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf(">>SendMsg() upgrade failed! Err: [%v]", err)
		c.IndentedJSON(http.StatusBadRequest, JSONResult{500, "内部错误", nil})
		return
	}
	defer conn.Close()

	msgHandler(conn, c)
}

func msgHandler(ws *websocket.Conn, c *gin.Context) {
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


