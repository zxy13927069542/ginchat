package service

import (
	"ginchat/models"
)

type JSONResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JSONBResult struct {
	Code    int
	Message string `json:"msg"`
	Data    interface{}
}

type ListResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Rows    interface{}
	Total   int
}

type PhoneS struct {
	Phone string `json:"phone" valid:"mobile"`
}

type EmailS struct {
	Email string `json:"email" valid:"email"`
}

type UserRegisterReq struct {
	UserName   string `form:"name" json:"userName" binding:"required"`     //	用户名
	Password   string `form:"password" json:"password" binding:"required"` //	密码
	RePassword string `form:"Identity" json:"Identity" binding:"required"` //	密码确认
}

type UserRegisterResp struct {
	models.UserBasic
}

type UserDelReq struct {
	ID uint `json:"id" binding:"required"`
}

type UserGetReq struct {
	ID uint `json:"id" binding:"required"`
}

type UserUpdateReq struct {
	ID       uint   `json:"id" binding:"required"`
	Name     string `json:"name"`
	PassWord string `json:"password"`
	PhoneS
	EmailS
}

type UserLoginReq struct {
	Name     string `json:"name" binding:"required"`
	PassWord string `json:"password" binding:"required"`
}

type UserLoginResp struct {
	Token        string `json:"token"`
	AccessExpire string `json:"accessExpire"`
	RefreshAfter string `json:"refreshAfter"`
}

type FindUserReq struct {
	Name     string `form:"name" binding:"required"`
	PassWord string `form:"password" binding:"required"`
}

type FindUserResp struct {
	models.UserBasic
}

type ToChatReq struct {
	UserId uint   `form:"userId"`
	Token  string `form:"token"`
}

type SearchFriendReq struct {
	UserId uint `form:"userId"`
}

type SearchFriendResp struct {
	JSONResult
	Rows  []models.UserBasic
	Total int
}

type ChatReq struct {
	UserId int64  `json:"userId" form:"userId"`
	Token  string `json:"token" form:"token"`
}

// UploadReq 文件上传
type UploadReq struct {
	UserId   int64  `json:"userid" form:"userid"`
	FileType string `json:"filetype" form:"filetype"`
}

// AddFriendReq 添加好友
type AddFriendReq struct {
	UserId     uint   `form:"userId"`
	TargetName string `form:"targetName"`
}

// CreateGroupReq 创建群聊
type CreateGroupReq struct {
	OwnerId string `form:"ownerId"`
	Name    string `form:"name"`
	Icon    string `form:"icon"`
	Desc    string `form:"desc"`
}

// JoinGroupReq 加入群聊
type JoinGroupReq struct {
	UserId string `form:"userId"`
	Group  string `form:"comId"`
}

// LoadGroupReq 加载群聊列表
type LoadGroupReq struct {
	OwnerId string `form:"ownerId"`
}

//	RedisMsgReq 读取历史消息
type RedisMsgReq struct {
	UserIdA string `form:"userIdA"`
	UserIdB string `form:"userIdB"`
	Start   string `form:"start"`
	End     string `form:"end"`
	IsRev   string `form:"isRev"`
}
