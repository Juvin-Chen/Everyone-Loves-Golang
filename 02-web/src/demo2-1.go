// 如何处理（Handle）Web 请求
// http.Handle函数 / http.HandleFunc函数

/*
1.创建 Web Server
http.ListenAndServe(x, x)
-第一个参数是网络地址
  -如果为 ""，那么就是所有网络接口的 80 端口
-第二个参数是 handler
  -如果为 nil，那么就是 DefaultServeMux
  -DefaultServeMux 是一个 multiplexer（多路路由器），它也是个Handler，实现了接口 ServeHTTP () 方法


2.http.ListenAndServe() 源代码解析
(1). 顶层 ListenAndServe() 函数源码
// ListenAndServe 监听指定的 TCP 网络地址 addr，然后调用 Serve 方法处理传入的连接
// 参数说明：
// - addr: 监听地址（如 ":8080"）
// - handler: 处理请求的处理器，nil 表示使用默认的 DefaultServeMux（默认路由器）

func ListenAndServe(addr string, handler Handler) error {
    // 1. 创建 http.Server 结构体实例（这是 Web 服务的核心配置）
    server := &Server{Addr: addr, Handler: handler}
    // 2. 调用 Server 的 ListenAndServe 方法启动服务
    return server.ListenAndServe()
}

http.Server 这是一个 struct，Addr 字段表示网络地址，如果为 ""，那么就是所有网络接口的 80 端口
Handler 字段，如果为 nil，那么就是 DefaultServeMux
本质上是调用了 Server 的 ListenAndServe () 函数，比调用http.ListenAndServe()更灵活


3.https，add SSL/TLS
如果要启动 HTTPS 服务，使用的是 http.ListenAndServeTLS(addr, certFile, keyFile, handler)。它比普通的启动多了两个参数：公钥证书文件和私钥文件。
http.ListenAndServeTLS() / server.ListenAndServeTLS()


4.Handler到底是什么？
handler 是一个接口（interface），handler 定义了一个方法 ServeHTTP ()
type Handler interface {
	// http.ResponseWriter 是一个接口（interface），*http.Request 是一个结构体指针（struct pointer）。
    ServeHTTP(ResponseWriter, *Request)
}


Q & A
1. 核心问题：myHandler 算不算“路由”？
严格来说，代码里写的 myHandler 不叫“路由”（Router），它只是一个纯粹的“处理器”（Handler）。
参考 demo2.md 中
上半部分（HTTP 请求 -> MyHandler）： 对应你的 demo2_2()。所有的请求（不管你是访问 /login 还是 /pay），全部都一股脑交给了 myHandler。它没有做任何“分发”工作，就是个无情的打桩机，所有请求都返回 "Hello web"。这不叫路由。
下半部分（HTTP 请求 -> DefaultServeMux -> Handler 1, 2, 3）： 这个 DefaultServeMux 才是真正的路由（多路复用器 Multiplexer）。它的工作是“看人下菜碟”：如果请求是 /a，它就交给 Handler 1；如果是 /b，就交给 Handler 2。

总结：
路由（ServeMux）的本质，其实也是一个特殊的 Handler（因为它也实现了 ServeHTTP 方法）。只不过它的 ServeHTTP 内部逻辑是：根据 URL 路径，把请求转发给其他具体的 Handler。
*/

package main

import "net/http"

type myHandler struct{}

// 自己定义一个handler
func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello web"))
}

func demo2() {
	// 根据 2. & 3.
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: nil,
	}
	server.ListenAndServe()
	// 上面的写法更灵活一些，上下两个写法等价
	http.ListenAndServe("localhost:8080", nil)
}

func demo2_() {
	// 这种自定义路由器的方式，不同路径都会走到这个Handler，都会把这个hello world 输出
	mh := myHandler{}
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: &mh,
	}
	server.ListenAndServe()
}
