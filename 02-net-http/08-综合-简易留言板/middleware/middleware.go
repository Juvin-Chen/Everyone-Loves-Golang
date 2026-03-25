/*
实现三个中间件：
	Logger：记录请求方法和路径，以及耗时。
	Recovery：捕获 panic，返回 500 并记录日志。
	RequestID：为每个请求生成唯一 ID 并放入 Context，供后续使用。
*/

package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

// Logger 中间件：记录请求日志
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

// Recovery 中间件：捕获 panic
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RequestID 中间件：生成请求 ID 并放入 context
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 简单生成唯一 ID（实际可用 uuid）
		id := fmt.Sprintf("%d", time.Now().UnixNano())
		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
