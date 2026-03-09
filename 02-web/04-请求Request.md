# Go Web 编程：HTTP Request 深度解析

------

## 一、HTTP 消息基础

### 1. HTTP 消息结构

HTTP 请求（Request）和响应（Response）遵循完全相同的四层结构：

- **请求 / 响应行**：第一行，包含方法、路径、协议版本（请求）或状态码、状态信息（响应）
- **0 个或多个 Header**：键值对，描述请求 / 响应的元数据（如内容类型、编码等）
- **空行**：分隔 Header 和 Body，是协议规定的分隔符
- **可选的消息体（Body）**：承载请求 / 响应的实际数据（如表单、JSON、文件等）

**示例（HTTP 请求）**：

plaintext











```
GET /Protocols/rfc2616/rfc2616.html HTTP/1.1
Host: www.w3.org
User-Agent: Mozilla/5.0

(空行)
```

### 2. net/http 包的抽象

Go 的 `net/http` 包提供了 `Request` 和 `Response` 两个核心结构体，用于在代码中**结构化地表示和操作 HTTP 消息**，无需手动解析原始文本。

------

## 二、`http.Request` 结构体

### 1. 定义

`http.Request` 是一个结构体，代表客户端发送到服务器的 HTTP 请求（也可用于客户端主动发起请求）。它封装了 HTTP 消息的所有部分。

### 2. 核心字段（来自标准库 `net/http/request.go`）

go



运行









```
// 简化版 Request 结构体定义
type Request struct {
    // 请求方法，如 "GET", "POST", "PUT"
    Method string
    // 请求的 URL，类型为 *url.URL
    URL *url.URL
    // 协议版本，如 "HTTP/1.1"
    Proto string
    // 请求头，类型为 http.Header（本质是 map[string][]string）
    Header Header
    // 请求体，实现了 io.ReadCloser 接口
    Body io.ReadCloser
    // 表单数据（解析后）
    Form url.Values
    PostForm url.Values
    MultipartForm *multipart.Form
    // ... 其他辅助字段
}
```

### 3. 常用便捷方法

`Request` 提供了封装好的方法，方便直接访问常用信息：

- `Cookie(name string) (*Cookie, error)`：获取指定名称的 Cookie
- `UserAgent() string`：获取请求的 User-Agent 头
- `Referer() string`：获取请求的 Referer 头
- `Host() string`：获取请求的 Host 头

------

## 三、`Request.URL` 字段

### 1. 定义

`Request.URL` 是一个指向 `url.URL` 结构体的指针，它代表了 HTTP 请求行中的 URL 部分。

### 2. `url.URL` 结构体（来自标准库 `net/url/url.go`）

go



运行









```
type URL struct {
    // URL 协议，如 "http", "https", "ftp"
    Scheme string
    // 不透明的 URL 部分（极少使用）
    Opaque string
    // 认证信息（用户名和密码）
    User *Userinfo
    // 主机名和端口，如 "localhost:8080", "www.example.com"
    Host string
    // URL 路径，如 "/home", "/api/users"
    Path string
    // 原始查询字符串，如 "id=123&name=foo"
    RawQuery string
    // URL 片段（# 后面的部分）
    Fragment string
}
```

### 3. 关键子字段解析

#### (1) `RawQuery`

- 表示 URL 中 `?` 后面的查询字符串，是**原始字符串形式**。
- 例子：URL `http://example.com/post?id=123&thread_id=456` 中，`RawQuery` 的值是 `"id=123&thread_id=456"`。
- 可以通过 `r.URL.Query()` 方法将其解析为 `url.Values`（键值对），方便读取。

#### (2) `Fragment`

- 表示 URL 中 `#` 后面的片段，用于定位文档内的位置（如页面锚点）。
- **重要特性**：浏览器在发送 HTTP 请求时，会**自动去掉 Fragment 部分**，因此服务器端无法从浏览器的请求中获取到 Fragment 值。但非浏览器客户端（如 Go 的 `http.Client`）可以发送包含 Fragment 的请求。

------

## 四、`Request.Header` 字段

### 1. 定义

`Request.Header` 的类型是 `http.Header`，它本质上是一个 `map[string][]string`，用于存储 HTTP 请求头的键值对。

### 2. `http.Header` 类型定义（来自标准库 `net/http/header.go`）

