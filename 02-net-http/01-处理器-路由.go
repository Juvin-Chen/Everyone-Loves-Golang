/*
lesson1:

核心内容：
  1. 最简HTTP服务器实现（基础路由+启动）
  2. Handler与ServeMux核心概念拆解
  3. 自定义Handler和Server的综合示例
  4. 解答核心疑惑：
     - http.HandleFunc与mux.HandleFunc的层级问题
     - fmt.Fprintf与w.Write的等价性问题
     - Handler与路由的本质区别
*/

/*
一、HTTP 服务器核心定义
什么是 HTTP 服务器？
- 核心职责：监听指定端口 → 接收客户端（浏览器/Postman）HTTP 请求 → 根据 URL/请求方法分发处理 → 返回 HTTP 响应
- Go 实现基础：net/http 包提供了所有核心组件（Handler、ServeMux、Server 等）
*/

package main

import (
	"fmt"
	"net/http"
)

/*
二、基础示例：最简 HTTP 服务器
核心知识点：
1. fmt.Fprintf 与 w.Write 的区别（语法/逻辑层面）
2. 默认路由器 http.DefaultServeMux 的本质
3. http.HandleFunc 是默认路由器的快捷写法
*/
func testHttpServer() {
	// 1. 注册路由处理器：访问根路径 "/" 时执行匿名函数
	// http.HandleFunc 本质 = http.DefaultServeMux.HandleFunc（默认路由器的快捷写法）
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// w: 响应写入器（实现 io.Writer 接口），用于向客户端返回数据
		// r: 请求对象，包含客户端的请求信息（URL、参数、方法等）

		/*
			fmt.Fprintf 核心定义:
			1. 函数归属与作用
			核心作用：将格式化后的内容写入指定的 “可写入对象”（不是直接打印到控制台，而是写入目标载体）
			2. 函数签名
			func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)

			前提：fmt.Fprintf 对第一个参数的要求:fmt.Fprintf 要求第一个参数 w 必须实现 Go 标准库 io 包中的 io.Writer 接口 —— 这是一个 “写入标准”。
			io.Writer 接口的定义:
				type Writer interface {
					// 只要实现了这个 Write 方法，就符合 io.Writer 标准
					Write(p []byte) (n int, err error)
				}
			任何对象，只要有 Write 方法（能接收字节数组并写入），就能被 fmt.Fprintf 使用。
			http.ResponseWriter 满足这个标准:Go 底层已经帮我们实现了 http.ResponseWriter 的 Write 方法 —— 它本质上是一个 “符合 io.Writer 标准的响应写入工具”。

			如果只是固定字符串，就直接写；如果需要动态内容（比如拼接变量），可以用占位符，比如：
				name := "xiaoming"
				fmt.Fprintf(w, "Hello, %s!", name) // 最终输出：Hello, xiaoming!
		*/
		fmt.Fprintf(w, "Hello World!") // 推荐：固定字符串/格式化场景都能用
		// 等价写法（原生版）：w.Write([]byte("Hello World!"))，Fprintf底层一样的也是调用Write()
	})

	// 2. 启动 HTTP 服务器
	fmt.Println("基础服务器启动：http://localhost:8080")
	/*
	   http.ListenAndServe 核心参数：
	   - 第一个参数：监听地址（":8080" 表示本机所有网卡的 8080 端口）
	   - 第二个参数：处理器（Handler），传 nil 表示使用默认路由器 http.DefaultServeMux
	   - 底层逻辑：http.ListenAndServe(addr, nil) [ 这个方法本质上是在内部调用 server.ListenAndServe() ] ≡ (&http.Server{Addr:addr, Handler:http.DefaultServeMux}).ListenAndServe()
	*/
	http.ListenAndServe(":8080", nil)
}

/*
三、核心概念：Handler 与 ServeMux（解决「不同级」的疑惑）
1. Handler（处理器）：处理具体请求的「执行者」
  - 定义：http.Handler 是接口，仅需实现 ServeHTTP(w ResponseWriter, r *Request) 方法
  - 作用：接收请求 → 执行业务逻辑 → 返回响应
  - 注意：Handler 不是路由，只是「处理请求的函数/结构体」

2. ServeMux（多路复用器/路由器）：分发请求的「调度员」
  - 定义：本质是实现了 Handler 接口的特殊处理器
  - 核心作用：根据请求 URL 路径，将请求分发到不同的 Handler
  - 两种使用方式：
    a) 默认路由器：http.DefaultServeMux（全局单例，http.HandleFunc 是其快捷写法）
    b) 自定义路由器：http.NewServeMux()（独立实例，避免全局冲突）
*/
func practice1() {
	// 自定义路由器 vs 默认路由器
	// 1. 创建自定义路由器实例（独立于 DefaultServeMux）
	mux := http.NewServeMux()

	// 2. 给自定义路由器注册路由（mux.HandleFunc 是原生写法）
	/*
		对比：http.HandleFunc → 给 DefaultServeMux 注册；mux.HandleFunc → 给自定义 mux 注册
		两者是「同一层级」，只是操作的路由器实例不同

		Go 标准库 net/http 源码里的 http.HandleFunc 定义
			func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
				DefaultServeMux.HandleFunc(pattern, handler) // 本质是调用默认路由器的 HandleFunc
			}
	*/
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Home page")
	})
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "About page")
	})

	// 3. 启动服务器：使用自定义路由器（第二个参数传 mux）
	// 对比：传 nil → 用 DefaultServeMux；传 mux → 用自定义路由器
	http.ListenAndServe(":8080", mux)
}

