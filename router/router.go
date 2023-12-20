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

	userService := service.NewUserBasicService()
	chatService := service.NewChatService()

	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/index", service.GetIndex)
	r.POST("/user/register", userService.Register)
	r.POST("/user/login", userService.Login)
	r.POST("/user/delete", userService.Delete)
	r.POST("/user/get", userService.Get)
	r.POST("/user/update", userService.Update)
	r.GET("/user/list", userService.List)
	r.GET("/user/SendMsg", userService.SendMsg)
	r.GET("/user/chat", chatService.Chat)
	r.Static("/user/chat/", "html/")

	return r
}
