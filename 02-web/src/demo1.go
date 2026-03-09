// 创建一个最简单的 web 应用程序
// go 语言仅靠标准库就可构建 web 端应用程序，相对别的语言要简单很多

package main

import "net/http"

func demo1() {
	// 第一个参数是路由地址，/为根地址
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	// 设置 web 服务器
	http.ListenAndServe("localhost:8080", nil) // default servemux
}
