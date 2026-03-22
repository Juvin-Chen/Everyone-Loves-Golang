/*
lesson 2
*/

/*
一、一个 HTTP 响应由哪几部分组成？
当你的浏览器向服务器发起请求，服务器会返回一段数据，这就是 HTTP 响应。
一个完整的 HTTP 响应包含三部分（就像一封信包含信封、信头、正文）：
1.状态行
例如：HTTP/1.1 200 OK，它告诉我们 HTTP 协议的版本、状态码（200）和状态文字（OK）。

2.响应头（Headers）
一些键值对，描述响应的附加信息，比如内容类型、内容长度、服务器名称等。例如：
Content-Type: text/html; charset=utf-8
Content-Length: 1234
Server: nginx

3.响应体（Body）
实际的数据，比如 HTML 代码、JSON 字符串、图片二进制等。

浏览器收到响应后，会根据状态行判断是否成功，根据响应头决定如何解析响应体，最后显示内容。

二、响应头的作用（为什么要设置它们？）
你可以把响应头想象成快递包裹上的标签。
包裹里是商品（响应体），标签上写着“易碎品”“冷藏”“收件人地址”等信息，快递员（浏览器）根据标签来决定如何搬运和处理。没有标签，快递员可能不知道里面是什么，就会用错误的方式处理。

常见响应头示例：
响应头				作用						举例
Content-Type		告诉浏览器响应体的数据格式	  text/html 表示 HTML，application/json 表示 JSON，image/png 表示 PNG 图片
Content-Length		响应体的大小（字节数）	      Content-Length: 348
Server				服务器软件的名称			 Server: nginx/1.18.0
Cache-Control		控制浏览器如何缓存			 Cache-Control: no-cache 表示不缓存
Set-Cookie			让浏览器保存 				Cookie	Set-Cookie: sessionId=abc123
如果没有正确设置 Content-Type，浏览器可能会把 JSON 当作纯文本显示，或者把 HTML 当作纯文本，导致页面乱码或无法渲染。

三、在 Go 中如何操作响应头？
在 net/http 包中，处理函数接收的 http.ResponseWriter 参数（我们习惯命名为 w）提供了设置响应头的方法。

ResponseWriter 是 Go 为 HTTP 服务端定义的 “响应写入器” 接口，而非具体的结构体。
标准库的设计思路是：暴露接口给开发者使用，内部用不同结构体实现该接口，适配不同场景（比如普通响应、分块响应、缓冲响应）。

ResponseWriter 接口的官方定义（简化版）:

package http

// ResponseWriter 是服务器向 HTTP 响应写入数据的核心接口
type ResponseWriter interface {
    // 核心：Write 方法，必须实现（满足 io.Writer 接口）
    Write([]byte) (int, error)

    // 发送 HTTP 状态码（如 404、500），可省略（默认 200 OK）
    WriteHeader(int)

    // 获取/设置响应头（比如 Content-Type、Set-Cookie）
    Header() Header
}

1. w/r 参数来源：handler 函数 `func(w http.ResponseWriter, r *http.Request)` 中的 w/r 并非手动创建，而是 `net/http` 库作为回调参数自动传入：
   - r：库解析 TCP 连接中的 HTTP 报文，生成 `http.Request` 实例；
   - w：库基于 TCP 连接创建 `defaultResponseWriter`（`ResponseWriter` 接口的默认实现），以接口形式传入。
2. defaultResponseWriter 角色：日常开发中 handler 拿到的 w，底层默认是该私有结构体，实现了 `ResponseWriter` 接口的 `Write/WriteHeader/Header` 方法，负责实际向客户端写入响应、管理状态码/响应头。
3. 核心逻辑：开发者只需定义 handler 处理逻辑，`net/http` 库负责：监听连接 → 创建 w/r → 调用 handler 并传参 → 最终通过 w 的底层实现完成响应发送。


3.1 w.Header() 返回一个 http.Header 对象
	http.Header 其实是一个 map[string][]string，用来存储所有的响应头。
	你可以通过 w.Header().Set(key, value) 来设置一个键值对。

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	Set 方法：如果该键不存在，则添加；如果已存在，则 覆盖 原来的值（只保留一个值）。

	如果你希望同一个键有多个值（例如 Set-Cookie 可能会设置多个 Cookie），可以用 Add 方法：
		w.Header().Add("Set-Cookie", "session=abc")
		w.Header().Add("Set-Cookie", "user=zhang")
	这样 Set-Cookie 对应的值就是 ["session=abc", "user=zhang"]。

	删除某个头：
		w.Header().Del("Content-Type")。

3.2 设置状态码：w.WriteHeader(statusCode)
状态码是三位数字，表示请求的处理结果。
常见的状态码有：
	200 OK：成功
	301 Moved Permanently：永久重定向
	302 Found：临时重定向
	400 Bad Request：客户端请求有误
	404 Not Found：资源不存在
	500 Internal Server Error：服务器内部错误

在 Go 中，我们可以用 w.WriteHeader(200) 来发送状态码。
重要：必须在写入响应体之前调用 WriteHeader，因为状态行和响应头是在响应体之前发送的。
如果你没有显式调用 WriteHeader，Go 会在你第一次调用 w.Write() 或 fmt.Fprintf(w, ...) 时自动发送 200 OK。
w.WriteHeader(http.StatusNotFound) // 发送 404
这里 http.StatusNotFound 是 Go 定义好的常量，值为 404，这样写比直接写数字更易读。

3.3 写入响应体
你可以用 fmt.Fprintf(w, ...) 或 w.Write([]byte(...)) 来写入响应体。
一旦开始写入响应体，就不能再修改响应头了（因为头已经发送出去了）。

提示： HTTP 协议规定：响应行（状态行）和响应头必须在响应体之前发送。
一旦你开始写入响应体（例如调用 fmt.Fprintf(w, ...)），Go 会自动把之前设置的响应头发送出去，然后发送响应体。
*/

package main

import (
	"fmt"
	"log"
	"net/http"
)

func testResponseWriter() {
	/*
		源码解析：
		// net/http 包中 Header 的核心定义
		package http

		// Header 是 HTTP 头的类型别名，本质就是 map[string][]string
		type Header map[string][]string

		w.Header() 返回一个 http.Header 对象,http.Header 其实是一个 map[string][]string，用来存储所有的响应头。
	*/

	// 注册根路径的处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 设置响应头，告诉浏览器返回的是纯文本，UTF-8 编码
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// 告诉浏览器不要缓存这个页面（仅作演示）
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		// 也可以添加一个自定义头（虽然不常见，但合法）
		w.Header().Set("X-My-Custom-Header", "Hello from Go")

		// 设置状态码
		// 显式发送 200 状态码
		w.WriteHeader(http.StatusOK) // http.StatusOK 是常量，值为 200

		// 写入响应体
		// 注意：必须在上面的头设置完成之后再写入响应体
		fmt.Fprintf(w, "Hello, World! 这是一段纯文本响应。")
	})

	// 启动服务器，监听 8080 端口
	fmt.Println("服务器已启动，访问 http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}

/*
总结:
	HTTP 响应由 状态行、响应头、响应体 三部分组成。
		1.响应头是键值对，告诉浏览器如何解析和处理响应体。
		2.在 Go 中，通过 w.Header().Set(key, value) 设置响应头。
		3.通过 w.WriteHeader(code) 设置状态码（必须在写入响应体之前调用）。
		4.用 fmt.Fprintf(w, ...) 或 w.Write([]byte(...)) 写入响应体。
		5.用浏览器开发者工具或 curl 可以查看响应头，验证你的设置。
*/
