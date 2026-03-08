package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("开始处理请求：%s %s", r.Method, r.URL.Path)

		// 执行下一个处理器
		next(w, r)

		log.Printf("请求处理完成，耗时：%v", time.Since(start))
	}
}
