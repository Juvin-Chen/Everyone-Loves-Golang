/*
lesson 4

1.POST 请求与请求体的概念。
2.如何解析 application/x-www-form-urlencoded 表单数据。
3.ParseForm、Form、PostForm、FormValue、PostFormValue 的区别与用法。
4.简单的 JSON 请求体解析。
5.区分 URL 参数和表单参数的重要性。
*/

/*
1. 为什么要用 POST？
前面我们学习了 GET 请求，它把参数放在 URL 的查询字符串里（?name=张三）。
但 GET 有几个限制：
	数据长度有限：URL 长度通常有限制（浏览器限制 2KB 左右），不能提交大量数据。
	数据可见：所有参数都暴露在 URL 上，不适合传输密码等敏感信息。
	数据用途：GET 应当用于“获取数据”，而不应该用于修改数据（比如提交订单、修改资料等）。

POST 请求解决了这些问题：
	数据放在请求体（Body）中，不暴露在 URL 上，相对安全。
	没有长度限制（理论上由服务器配置决定）。
	语义明确：POST 用于“创建”或“提交”数据，符合 RESTful 规范。

在 Web 开发中，表单提交（登录、注册、评论等）几乎都用 POST。

2.HTTP 请求体（Body）是什么？
HTTP 请求由三部分组成：请求行、请求头、请求体。

!!! 请求体是可选的，GET 请求通常没有请求体，而 POST 请求通常包含请求体。

请求体可以是多种格式：
	application/x-www-form-urlencoded：表单数据，格式和查询字符串一样，比如 name=张三&age=25。
	multipart/form-data：用于文件上传，可以包含二进制数据。
	application/json：JSON 格式的数据。
	其他格式如纯文本、XML 等。

我们这节课先学习最常用的 application/x-www-form-urlencoded 表单数据。

3. 如何在 Go 中获取表单数据？
	在 net/http 中，处理 POST 请求的表单数据主要分三步：
	1调用 r.ParseForm() 解析请求体（可选，但建议显式调用以捕获错误）。
	2通过 r.FormValue(key) 或 r.PostFormValue(key) 获取值。
	3如果涉及文件上传，还需要调用 r.ParseMultipartForm()。

*/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 先看一个最简单的例子：接收用户名和密码，并返回一段欢迎信息。
func testPost() {
	// 注册处理 /login 路径的函数
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// 只允许 POST 方法
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "只支持post请求")
			return
		}

		// 解析表单（将请求体中的数据解析到 r.Form 和 r.PostForm）
		err := r.ParseForm()
		if err != nil {
			// 如果解析失败，返回 400 错误
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "表单解析失败: %v", err)
			return
		}

		//  获取表单字段的值
		username := r.FormValue("username") // 从表单中获取 username 字段
		password := r.FormValue("password") // 获取 password 字段

		// 简单的验证（演示用，实际应该查数据库）
		if username == "admin" && password == "123456" {
			fmt.Fprintf(w, "登录成功！欢迎 %s", username)
		} else {
			fmt.Fprintf(w, "用户名或密码错误")
		}
	})
}

/*
测试步骤（打开终端，依次敲命令）：
先启动上面的 Go 程序；
	1.敲第一个命令（发 GET 指令）：
	curl http://localhost:8080
	→ 输出：你发给我的操作指令是：GET（因为 curl 默认发 GET）；

	2.敲第二个命令（发 POST 指令）：
	curl -X POST http://localhost:8080
	→ 输出：你发给我的操作指令是：POST（-X 就是告诉 curl “我要发 POST 指令”）；

	3.敲第三个命令（发 PUT 指令）：
	curl -X PUT http://localhost:8080
	→ 输出：你发给我的操作指令是：PUT；

	如果忘记加 -d 或者用 GET 请求，程序会返回“只支持 POST 请求”。
*/

/*
4. 深入理解 ParseForm、Form、PostForm、FormValue 和 PostFormValue
这几个概念容易混淆，我们逐个拆解：

4.1 r.ParseForm()
	作用：解析请求体中的表单数据（application/x-www-form-urlencoded 格式），并将结果存入 r.Form 和 r.PostForm 中。
	调用时机：必须在访问 r.Form 或 r.PostForm 之前调用。
		如果直接调用 r.FormValue，Go 会 自动调用 ParseForm（但不会返回解析错误），所以显式调用是为了捕获错误。
	错误处理：如果请求体格式错误（比如不是有效的表单数据），ParseForm 会返回错误，我们应该检查并返回 400 状态码。

4.2 r.Form 和 r.PostForm
	r.Form 是一个 url.Values 类型（map[string][]string），包含 URL 查询参数 + 请求体中的表单参数。
		如果同一个 key 既出现在 URL 中又出现在请求体中，URL 中的值会被优先保留？
		实际上，ParseForm 会合并：先解析 URL 查询参数，再解析请求体，请求体中的同名参数会 追加 到切片中（不覆盖），所以 r.Form 中该 key 对应的值是 [url_value, body_value]。

	r.PostForm 只包含 请求体中的表单参数（不包含 URL 参数）。它也是在调用 ParseForm 后填充的。

示例（假设我们既有 URL 参数又有表单参数）：
http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    fmt.Fprintf(w, "Form: %v\n", r.Form)
    fmt.Fprintf(w, "PostForm: %v\n", r.PostForm)
})

请求：
POST /test?name=url_name HTTP/1.1
Content-Type: application/x-www-form-urlencoded
name=body_name

输出：
Form: map[name:[url_name body_name]]
PostForm: map[name:[body_name]]

4.3 r.FormValue(key)
	便捷方法：先自动调用 ParseForm（如果有必要），然后从 r.Form 中取 第一个值（即 r.Form[key][0]）。
	如果 key 不存在，返回空字符串。
	注意：它会同时查找 URL 参数和表单参数，URL 参数优先（因为 r.Form 中 URL 参数排在切片前面）。

4.4 r.PostFormValue(key)
	同样自动调用 ParseForm，但只从 r.PostForm 中取第一个值（仅请求体中的表单数据）。
	如果 key 不存在，返回空字符串。
	这个更常用于只关心提交的表单数据，避免被 URL 参数污染。

建议：在纯 POST 表单处理中，使用 r.PostFormValue 更安全，因为它不会误取 URL 参数。
*/

