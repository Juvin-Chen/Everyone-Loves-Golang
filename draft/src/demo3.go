// 内置的 Handlers

package main

import "net/http"

// 1.404 处理器，返回一个预定义的 Handler，对所有请求统一返回 404 page not found 响应。
// func NotFoundHandler() Handler
func demo3_1() {
	// 所有访问 /404 的请求，都返回 404 页面
	http.Handle("/404", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

// 2.http.RedirectHandler：重定向处理器，返回一个 Handler，将所有请求重定向到指定 URL，并使用指定的 3xx 状态码。
// func RedirectHandler(url string, code int) Handler
func demo3_2() {
	// 给/new路径加个处理器，让跳转后能看到内容
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("我是新地址 /new 的内容！"))
	})

	// 访问 /old 时，永久重定向到 /new
	http.Handle("/old", http.RedirectHandler("/new", http.StatusMovedPermanently))
	http.ListenAndServe(":8080", nil)
}

// 3.http.StripPrefix：URL 前缀处理器，返回一个 Handler，先从请求 URL 中移除指定前缀，再将修改后的请求交给下一个 Handler 处理；如果 URL 不匹配前缀，返回 404。
// func StripPrefix(prefix string, h Handler) Handler
func demo3_3() {

}
