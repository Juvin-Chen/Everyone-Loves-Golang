/*
lesson 3
获取请求信息（URL、方法、查询参数）
*/

/*
一、一个 HTTP 请求由哪几部分组成？
在上一课，我们知道了 HTTP 响应由 状态行、响应头、响应体 组成。

对应的，HTTP 请求也由三部分组成：
1.请求行
例如：GET /hello?name=张三 HTTP/1.1
包含：方法（GET）、路径（/hello）、查询参数（?name=张三）、HTTP 版本（HTTP/1.1）。

2.请求头
键值对，描述请求的附加信息，比如 User-Agent（浏览器类型）、Content-Type（请求体的类型）、Accept（客户端能接收的格式）等。

3.请求体
只在 POST、PUT 等方法中携带数据，比如表单数据、JSON 字符串等。

我们这节课先关注 1.请求行 里的信息，即：
请求方法（Method）
URL 路径（Path）
查询参数（Query String，即 ? 后面的部分）

二、Go 中的 http.Request 结构体
在 net/http 包中，处理函数接收的第二个参数是 *http.Request（通常命名为 r），它包含了客户端请求的所有信息。
这个结构体有很多字段，我们目前主要关注三个：
r.Method：字符串，表示请求方法（如 "GET"、"POST"）。
r.URL：类型是 *url.URL，它又包含多个字段，我们主要用 r.URL.Path 和 r.URL.Query()。
*/

package main

import (
	"fmt"
	"net/http"
	"strings"
)

/*
获取请求方法（Method）
请求方法指明了客户端想对资源进行什么操作。最常见的有：

	GET：获取资源
	POST：创建资源或提交数据
	PUT：更新资源
	DELETE：删除资源
	等等

在 Go 中，直接读取 r.Method 就可以得到方法的字符串。
*/
func testGetMethod() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 获取请求方法，是一个字符串
		method := r.Method
		fmt.Fprintf(w, "你使用的请求方法是: %s\n", method)
	})
	http.ListenAndServe(":8080", nil)
}

/*
注意：r.Method 始终是大写的（如 "GET"），Go 在 net/http 包中定义了常量方便比较：http.MethodGet、http.MethodPost 等。例如：
if r.Method == http.MethodPost {
    // 处理 POST 请求
}
*/

// 获取 URL 路径,URL 路径就是域名后面、? 之前的部分。例如：http://localhost:8080/hello?name=张三 的路径是 /hello。在 Go 中，通过 r.URL.Path 获取。
func testGetUrl() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 获取 URL 路径
		path := r.URL.Path
		fmt.Fprintf(w, "你请求的路径是: %s\n", path)
	})
	http.ListenAndServe(":8080", nil)
}

// 获取查询参数（Query String）
/*
查询参数是 URL 中 ? 之后的部分，例如 ?name=张三&age=25。
这些参数通常用于 GET 请求的筛选、分页等。
在 Go 中，通过 r.URL.Query() 获取所有查询参数，它返回一个 url.Values 类型。
url.Values 本质是 map[string][]string，（[]string代表一个字符串切片，一个key可存多个值，区分于 map[string]string），因为一个参数名可能对应多个值（如 ?hobby=篮球&hobby=足球）。
*/
func testGetQuery() {
	/*
		基本用法：获取单个参数值
		url.Values 提供了 Get(key) 方法，返回该 key 对应的 第一个值，如果 key 不存在则返回空字符串。
	*/
	http.HandleFunc("/great", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		name := query.Get("name")
		if name == "" {
			name = "陌生人"
		}
		fmt.Fprintf(w, "你好, %s!", name)
	})
}

