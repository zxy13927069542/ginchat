package router

import (
	"ginchat/docs"
	"ginchat/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	//	服务
	userService := service.NewUserBasicService()
	chatService := service.NewChatService()

	//	swagger api文档配置
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//静态资源
	r.Static("/asset", "asset/")
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	//	r.StaticFS()
	r.LoadHTMLGlob("views/**/*")

	//首页
	r.GET("/", service.Index)
	r.GET("/index", service.Index)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", chatService.Chat)
	r.POST("/searchFriends", userService.SearchFriends)

	//	用户模块
	r.POST("/user/createUser", userService.Register)
	r.POST("/user/login", userService.Login)
	r.POST("/user/findUserByNameAndPwd", userService.FindUserByNameAndPwd)
	r.POST("/user/delete", userService.Delete)
	r.POST("/user/get", userService.Get)
	r.POST("/user/update", userService.Update)
	r.GET("/user/list", userService.List)
	r.GET("/user/SendMsg", userService.SendMsg)

	return r
}