go



运行









```
type Header map[string][]string

// 常用操作方法
func (h Header) Get(key string) string      // 获取值（字符串形式）
func (h Header) Set(key, value string)     // 设置值（覆盖）
func (h Header) Add(key, value string)     // 添加值（追加）
func (h Header) Del(key string)            // 删除键
```

### 3. 核心特性

- **一个键对应多个值**：HTTP 标准允许同一个 Header 键出现多次，因此值的类型是 `[]string`。
- **键名不区分大小写**：`Accept-Encoding` 和 `accept-encoding` 会被视为同一个键。

### 4. 访问 Header 的两种方式

假设请求头包含 `Accept-Encoding: gzip, deflate`：

#### (1) 直接访问 map

go



运行









```
// 返回 []string{"gzip", "deflate"}，获取所有值
values := r.Header["Accept-Encoding"]
```

#### (2) 使用 `Get` 方法（推荐）

go



运行









```
// 返回字符串 "gzip, deflate"（用逗号拼接所有值）
value := r.Header.Get("Accept-Encoding")
```

------

## 五、`Request.Body` 字段

### 1. 定义

`Request.Body` 是一个实现了 `io.ReadCloser` 接口的对象，用于读取 HTTP 请求体的数据流。

### 2. `io.ReadCloser` 接口（来自标准库 `io/io.go`）

go



运行









```
type ReadCloser interface {
    // Reader 接口：从流中读取数据到字节数组
    Read(p []byte) (n int, err error)
    // Closer 接口：关闭流，释放资源
    Close() error
}
```

### 3. 核心特性

- **请求体是流**：Body 不是一次性加载到内存的字节数组，而是一个数据流，需要通过 `Read` 方法逐步读取。
- **必须关闭**：使用完 Body 后必须调用 `Close()` 方法，否则会造成资源泄漏（推荐用 `defer` 确保关闭）。
- **只能读取一次**：Body 流只能被读取一次，读取后无法再次读取（除非先保存到内存）。

### 4. 读取 Body 的示例

go



运行









```
func handler(w http.ResponseWriter, r *http.Request) {
    // 确保 Body 被关闭，避免资源泄漏
    defer r.Body.Close()

    // 读取 Body 所有数据到内存
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "读取请求体失败", http.StatusInternalServerError)
        return
    }

    fmt.Printf("请求体内容：%s\n", string(body))
}
```

------

## 六、完整示例：解析 HTTP 请求

go



运行









```
package main

import (
    "fmt"
    "io"
    "net/http"
)

func demoHandler(w http.ResponseWriter, r *http.Request) {
    // 1. 解析 URL 信息
    fmt.Println("=== URL 信息 ===")
    fmt.Println("路径:", r.URL.Path)
    fmt.Println("原始查询字符串:", r.URL.RawQuery)
    fmt.Println("主机:", r.URL.Host)
    fmt.Println("协议:", r.URL.Scheme)

    // 2. 解析 Header 信息
    fmt.Println("\n=== Header 信息 ===")
    fmt.Println("User-Agent:", r.Header.Get("User-Agent"))
    fmt.Println("Accept-Encoding:", r.Header["Accept-Encoding"]) // 数组形式
    fmt.Println("Host:", r.Host())

    // 3. 解析 Body 信息
    fmt.Println("\n=== Body 信息 ===")
    defer r.Body.Close()
    body, _ := io.ReadAll(r.Body)
    fmt.Println("请求体:", string(body))
}

func main() {
    http.HandleFunc("/", demoHandler)
    fmt.Println("服务器启动在 :8080 端口")
    http.ListenAndServe(":8080", nil)
}
```

------

## 七、核心总结

1. **HTTP 消息结构**：行 → 头 → 空行 → 体，`net/http` 包用 `Request` 和 `Response` 封装了这一结构，让开发者无需手动解析。
2. **`Request.URL`**：封装了 URL 的各个部分，`RawQuery` 是原始查询字符串，`Fragment` 在浏览器请求中不可用。
3. **`Request.Header`**：本质是 `map[string][]string`，一个键对应多个值，`Get` 方法是便捷的访问方式。
4. **`Request.Body`**：是 `io.ReadCloser` 流，必须读取后关闭，且只能读取一次，避免资源泄漏。