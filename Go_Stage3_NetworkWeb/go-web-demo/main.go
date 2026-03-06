// 在底层的 Socket 中，我们只知道接收字符串。
// 但在 Web 框架中，服务端需要根据用户请求的不同 URL 路径（如 /login, /register）和 方法（如 GET, POST），执行不同的 Go 函数。
// 这个分发请求的机制，就叫 路由。

package main

import (
	"fmt"
	"go-web-demo/handler"    // 导入处理器包
	"go-web-demo/middleware" // 导入中间件包
	"net/http"
)

func main() {
	// 注册基础路由
	http.HandleFunc("/hello", handler.HelloHandler)
	http.HandleFunc("/user", handler.UserHandler)
	http.HandleFunc("/json", handler.JSONResponseHandler)

	// 注册带中间件的路由
	http.HandleFunc("/hello-middleware", middleware.LoggingMiddleware(handler.HelloWithMiddleware))

	fmt.Println("Web 服务已启动，监听 8080 端口...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("服务器启动失败：", err)
	}
}
