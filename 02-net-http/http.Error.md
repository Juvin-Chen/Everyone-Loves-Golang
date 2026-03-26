# Go Web 错误处理：`http.Error` 完全指南

## 核心结论

`http.Error` 是 Go 通用的 HTTP 错误返回函数，**但不是所有错误都用 500 状态码**，必须根据错误类型返回**标准 HTTP 状态码**。

------

## 一、`http.Error` 函数详解

### 1. 作用

向客户端返回标准化的 HTTP 错误响应，是 Go Web 开发最基础的错误处理工具。

### 2. 函数参数

```
http.Error(
    w http.ResponseWriter,  // 固定：响应写入器
    errMsg string,           // 错误提示信息（给用户/前端看）
    statusCode int           // 核心：HTTP 状态码
)
```

### 3. 常用的默认写法（500 服务器错误）

```
// 仅用于：代码崩溃、模板加载失败、数据库连接失败等服务器内部错误
http.Error(w, err.Error(), http.StatusInternalServerError)
```

------

## 二、标准错误场景 + 状态码对照表

### 分类规则

- **4xx**：客户端错误（用户 / 前端的问题）
- **5xx**：服务器错误（后端代码的问题）

| 错误类型          | 状态码常量                       | 状态码 | 适用场景                    | 代码示例                                                     |
| :---------------- | :------------------------------- | :----- | :-------------------------- | :----------------------------------------------------------- |
| 客户端参数错误    | `http.StatusBadRequest`          | 400    | 用户传参错误、请求格式非法  | `http.Error(w, "请求参数错误", http.StatusBadRequest)`       |
| 页面 / 资源不存在 | `http.StatusNotFound`            | 404    | 访问了不存在的路由 / 文件   | `http.Error(w, "页面不存在", http.StatusNotFound)`           |
| 未登录授权        | `http.StatusUnauthorized`        | 401    | 需要登录才能访问的接口      | `http.Error(w, "请先登录", http.StatusUnauthorized)`         |
| 无访问权限        | `http.StatusForbidden`           | 403    | 已登录但无权限操作          | `http.Error(w, "没有访问权限", http.StatusForbidden)`        |
| 服务器内部错误    | `http.StatusInternalServerError` | 500    | 代码报错、模板 / 数据库失败 | `http.Error(w, "服务器异常", http.StatusInternalServerError)` |

------

## 三、完整可运行代码示例

```
package main

import (
	"html/template"
	"net/http"
)

// 首页处理器（包含错误处理）
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 场景1：加载模板失败 → 500 服务器错误
	tpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "模板加载失败："+err.Error(), http.StatusInternalServerError)
		return
	}

	// 渲染模板
	err = tpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "页面渲染失败", http.StatusInternalServerError)
		return
	}
}

// 测试接口（演示多种错误处理）
func testHandler(w http.ResponseWriter, r *http.Request) {
	// 获取用户参数
	name := r.URL.Query().Get("name")

	// 场景2：参数为空 → 400 客户端错误
	if name == "" {
		http.Error(w, "参数 name 不能为空", http.StatusBadRequest)
		return
	}

	// 正常响应
	w.Write([]byte("你好：" + name))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/test", testHandler)

	// 启动服务
	http.ListenAndServe(":8080", nil)
}
```

------

## 四、关键总结

1. **`http.Error` 是通用工具**，所有 HTTP 错误都可以用它返回
2. **状态码不能乱写**：客户端问题用 `4xx`，服务器问题用 `5xx`
3. **禁止所有错误都返回 500**，不符合 Web 规范，前端无法正常处理
4. 最常用三件套：
   - 参数错误 → `400`
   - 页面不存在 → `404`
   - 代码报错 → `500`