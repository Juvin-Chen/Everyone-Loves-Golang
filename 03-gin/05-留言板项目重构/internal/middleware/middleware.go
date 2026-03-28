package middleware

import "github.com/gin-gonic/gin"

// Logger 统一使用 Gin 内置日志中间件。
func Logger() gin.HandlerFunc {
	return gin.Logger()
}

// Auth 预留认证中间件入口。
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
