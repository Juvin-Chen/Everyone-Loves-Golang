package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的 Gin 引擎（自带 Logger 和 Recovery 中间件）
	r := gin.Default()

	// 定义一个 GET 路由，路径为 /ping
	r.GET("/ping", func(c *gin.Context) {
		// 返回json数据
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 启动服务
	r.Run()
}
