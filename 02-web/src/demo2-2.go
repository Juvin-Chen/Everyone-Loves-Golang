// Handle 与 HandleFunc 结合使用

package main

import "net/http"

// 自定义 Handler 类型
type HelloHandler struct{}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Handle!"))
}

// 普通函数类型
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home fron HandleFunc!"))
}

func demo2_2() {
	// 使用 Handle 注册自定义 Handler
	http.Handle("/hello", &HelloHandler{})

	// 使用 HandleFunc 注册普通函数，会内部进行转换成HandlerFunc类型
	http.HandleFunc("/home", homeHandler)

	// 使用 HandleFunc 注册匿名函数
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("About from HandleFunc!"))
	})

	// 手动将函数转换为 HandlerFunc，再用 Handle 注册
	welcomeHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome from HandlerFunc!"))
	}
	http.Handle("/welcome", http.HandlerFunc(welcomeHandler))

	// 启动服务
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	server.ListenAndServe()
}
