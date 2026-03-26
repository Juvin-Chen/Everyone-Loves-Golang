/*
lesson 10
表单验证与重定向
*/

/*
1. 为什么需要表单验证？
	用户提交的数据可能：
		为空（未填写）
		格式不正确（如邮箱、手机号）
		长度超出限制
		包含恶意内容（如 XSS 攻击）
	服务端验证是安全的第一道防线。虽然前端也可以做验证（提升用户体验），但服务端验证绝对不可省略，因为前端验证可以被绕过。

2. 表单验证的基本流程
	接收用户提交的 POST 请求
	解析表单数据（ParseForm）
	对每个字段进行校验
	如果校验失败，返回错误信息（通常携带错误信息回到表单页面）
	如果校验通过，执行业务逻辑（如保存到数据库），然后重定向到成功页面

3. 重定向（Redirect）
HTTP 重定向是服务器告诉浏览器“你要访问的页面在另一个地址”，浏览器会自动跳转。
在 Go 中，使用 http.Redirect 函数：
http.Redirect(w, r, "/success", http.StatusSeeOther)
参数1：ResponseWriter
参数2：*Request
参数3：目标 URL
参数4：状态码（常见的有 http.StatusFound (302) 临时重定向，http.StatusSeeOther (303) 常用于 POST 后重定向，防止重复提交）

为什么 POST 后要重定向？
如果用户提交表单后直接返回 HTML 页面，刷新浏览器会再次提交表单（重复提交）。通过重定向到另一个页面（GET 请求），刷新就不会重复提交数据。

4. 完整示例：简单的登录表单验证
4.1 项目结构
text
login/
├── main.go
├── templates/
│   ├── login.html
│   └── dashboard.html
*/

package main

import (
	"html/template"
	"net/http"
)

func testReDirect() {
	// 登录页面（GET）
	http.HandleFunc("/login", loginHandler)
	// 处理登录提交（POST）
	http.HandleFunc("/login/submit", loginSubmitHandler)
	// 仪表盘页面（需要登录）
	http.HandleFunc("/dashboard", dashboardHandler)
	// 退出登录
	http.HandleFunc("/logout", logoutHandler)

	http.ListenAndServe(":8080", nil)
}

// 这是 Go 提供的便捷工具函数，专门用来简化模板加载的错误处理
/*
如果模板文件不存在 / 写错 / 损坏，Must 会直接让程序崩溃报错（panic）
如果加载成功，就正常返回模板对象
一次性加载 2 个模板文件，结果：生成一个模板集合对象（里面装着这两个页面）

这行代码的3 个好处（为什么要这么写）
	一次加载，全局使用
	程序启动时加载一次模板，后面所有页面渲染都用它，性能极高
	同时加载多个模板
	登录页、后台页都在一个集合里，渲染时指定文件名即可
	自动错误检查
	模板写错了，启动程序就会报错，不用等到运行时才发现

	后续怎么用？（对应你之前的渲染代码）
	加载完成后，渲染页面直接用 templates.ExecuteTemplate：
	// 渲染登录页
	templates.ExecuteTemplate(w, "login.html", data)
	// 渲染后台页
	templates.ExecuteTemplate(w, "dashboard.html", data)
*/
var templates = template.Must(template.ParseFiles("templates/login.html", "templates/dashboard.html"))

// 模拟用户存储
var validUsers = map[string]string{
	"admin": "123456",
	"user":  "password",
}

// 辅助函数：检查是否登录
func isLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}
	_, ok := validUsers[cookie.Value]
	return ok
}

// 辅助函数：获取当前登录用户名
func getUsername(r *http.Request) string {
	cookie, _ := r.Cookie("session")
	return cookie.Value
}

// 登录页面
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// 如果已经登录，直接跳转到仪表盘
	if isLoggedIn(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	// 显示登录表单
	data := struct {
		Error    string
		Username string
	}{}
	templates.ExecuteTemplate(w, "login.html", data)
}

// 处理登录提交
func loginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析表单
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// 验证
	errorMsg := ""
	if username == "" {
		errorMsg = "用户名不能为空"
	} else if password == "" {
		errorMsg = "密码不能为空"
	} else {
		expectedPwd, ok := validUsers[username]
		if !ok || expectedPwd != password {
			errorMsg = "用户名或密码错误"
		}
	}

	// 如果验证失败，回到登录页并显示错误
	if errorMsg != "" {
		data := struct {
			Error    string
			Username string
		}{
			Error:    errorMsg,
			Username: username,
		}
		templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	// 验证成功：设置登录状态（这里用 cookie 模拟）
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: username,
		Path:  "/",
	})

	// 重定向到仪表盘
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// 仪表盘页面
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// 检查是否登录
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := getUsername(r)
	data := struct {
		Username string
	}{Username: username}
	templates.ExecuteTemplate(w, "dashboard.html", data)
}

// 退出登录
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// 清除 cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
