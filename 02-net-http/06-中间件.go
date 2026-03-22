/*
lesson 6
这是构建可维护 Web 应用的核心概念，它让你可以在处理请求前后插入通用逻辑，比如日志、鉴权、恢复 panic、跨域处理等，而无需在每个处理函数里重复编写相同代码。

1.中间件的定义：接收 http.Handler，返回 http.Handler 的函数。
2.如何用 http.HandlerFunc 快速实现中间件。
3.链式调用中间件的方法。
4.通过 context 在中间件之间传递数据。
5.标准库中已有的中间件式函数。
6.使用第三方中间件简化开发。
*/

/*
1. 什么是中间件？
中间件（Middleware）就是一个 函数，它接收一个 http.Handler，并返回一个新的 http.Handler。
这个新的 Handler 可以在调用原始的 Handler 之前或之后执行一些额外操作，比如：

	记录请求开始和结束时间（日志）
	验证用户身份（认证）
	如果发生 panic，恢复程序（recovery）
	添加公共响应头（如 CORS）
	压缩响应内容（gzip）

在 Go 中，中间件充分利用了 http.Handler 接口和函数组合的特性。

2. 回顾 http.Handler 和 http.HandlerFunc
http.Handler 是一个接口：

	type Handler interface {
	    ServeHTTP(ResponseWriter, *Request)
	}

任何实现了 ServeHTTP 方法的类型都可以作为 Handler。
而 http.HandlerFunc 是一个函数类型，它实现了 ServeHTTP 方法，因此可以将普通函数转换为 Handler：
type HandlerFunc func(ResponseWriter, *Request)

	func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	    f(w, r)
	}

这意味着我们可以将 func(w http.ResponseWriter, r *http.Request) 这种函数直接当作 http.Handler 使用。
中间件通常也遵循这种模式：它们接收一个 http.Handler，然后返回一个新的 http.Handler（通常是 http.HandlerFunc）。
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// 一个最简单的中间件示例:写一个日志中间件，它会在请求处理前后打印日志。
func loggerMiddleware(next http.Handler) http.Handler {
	// 返回一个 HandlerFunc，它会在调用 next 之前和之后做日志记录
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		// 调用下一个处理器（可能是另一个中间件，也可能是最终的业务处理器）
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

func testLogMiddleware() {
	// 最终的业务处理器
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})

	// 将业务处理器包裹在中间件中
	loggedHandler := loggerMiddleware(helloHandler)

	// 注册并启动
	http.Handle("/", loggedHandler)
	http.ListenAndServe(":8080", nil)
}

/*
运行测试：
访问 http://localhost:8080，控制台会输出类似：
2024/03/22 15:04:05 Started GET /
2024/03/22 15:04:05 Completed / in 123.456µs

中间件中闭包核心笔记
1. 闭包载体：`loggerMiddleware` 返回的 `http.HandlerFunc` 是闭包，捕获外层的 `next`（核心业务处理器）；
2. 生命周期延长：`next` 本是栈变量，因被闭包捕获“逃逸”到堆，只要闭包（`loggedHandler`）存在，`next` 就不会被GC销毁；
3. 执行逻辑：闭包包裹核心逻辑，按「前置日志 → 调用 `next.ServeHTTP` 执行业务 → 后置日志」顺序执行；
4. 核心价值：解耦通用逻辑（日志）与业务逻辑，无需修改原处理器代码，实现功能增强。
*/
