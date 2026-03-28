package router

import (
	"message-board-gin/internal/handler"
	"message-board-gin/internal/middleware"
	"message-board-gin/internal/repository"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()

	// 静态文件
	r.Static("/static", "./static")

	// 加载模板
	r.LoadHTMLGlob("templates/*")

	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// 依赖注入
	repo := repository.NewMessageRepository()
	msgHandler := handler.NewMessageHandler(repo)

	// 路由
	r.GET("/", msgHandler.Index)
	r.GET("/new", msgHandler.NewForm)
	r.POST("/messages", msgHandler.Create)
	r.POST("/delete", msgHandler.Delete)

	return r
}
