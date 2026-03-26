/*
lesson 8
模板渲染（html/template）
可参考-Standard-Library目录下html/template
*/

/*
1. 什么是模板引擎？
模板引擎是一种将 静态模板文件 与 动态数据 结合，生成最终 HTML 文本的工具。
例如，你有一个 HTML 文件，里面有一些占位符 {{.Name}}，模板引擎会把 Name 的实际值替换进去，生成完整的 HTML 返回给浏览器。

Go 标准库提供了两个模板包：
text/template：用于生成任意文本（如邮件、配置文件）
html/template：专门用于生成 HTML，会自动对输出进行转义，防止 XSS 攻击

在 Web 开发中，我们通常使用 html/template。
*/

package main

import (
	"html/template"
	"net/http"
	"time"
)

// 第一个模板示例
// 定义一个字符串模板，填充数据后输出到浏览器
func testTemplateString() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 定义模板字符串
		const tmplStr = `<h1>Hello,{{.Name}}!<h1>`

		// 创建模板对象并解析
		// template.New("hello") 创建一个名为 "hello" 的新模板对象。
		// Parse 方法解析模板字符串，返回一个模板对象。
		tmpl, err := template.New("hello").Parse(tmplStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// 准备数据
		data := struct {
			Name string
		}{
			Name: "zhangsan",
		}

		// 执行模板，将结果写入 ResponseWriter
		// tmpl.Execute(w, data) 将模板与数据合并，结果写入 w。
		// 模板中的 {{.Name}} 表示访问数据中的 Name 字段。
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.ListenAndServe(":8080", nil)
}

// 以上示例1 直接将模板定义在函数内部

// 实际项目中，模板通常放在单独的文件中，便于维护。

// 从文件加载模板
func testTmplInHtml() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 解析模板文件
		/*
			ParseFiles 可以同时加载多个模板文件，但这里我们只加载一个。
			执行时，模板文件的路径是相对于程序运行时的当前目录（通常是项目根目录），所以要确保 templates/index.html 存在。
		*/
		tmpl, err := template.ParseFiles("08-index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Name string
			Age  int
		}{
			Name: "lisi",
			Age:  25,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	http.ListenAndServe(":8080", nil)
}

/*
4. 模板语法基础
4.1 变量
{{.Field}}：访问数据对象的字段
{{.}}：访问整个数据对象

4.2 条件判断
html
{{if .IsAdmin}}
    <p>您是管理员</p>
{{else}}
    <p>普通用户</p>
{{end}}

4.3 循环遍历
html
<ul>
{{range .Items}}
    <li>{{.}}</li>
{{end}}
</ul>

4.4 自定义函数
可以注册自定义函数到模板，例如格式化日期、截断字符串等。
极简示例：
	// 1. 注册函数：把年龄转成 "成年/未成年"
	tmpl = tmpl.Funcs(template.FuncMap{
		"checkAge": func(age int) string {
			if age >= 18 {
				return "成年"
			}
			return "未成年"
		},
	})
	html
	预览
	<!-- 2. 模板里用 -->
	<td>{{checkAge .Age}}</td>

5. 安全性：自动转义
html/template 会自动对输出进行转义，防止 XSS 攻击。
例如，如果你传递 <script>alert('xss')</script> 给模板，它会被转义为 &lt;script&gt;alert('xss')&lt;/script&gt;，浏览器不会执行脚本。
如果确实需要输出原始 HTML（比如从后台编辑的富文本），可以使用 {{.Content | safeHTML}}，但必须确保内容是可信的。
*/

// 对应 user.html
type User struct {
	ID   int
	Name string
	Age  int
}

func testUserHtml() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("08-user.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users := []User{
			{ID: 1, Name: "张三", Age: 20},
			{ID: 2, Name: "李四", Age: 25},
		}
		data := struct {
			Users []User
		}{
			Users: users,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// test 部分
/*
要求：
修改上面的用户列表，增加一个条件：如果年龄大于 18，在表格中显示“成年”，否则显示“未成年”（提示：在模板中使用 if 判断）。
尝试在模板中使用自定义函数：定义一个 formatDate 函数，将 time.Time 格式化为 2006-01-02，并在模板中调用它（提示：使用 template.FuncMap）。
*/

// 自定义函数 - 格式化时间为 2006-01-02
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// 全局模板：注册自定义函数 + 加载所有模板文件
var tpl = template.Must(
	template.New("").
		// 注册自定义函数到模板
		Funcs(template.FuncMap{
			"formatDate": formatDate,
		}).
		// 同时加载多个模板（核心知识点）
		ParseFiles("templates/greet.html", "templates/user.html"),
)

func testTestHtml() {

}
