/*
Gin 简介与第一个程序
1. Gin 是什么？
Gin 是 Go 语言中一个非常流行的 Web 框架，它基于 Go 标准库 net/http 进行了封装和增强。
它的主要特点：
	高性能：路由使用了 httprouter，速度非常快。
	轻量：核心代码很少，但功能齐全。
	中间件机制：可以方便地插入日志、认证、恢复等中间件。
	参数绑定与验证：自动将请求参数（路径参数、查询参数、表单、JSON）绑定到结构体，并支持验证。
	JSON 渲染：返回 JSON 只需要一行代码。
	错误管理：统一错误处理机制。

可以把 Gin 看作是 net/http 的“加强版”，它省去了很多重复劳动，让代码更简洁、可读性更高。

2. 安装 Gin
在开始之前，确保你的 Go 版本 ≥ 1.16（推荐 1.18+）。
打开终端，创建一个新项目并安装 Gin：
	mkdir gin-tutorial
	cd gin-tutorial
	go mod init gin-tutorial          # 初始化 Go Module
	go get -u github.com/gin-gonic/gin # 下载 Gin 框架
安装完成后，go.mod 文件中会多出 github.com/gin-gonic/gin 的依赖。
*/

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func test1_ReturnJson() {
	// 创建 gin 引擎
	/*
		gin.Default() 创建一个默认的 Gin 引擎。
		这个引擎已经内置了两个中间件：Logger（记录请求日志）和 Recovery（捕获 panic 并返回 500）。
		你也可以用 gin.New() 创建一个不带任何中间件的引擎，然后自己添加需要的中间件。
	*/
	r := gin.Default()

	// 注册一个 GET 路由
	/*
		r.GET("/ping", func(c *gin.Context) { ... })
		r.GET 是 Gin 提供的方法，用于注册处理 GET 请求的路由。
		第一个参数是路径（"/ping"），第二个参数是一个处理函数。
		处理函数的类型是 gin.HandlerFunc，它接收一个 *gin.Context 参数。
		gin.Context 是 Gin 的核心，它封装了 http.ResponseWriter 和 *http.Request，并提供了很多便捷方法（如 JSON、Bind、Param 等）。
	*/
	r.GET("/ping", func(c *gin.Context) {
		// 返回 JSON 响应
		/*
			c.JSON(http.StatusOK, gin.H{ "message": "pong" })
			c.JSON 方法返回一个 JSON 格式的响应。
				第一个参数是 HTTP 状态码，这里使用 http.StatusOK（常量 200）。
				第二个参数是要返回的数据。gin.H 是 map[string]interface{} 的类型别名，用于快速构建 JSON 对象。
			Gin 会自动设置响应头 Content-Type: application/json，并将数据序列化为 JSON 字符串。
		*/
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 启动服务器，默认监听 8080 端口
	/*
		启动 HTTP 服务器，默认监听 0.0.0.0:8080。
		你也可以传入地址参数，比如 r.Run(":9090") 监听 9090 端口。
		Run 方法会阻塞，直到服务器关闭。
	*/
	r.Run()
}

/*
如果用 net/http 实现同样的功能，代码会是这样的：
import (
    "encoding/json"
    "net/http"
)

func main() {
    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
    })
    http.ListenAndServe(":8080", nil)
}

Gin 的优势很明显：
	不需要手动设置 Content-Type。
	不需要手动编码 JSON。
	路由注册更简洁，支持多种 HTTP 方法（GET、POST、PUT、DELETE 等）。
	处理函数中直接使用 c 对象，提供了丰富的辅助方法。
*/

// Gin 也支持返回纯文本和 HTML
func test1_More() {
	r := gin.Default()
	// 返回纯文本
	r.GET("/text", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello World!")
	})
	// 返回HTML
	// 注意：使用 c.HTML 之前需要配置模板目录
	r.GET("/html", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Gin示例",
		})
	})
}

// 关于路由匹配
// Gin 的路由支持动态参数、通配符等，例如：
func test1_Path() {
	r := gin.Default()
	// 路径参数 :name
	r.GET("/user/name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "Hello %s", name)
	})
	// 查询参数
	r.GET("/search", func(ctx *gin.Context) {
		query := ctx.DefaultQuery("q", "default")
		ctx.String(http.StatusOK, "Search: %s", query)
	})
}
