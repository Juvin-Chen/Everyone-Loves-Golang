# Go html/template 

------

## 一、核心概念

1. **模板**：HTML 页面 + 占位符，用于后端渲染动态数据
2. **作用**：后端把数据传给 HTML，生成最终页面
3. **核心规则**：**Go 结构体字段名 = 模板中 `.字段名`**（大小写 / 拼写必须完全一致）

------

## 二、最常用核心函数

这是项目中**100% 会用到**的 4 个函数，按使用顺序记：

### 1. `template.ParseFiles()`

- **作用**：**加载模板文件**，支持**同时加载多个**

- 语法：
  ```
  // 加载单个
  template.ParseFiles("templates/index.html")
  // 加载多个
  template.ParseFiles("templates/login.html", "templates/dashboard.html")
  ```
- **返回值**：模板对象 + 错误

### 2. `template.Must()`

- **作用**：**简化模板错误处理**，加载失败直接报错终止程序，如果模板文件不存在 / 写错 / 损坏，Must 会直接让程序**崩溃报错（panic）**

- 语法：

  ```
  template.Must(模板加载函数) 
  var templates = template.Must(template.ParseFiles(...))
  ```

- **场景**：全局加载模板时使用，不用写 `if err != nil`

### 3. `Execute()`

- **作用**：渲染**单个模板**，直接输出页面

- 语法：
  ```
  tpl.Execute(w, 数据)
  ```

  

### 4. `ExecuteTemplate()`

- **作用**：渲染**多个模板中的某一个**（配合多模板加载使用）

- 语法：
  ```
  // 指定模板文件名渲染
  tpl.ExecuteTemplate(w, "login.html", 数据)
  ```

  

------

## 三、全局模板定义（我的一个项目的写法）

**标准全局模板写法**：

```
// 定义全局变量，程序启动时一次性加载所有模板
var templates = template.Must(template.ParseFiles(
	"templates/login.html",
	"templates/dashboard.html",
	"templates/index.html",
))
```

优点：

- 启动时加载一次，后续直接用，性能高
- 多模板共用，随时切换渲染
- 自动检查模板错误

------

# 四、模板基础语法（超详细版 + 完整前后端示例）

结合你实战的**留言板、用户列表**场景，每个语法都配：**用法说明 + Go 后端代码 + 模板代码 + 渲染效果**，一看就懂！

------

## 1. 变量渲染

模板通过 `.` 符号获取后端传递的数据，分两种最常用场景：

### 场景 1：获取对象 / 结构体的指定字段（90% 业务用法）

- 用法：`{{.字段名}}`
- 规则：Go 结构体字段必须**大写**，模板名称严格对应


```
// 👉 Go 后端代码（定义数据）
type User struct {
	Name string // 大写导出字段
}
data := User{Name: "张三"}
```

```
<!-- 👉 模板代码 -->
<!-- 取数据中的 Name 字段 -->
<h1>你好，{{.Name}}</h1>
```


```
<!-- 👉 最终渲染到页面的效果 -->
<h1>你好，张三</h1>
```

### 场景 2：获取整个数据对象（仅支持简单数据：字符串 / 数字）

- 用法：`{{.}}`
- 禁忌：**不能直接打印结构体**，仅用于简单数据



```
// 👉 Go 后端代码（直接传字符串）
data := "欢迎来到留言板"
```

```
<!-- 👉 模板代码 -->
<!-- 取整个数据对象（简单数据专用） -->
<h1>{{.}}</h1>
```


```
<!-- 👉 最终渲染效果 -->
<h1>欢迎来到留言板</h1>
```

------

## 2. 条件判断 `if`

用于逻辑判断，搭配 Go 模板**内置比较函数**使用

常用内置函数：`gt`(大于)、`lt`(小于)、`eq`(等于)


```
// 👉 Go 后端代码
type User struct {
	Age int
}
data := User{Age: 20}
```

```
<!-- 👉 模板代码（带完整注释） -->
<!-- gt = 内置函数，含义：大于 → 判断年龄是否大于18 -->
{{if gt .Age 18}}
    <span>成年</span>
{{else}}
    <span>未成年</span>
{{end}}
```

```
<!-- 👉 最终渲染效果 -->
<span>成年</span>
```

------

## 3. 循环遍历 `range`

用于遍历**数组 / 切片**（你的留言列表、用户列表核心用法）

- `{{range 列表}}`：开始循环
- 循环内 `.`：代表**当前遍历的单条数据**
- `{{else}}`：列表为空时显示

### 示例：遍历留言列表（结构体格式，带 Content 字段）


```
// 👉 Go 后端代码
type Message struct {
	Content string // 留言内容
}
// 模拟留言数据
messages := []Message{
	{Content: "第一条留言"},
	{Content: "学习Go模板"},
}
// 封装数据
data := struct {
	Messages []Message
}{Messages: messages}
```

```html
<!-- 👉 模板代码（带完整注释） -->
<!-- 遍历留言列表 Messages -->
{{range .Messages}}
    <!-- . = 当前单条留言，.Content 获取留言内容 -->
    <p>{{.Content}}</p>
{{else}}
    <!-- 没有留言时显示 -->
    <p>暂无留言</p>
{{end}}
```

```html
<!-- 👉 最终渲染效果 -->
<p>第一条留言</p>
<p>学习Go模板</p>
```

### 补充：空列表效果

```html
// 后端传空切片
messages := []Message{}
```

```html
<!-- 页面渲染结果 -->
<p>暂无留言</p>
```

------

## 五、自定义函数（进阶）

### 1. 定义函数
```
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}
```

### 2. 注册函数

```
template.New("").Funcs(template.FuncMap{
    "formatDate": formatDate,
})
```

### 3. 模板使用
```
{{formatDate .Now}}
```

------

## 六、完整使用流程

```
package main

import (
	"html/template"
	"net/http"
)

// 1. 全局加载所有模板
var templates = template.Must(template.ParseFiles(
	"templates/login.html",
	"templates/index.html",
))

// 2. 页面处理器
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 3. 准备数据
	data := struct {
		Messages []string
	}{
		Messages: []string{"第一条留言", "第二条留言"},
	}

	// 4. 渲染指定模板
	templates.ExecuteTemplate(w, "index.html", data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
```

------

## 七、新手必避坑

1. **字段大小写**：Go 结构体必须**大写**，模板才能访问
2. **名称对应**：Go 字段名 = 模板 `.字段名`，错一个字母都不行
3. **文件路径**：`ParseFiles` 路径是**程序运行根目录**相对路径
4. **多模板**：必须用 `ExecuteTemplate` 指定文件名渲染

------

## 八、极简总结

1. **加载模板**：`ParseFiles`（支持多个）
2. **错误处理**：`Must` 一键简化
3. **渲染页面**：`ExecuteTemplate`（多模板专用）
4. **页面语法**：变量、`if`、`range` 三大件够用
5. **全局模板**：启动加载一次，全程复用