/*
四、自定义 Handler 示例（实现 Handler 接口）
- 步骤：定义结构体 → 实现 ServeHTTP 方法 → 注册/使用该 Handler
*/
// 自定义 Handler 结构体（可添加字段存储配置/数据）
type myHandler struct{}

// 实现 http.Handler 接口的 ServeHTTP 方法（核心）
func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello web")) // 原生写法：固定字符串用 Write 更高效
}

/*
五、综合示例1：显式创建 Server + 混合使用默认路由器/自定义 Handler
- 重点：Server 是 HTTP 服务的核心配置结构体，比 http.ListenAndServe 更灵活
- Handler 传 nil → 复用 DefaultServeMux（已注册自定义 Handler/普通路由）
*/
func demo1_() {
	// 1. 往默认路由器（DefaultServeMux）注册路由/Handler
	/*
		使用 http.Handle 注册自定义 Handler ：
		第一个参数：路由路径（比如 /hello，第二个参数：实现了 http.Handler 接口的结构体实例（必须传指针）
		本质：往默认路由器 http.DefaultServeMux 注册该 Handler
		http.Handle("/hello", helloHandler)

		补充：自定义路由器中使用 mux.Handle（和 http.Handle 逻辑一致，只是路由器不同）
		mux := http.NewServeMux()
		mux.Handle("/hello", helloHandler) // 往自定义路由器注册
	*/
	http.Handle("/", &myHandler{})                                           // 注册自定义 Handler（处理 / 路径）
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { // 注册普通路由
		fmt.Fprintf(w, "welcome to login in")
	})

	// 2. 显式创建 Server 实例（可配置超时、TLS 等，比 http.ListenAndServe 灵活）
	server := http.Server{
		Addr:    "localhost:8080", // 监听地址
		Handler: nil,              // nil → 使用 DefaultServeMux（等价写 http.DefaultServeMux）
	}

	// 3. 启动服务（阻塞调用，需放在最后）
	fmt.Println("综合示例1服务启动：http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("服务启动失败：%v\n", err)
	}
}

/*
六、综合示例2：纯自定义 Handler（无路由器）
- 特点：Server 的 Handler 直接传自定义 Handler 实例
- 结果：所有请求（无论 URL 是 /login 还是 /about）都由该 Handler 处理（无路由分发）
- 对比：demo1_ 有路由（DefaultServeMux），demo2_ 无路由（所有请求走同一个 Handler）
*/
func demo2_() {
	// 1. 创建自定义 Handler 实例
	mh := myHandler{}

	// 2. 创建 Server 实例：Handler 直接指定为自定义 Handler（无路由器）
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: &mh, // 所有请求都交给 mh 处理，无路由分发
	}

	// 3. 启动服务
	fmt.Println("综合示例2服务启动：http://localhost:8080")
	server.ListenAndServe()
}

/*
核心 Q&A（解决所有疑惑）
Q1：为什么 http.HandleFunc 和 mux.HandleFunc 感觉不是同级？
A：两者是同一层级！
   - http.HandleFunc = http.DefaultServeMux.HandleFunc（默认路由器的快捷写法）
   - mux.HandleFunc = 自定义路由器的原生写法
   差异仅在于操作的路由器实例（默认/自定义），核心逻辑完全一致。

Q2：fmt.Fprintf(w, "xxx") 和 w.Write("xxx") 等价吗？
A：语法上不等价（后者报错），逻辑上（结果）等价于 w.Write([]byte("xxx"))；
   - fmt.Fprintf 支持格式化，适合动态内容；
   - w.Write 是原生方法，适合固定字符串（效率更高）。

Q3：Handler 是路由吗？
A：严格来说不是
   - Handler 是「处理请求的执行者」（比如 myHandler）；
   - ServeMux（路由器）是「分发请求的调度员」（比如 DefaultServeMux/mux）；
   - 路由的本质是 ServeMux，它也是一种特殊的 Handler（实现了 ServeHTTP 方法）。

Q4：http.ListenAndServe 和 server.ListenAndServe 的关系？
A：http.ListenAndServe 是 server.ListenAndServe 的封装：
   http.ListenAndServe(addr, handler) ≡ (&http.Server{Addr:addr, Handler:handler}).ListenAndServe()
   - 简单场景用 http.ListenAndServe；
   - 复杂场景（需配置超时/TLS）用显式创建 Server。
*/
