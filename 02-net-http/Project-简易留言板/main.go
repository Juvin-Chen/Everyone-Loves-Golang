/*
1.初始化存储
2.注册路由
3.应用中间件
4.静态文件服务
5.启动服务器
*/

package main

import (
	"log"
	"message-board/handlers"
	"message-board/middleware"
	"message-board/store"
	"net/http"
)

func main() {
	// 初始化内存存储
	memStore := store.NewMemoryStore()
	// 创建路由器
	mux := http.NewServeMux()

	// 注册路由
	mux.HandleFunc("/", handlers.IndexHandler(memStore))
	mux.HandleFunc("/new", handlers.NewMessageFormHandler)
	mux.HandleFunc("/delete", handlers.DeleteMessageHandler(memStore))
	mux.HandleFunc("/messages", handlers.CreateMessageHandler(memStore))

	// 静态文件服务
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// 应用中间件（顺序：恢复 → 日志 → 请求ID → 路由）
	var handler http.Handler = mux
	handler = middleware.Recovery(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RequestID(handler)

	// 启动服务器
	log.Println("服务器启动在 http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
