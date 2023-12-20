package service

type JSONResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PhoneS struct {
	Phone string `json:"phone" valid:"mobile"`
}

type EmailS struct {
	Email string `json:"email" valid:"email"`
}

type UserRegisterReq struct {
	UserName   string `json:"userName" binding:"required"`   //	用户名
	Password   string `json:"password" binding:"required"`   //	密码
	RePassword string `json:"rePassword" binding:"required"` //	密码确认
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