// 获取所有值（当参数有多个时）
// 如果你需要获取一个参数的所有值（比如用户选择了多个爱好），可以直接通过 query["hobby"] 获取 []string 切片。
func testGetMoreQuery() {
	http.HandleFunc("/hobby", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// 获取 hobby 参数的所有值
		hobbies := query["hobby"] // 返回 []string

		if len(hobbies) == 0 {
			fmt.Fprintf(w, "你没有提供任何爱好。")
			return
		}

		result := strings.Join(hobbies, "、")
		fmt.Fprintf(w, "你的爱好有: %s", result)
	})
}

/*
注意：直接使用 query["hobby"] 时，如果参数不存在，会返回一个空切片（长度为0），不会报错。

查询参数的值可能包含特殊字符（URL 编码）
在 URL 中，某些字符（如中文、空格、& 等）需要被编码，例如空格变成 %20，中文变成 %E5%BC%A0 等形式。
r.URL.Query() 会自动解码，所以你拿到的值已经是原始字符串了。
例如：?name=%E5%BC%A0%E4%B8%89 会被解码为 "张三"。
*/

// 综合示例：同时获取方法、路径和查询参数
func testWhole() {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method

		path := r.URL.Path

		query := r.URL.Query()

		// 开始构造响应
		fmt.Fprintf(w, "请求方法: %s\n", method)
		fmt.Fprintf(w, "请求路径: %s\n", path)
		fmt.Fprintf(w, "查询参数:\n")
		if len(query) == 0 {
			fmt.Fprintf(w, "  无\n")
		} else {
			for k, v := range query {
				fmt.Fprintf(w, "  %s = %v\n", k, v)
			}
		}
	})
	fmt.Println("服务器已启动，访问 http://localhost:8080/info?name=张三&age=25")
	http.ListenAndServe(":8080", nil)
}

/*
深入理解：r.URL 的结构
r.URL 是 *url.URL 类型，它包含很多字段，除了 Path 和 RawQuery，还有 Scheme（协议，如 http）、Host（域名+端口）等。
我们暂时不需要全掌握，但了解一下 RawQuery 会帮助我们理解 Query() 的工作方式。
r.URL.RawQuery 是原始的查询字符串（未解码），例如 "name=%E5%BC%A0%E4%B8%89&age=25"。
r.URL.Query() 内部就是解析 RawQuery 并返回 url.Values。
所以，如果只是获取参数值，直接用 Query() 最方便。
*/

/*
重要：路径参数（Path Parameters） vs 查询参数（Query String）

我们经常听到“参数”这个词，但要注意区分两种：
查询参数：在 ? 后面，以键值对形式存在，适合过滤、分页等。
例如：/users?page=2&limit=10

路径参数：是路径的一部分，常用于 RESTful API 中标识资源。
例如：/users/123，这里的 123 就是路径参数。

在 net/http 标准库中，没有直接支持路径参数，需要自己解析路径（比如用 strings.Split 或正则）。我们会在后面介绍如何实现。

补充：r.FormValue 的说明（提前了解）
你可能在网上的例子中见过 r.FormValue("name") 也能获取参数。
r.FormValue 不仅能获取查询参数，还能获取 POST 表单中的参数（请求体里的）。
在只有 GET 请求的情况下，r.FormValue 和 r.URL.Query().Get 效果一样，但内部机制略有不同。
我们会在学习表单处理时详细讲解 r.FormValue，目前先知道这个函数存在即可。
*/

// 注意：这是一个不好的示例，不要用一个处理器借助内部逻辑处理不同路径，最好的方式就是分开注册路由
func BadDemo() {
	// 只注册一个处理器（处理所有路径）
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 核心：通过 r.URL.Path 判断请求的路径
		path := r.URL.Path

		// 根据不同路径执行不同逻辑
		switch path {
		case "/":
			fmt.Fprintf(w, "欢迎访问首页！")
		case "/about":
			fmt.Fprintf(w, "这是关于页面！")
		case "/login":
			fmt.Fprintf(w, "欢迎登录！")
		default:
			// 路径不存在时返回 404
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "页面不存在：%s", path)
		}
	})
}
