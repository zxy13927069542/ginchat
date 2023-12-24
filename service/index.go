package service

import (
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"text/template"
)

// Index
//
//	@Tags		首页
//	@Success	200	{string}	pong
//	@Router		/ [get]
func Index(c *gin.Context) {
	temp := template.Must(template.ParseFiles("index.html", "views/chat/head.html"))
	temp.Execute(c.Writer, "index")
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
	// c.JSON(200, gin.H{
	// 	"message": "welcome !!  ",
	// })
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}
	var req ToChatReq
	if err := c.Bind(&req); err != nil {
		c.IndentedJSON(http.StatusOK, JSONResult{400, "参数错误", nil})
		return
	}
	user := models.UserBasic{}
	user.ID = req.UserId
	user.Identity = req.Token
	//fmt.Println("ToChat>>>>>>>>", user)
	ind.Execute(c.Writer, user)
	// c.JSON(200, gin.H{
	// 	"message": "welcome !!  ",
	// })
}
