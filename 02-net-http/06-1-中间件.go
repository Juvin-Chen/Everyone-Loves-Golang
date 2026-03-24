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

// 中间件的链式调用

/*
实际应用中我们会组合多个中间件，比如：日志 → 认证 → 业务处理。
这可以通过嵌套调用来实现：
handler := authMiddleware(loggerMiddleware(helloHandler))
*/

// 但是当中间件很多时，嵌套写法可读性差。我们可以写一个辅助函数来串联中间件：
// chain 将多个中间件串联起来，按顺序执行
// ...类型 放在函数参数最后，意思是：这个参数可以接收「任意数量」的同类型值（0 个、1 个、多个都可以），函数内部会把这些值自动打包成一个切片（slice） 来处理。
func chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middlewares {
		handler = mw(handler)
	}
	return handler
}

// 用法：
/*
	handler := chain(helloHandler, loggerMiddleware, authMiddleware)
	http.Handle("/", handler)
	注意：中间件的执行顺序是从左到右，即先应用 loggerMiddleware，再应用 authMiddleware。
	请求流程：请求 → logger 前置 → auth 前置 → 业务处理 → auth 后置 → logger 后置 → 响应。
*/

// 编写一个简单的认证中间件
// 假设我们要求某些路径必须带有 Authorization: Bearer <token> 头，且 token 必须等于一个固定值（演示用）。
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取 Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Missing Authorization header")
			return
		}

		// 假设我们期望 Bearer token 是 "secret-token"
		expectedToken := "Bearer secret-token"
		if authHeader != expectedToken {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Invalid token")
			return
		}
		// 认证通过，继续执行
		next.ServeHTTP(w, r)
	})
}

// 我们可以把它应用到需要保护的路由上：
func testProtected() {
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Secret data")
	})
	http.Handle("/secret", protectedHandler)
}

/*
测试：
# 不带 token
curl -v http://localhost:8080/secret
# 返回 401

# 带正确 token
curl -H "Authorization: Bearer secret-token" http://localhost:8080/secret
# 返回 "Secret data"
*/

// 在中间件中传递数据（使用 Context）
/*
有时我们希望在中间件中计算一些值，然后传递给后续的处理器。
例如，认证中间件解析出用户信息，业务处理器需要知道当前用户是谁。
这可以通过 context.Context 实现。
*/

// 补充内容
/*
使用 http.StripPrefix 和 http.FileServer 的中间件思想
标准库中已经有一些中间件风格的函数，比如 http.StripPrefix 和 http.TimeoutHandler。
	http.StripPrefix 接收一个前缀和一个 Handler，返回一个删除了该前缀再转发给 Handler 的新 Handler。
	这本质上也是一个中间件。

	// 将 /static/ 映射到 ./public 目录，并去掉 /static/ 前缀
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))


第三方中间件库
实际开发中，我们常使用现成的中间件，如 github.com/gorilla/handlers 提供的 LoggingHandler、CORS 等。
用法示例：
	import "github.com/gorilla/handlers"
	handler := handlers.LoggingHandler(os.Stdout, myHandler)
还有 github.com/rs/cors 专门处理跨域。
这些库提供了更健壮、易用的实现。
*/

/*
常见问题
Q：中间件的执行顺序为什么重要？
A：顺序决定了前置和后置操作的执行时机。例如，认证中间件应该在日志中间件之后？不一定，日志通常在最外层，这样能记录包括认证失败在内的所有请求。一般顺序：recovery → 日志 → 认证 → 业务。

Q：为什么中间件要返回 http.Handler 而不是直接修改原来的？
A：返回新 Handler 保持了函数的纯粹性和可组合性，便于链式调用，且不影响原始 Handler。

Q：可以在中间件中直接修改响应吗？
A：可以。比如认证失败时直接写回 401 并返回，不再调用 next.ServeHTTP。这就是中间件的“短路”效果。

Q：中间件中修改了 http.ResponseWriter 会怎样？
A：如果你传递的是 http.ResponseWriter 接口，修改会影响后续处理器。但要注意，有些中间件会包装 ResponseWriter（例如记录状态码、内容长度），需要实现一个自定义的 ResponseWriter 类型。
*/