// 完整示例：区分 URL 参数和表单参数
func testUrlAndPostform() {
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "只支持 POST")
			return
		}

		// 显式解析，可以捕获错误
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "解析失败：%v", err)
			return
		}

		// 1. 便捷方法：只取第一个值
		// 从 Form 中获取
		nameFromForm := r.FormValue("name")
		// 从 PostForm 中获取（仅表单参数）
		nameFromPostForm := r.PostFormValue("name")

		/*
			2. 直接访问切片：获取所有值
			//r.Form["name"] 返回 []string，包含 URL+表单的所有 name 值
			allNamesFromForm := r.Form["name"]
			// r.PostForm["name"] 返回 []string，仅包含表单的所有 name 值
			allNamesFromPostForm := r.PostForm["name"]
		*/

		fmt.Fprintf(w, "r.FormValue(\"name\") = %s\n", nameFromForm)
		fmt.Fprintf(w, "r.PostFormValue(\"name\") = %s\n", nameFromPostForm)
	})
	http.ListenAndServe(":8080", nil)
}

/*
测试：
# 在 URL 中带上 name，同时表单中也提交 name
curl -X POST -d "name=body_name" "http://localhost:8080/submit?name=url_name"

输出：
r.FormValue("name") = url_name
r.PostFormValue("name") = body_name
可以看到 FormValue 取了 URL 参数，而 PostFormValue 只取表单参数。
*/

// 处理 JSON 请求体
/*
除了表单，现在 API 也常用 JSON 格式提交数据。
处理 JSON 需要读取请求体，然后使用 encoding/json 解码。

这种数据不是 “键值对字符串”，而是 JSON 字符串 ——Go 本身不认识这种字符串，没法直接用 r.FormValue 获取。
所以需要把这个 JSON 字符串「转成 Go 能操作的结构体变量」，这个过程就是 JSON 解码（Decode）；反过来，把 Go 结构体转成 JSON 字符串返回给客户端，就是 JSON 编码（Encode）。
*/

// 定义一个结构体，和 JSON 数据的结构对应
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func testJson() {
	http.HandleFunc("/login_json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "只支持 POST")
			return
		}

		// 解析 JSON 请求体
		var req LoginRequest
		// 核心：把请求体里的 JSON 字符串，解析到 req 变量里
		/*
			r.Body是 HTTP 请求体的读取入口，里面存着客户端发过来的 JSON 字符串
			json.NewDecoder(r.Body)	造一个 “JSON 翻译器” 告诉 Go：我要从 r.Body 里读数据，并且这些数据是 JSON 格式的
			Decode(&req)	翻译器把 JSON 转成 Go 结构体
		*/
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "无效的 JSON:%v", err)
			return
		}

		// 最佳实践：延迟关闭 Body，函数结束前执行
		defer r.Body.Close()

		// 模拟验证
		if req.Username == "admin" && req.Password == "123456" {
			fmt.Fprintf(w, "登录成功")
		} else {
			fmt.Fprintf(w, "登录失败")
		}
	})
	http.ListenAndServe(":8080", nil)
}

/*
注意：JSON 解析时，必须确保请求头 Content-Type 是 application/json，否则服务端可能不知道如何解析（但上面的代码没检查，它会尝试解析任何请求体，如果格式不对会报错）。

现在学的是 “解码（客户端→服务器：JSON→结构体）”，反过来 “服务器→客户端：结构体→JSON” 就是编码，比如登录成功后返回 JSON 响应：
// 定义响应结构体
type LoginResponse struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 提示信息
}

// 在登录成功后编码并返回
successResp := LoginResponse{Code: 200, Message: "登录成功"}
// 设置响应头为 JSON 类型
w.Header().Set("Content-Type", "application/json")
// 把结构体编码成 JSON 字符串，写入响应体
json.NewEncoder(w).Encode(successResp)

客户端会收到：{"code":200,"message":"登录成功"}—— 这就是 JSON 编码的作用。
*/

/*
文件上传（multipart/form-data）简介
	文件上传使用 multipart/form-data 编码。
		处理步骤：
			1.调用 r.ParseMultipartForm(maxMemory) 设置内存缓冲大小。
			2.通过 r.FormFile(key) 获取上传的文件。
			3.读取文件内容并保存。
先做一个大致了解，详细地在后面再学习
*/
