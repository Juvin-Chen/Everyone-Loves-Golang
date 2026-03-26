/*
lesson 9
静态文件服务
学习如何在 Go 的 Web 应用中提供 CSS、JavaScript、图片等静态文件。这些资源通常不经过模板处理，直接返回给浏览器。
*/

/*
1. 为什么需要静态文件服务？
在 Web 开发中，除了动态生成的 HTML，还有大量静态资源：
# CSS：控制页面样式
# JavaScript：前端交互逻辑
图片：logo、背景等
字体文件、favicon 等
这些文件不需要经过后端逻辑，直接由 HTTP 服务器返回给客户端即可。
在 Go 标准库中，net/http 提供了 http.FileServer 来高效地提供静态文件服务。

2. http.FileServer 的基本用法
http.FileServer 是一个处理器（http.Handler），它接收一个 http.FileSystem 接口（通常用 http.Dir 实现），然后将指定目录下的文件映射到 HTTP 路径上。

最简单的例子：把当前目录下的 static 文件夹暴露在 /static/ 路径下。
设项目结构：
project/
├── main.go
└── static/

	└── style.css
*/
package main

import "net/http"

// 承接2.中的简单例子 /static/
func testStatic() {
	// 创建文件服务器处理器，提供 ./static 目录下的文件
	fs := http.FileServer(http.Dir("./static"))

	// 将 /static/ 路径映射到文件服务器
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 可选：添加一个简单的首页，引用 CSS
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
            <html>
                <head><link rel="stylesheet" href="/static/style.css"></head>
                <body><h1>Hello, 静态文件!</h1></body>
            </html>
        `))
	})

	http.ListenAndServe(":8080", nil)
}

/*
解释：
	http.Dir("./static") 将目录 ./static 包装成 http.FileSystem，文件系统根目录就是这个文件夹。
	http.FileServer(...) 创建一个处理器，它会处理该目录下的文件请求。
	http.Handle("/static/", http.StripPrefix("/static/", fs))：
		路径 /static/style.css 会被匹配到 /static/ 前缀。
		http.StripPrefix 会移除请求路径中的 /static/ 前缀，剩下 style.css，然后交给 fs 处理器去查找 ./static/style.css 文件。
		如果不用 StripPrefix，fs 会在 ./static/static/style.css 里找文件，导致 404。
访问 http://localhost:8080/static/style.css 就能直接看到 CSS 内容。
*/

/*
3.http.StripPrefix 的作用
http.StripPrefix 本身也是一个中间件函数，它接收一个前缀字符串和一个 http.Handler，返回一个新的 http.Handler。
这个新处理器在调用内部处理器之前，会先去掉请求路径中的指定前缀。

为什么要用？
因为 http.FileServer 期望访问路径直接对应文件路径。
如果你把文件服务器挂载在 /static/ 下，客户端请求 /static/style.css 时，服务器需要查找 ./static/style.css。
StripPrefix("/static/", fs) 会把 /static/style.css 变为 style.css，再交给 fs 处理，这样 fs 就会去找 ./static/style.css。
*/
