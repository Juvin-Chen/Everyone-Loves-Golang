# Go Web 编程：标准库 Handler 工具深度解析

------

## 一、`http.NotFoundHandler`：404 处理器

### 1. 定义

```go
func NotFoundHandler() Handler
```

返回一个预定义的 `Handler`，对所有请求统一返回 **404 page not found** 响应。

### 2. 标准库源码（简化版）

```go
// 来自 net/http/server.go
func NotFoundHandler() Handler {
    return HandlerFunc(func(w ResponseWriter, r *Request) {
        Error(w, "404 page not found", StatusNotFound)
    })
}

// Error 函数：设置状态码并返回错误文本
func Error(w ResponseWriter, error string, code int) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(code)
    w.Write([]byte(error))
}
```

### 3. 核心逻辑

- 本质是一个 `HandlerFunc` 适配器，内部调用 `Error` 函数；
- 自动设置 `Content-Type` 为纯文本，状态码为 `404`。

### 4. 代码示例

```go
package main

import "net/http"

func main() {
    // 所有访问 /404 的请求，都返回 404 页面
    http.Handle("/404", http.NotFoundHandler())
    http.ListenAndServe(":8080", nil)
}
```

------

## 二、`http.RedirectHandler`：重定向处理器

### 1. 定义

```go
// 生产一个 “重定向处理器”
func RedirectHandler(url string, code int) Handler
```

返回一个 `Handler`，将所有请求重定向到指定 URL，并使用指定的 3xx 状态码。

#### 补充：常用 3xx 重定向状态码（3 个）

| 状态码常量                     | 数值 | 名称               | 通俗含义                                                     | 浏览器行为                                            |
| :----------------------------- | :--- | :----------------- | :----------------------------------------------------------- | :---------------------------------------------------- |
| `http.StatusMovedPermanently`  | 301  | 永久重定向         | 旧地址**永久作废**，浏览器会记住新地址，下次直接访问新地址   | 跳转 + 缓存新地址（后续访问旧地址不再发请求，直接跳） |
| `http.StatusFound`             | 302  | 临时重定向         | 旧地址**暂时可用新地址**，浏览器不缓存                       | 跳转，但下次访问旧地址仍会先问服务器                  |
| `http.StatusTemporaryRedirect` | 307  | 临时重定向（严格） | 和 302 类似，但要求浏览器**保持原请求方法**（比如 POST 仍用 POST） | 跳转，且不改变请求方法（302 可能被浏览器改成 GET）    |

> 301（永久）、302（临时）足够日常用，307 是进阶的 “严格临时重定向”。

### 2. 标准库源码（简化版）

```go
// 来自 net/http/server.go
func RedirectHandler(url string, code int) Handler {
    return HandlerFunc(func(w ResponseWriter, r *Request) {
        w.Header().Set("Location", url)
        w.WriteHeader(code)
    })
}
```

### 3. 核心逻辑

- 在响应头中设置 `Location` 字段为目标 URL；
- 返回 3xx 状态码（如 `301` 永久重定向、`302` 临时重定向），触发浏览器跳转。

### 4. 代码示例

```go
package main

import "net/http"

func main() {
    // 访问 /old 时，永久重定向到 /new
    http.Handle("/old", http.RedirectHandler("/new", http.StatusMovedPermanently))
    http.ListenAndServe(":8080", nil)
}
```

------

## 三、`http.StripPrefix`：URL 前缀处理器

### 1. 定义

```go
func StripPrefix(prefix string, h Handler) Handler
```

返回一个 `Handler`，先从请求 URL 中移除指定前缀，再将修改后的请求交给下一个 `Handler` 处理；如果 URL 不匹配前缀，返回 404。

### 2. 标准库源码（简化版）

```go
// 来自 net/http/server.go
func StripPrefix(prefix string, h Handler) Handler {
    if prefix == "" {
        return h
    }
    return HandlerFunc(func(w ResponseWriter, r *Request) {
        // 检查 URL 是否包含指定前缀
        if !strings.HasPrefix(r.URL.Path, prefix) {
            NotFoundHandler().ServeHTTP(w, r)
            return
        }
        // 复制请求对象，修改 URL 路径（移除前缀）
        r2 := new(Request)
        *r2 = *r
        r2.URL = new(url.URL)
        *r2.URL = *r.URL
        r2.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
        // 交给下一个 Handler 处理
        h.ServeHTTP(w, r2)
    })
}
```

### 3. 核心逻辑

- 类似 “中间件”：先修改请求，再转发给目标处理器；
- 常用于静态文件服务，避免暴露真实文件路径。

### 4. 代码示例（配合静态文件服务）

```go
package main

import "net/http"

func main() {
    // 访问 /static/index.html → 实际查找 wwwroot/index.html
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("wwwroot"))))
    http.ListenAndServe(":8080", nil)
}
```

------

## 四、`http.TimeoutHandler`：超时处理器

### 1. 定义

```go
func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler
```

返回一个 `Handler`，限制传入的 `h` 最多执行 `dt` 时间；超时后返回指定消息 `msg`。

### 2. 标准库源码（简化版）

```go
// 来自 net/http/server.go
func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler {
    return HandlerFunc(func(w ResponseWriter, r *Request) {
        // 创建带超时的上下文
        ctx, cancel := context.WithTimeout(r.Context(), dt)
        defer cancel()
        r = r.WithContext(ctx)

        // 包装 ResponseWriter，用于捕获响应头
        tw := &timeoutWriter{w: w}
        done := make(chan struct{})

        // 启动 goroutine 执行目标 Handler
        go func() {
            h.ServeHTTP(tw, r)
            close(done)
        }()

        // 等待完成或超时
        select {
        case <-done:
            // 正常完成，复制响应头
            copyHeader(w, tw.Header())
            w.WriteHeader(tw.status)
        case <-ctx.Done():
            // 超时，返回错误信息
            http.Error(w, msg, StatusServiceUnavailable)
        }
    })
}
```

### 3. 核心逻辑

- 基于 `context.WithTimeout` 实现超时控制；
- 超时后返回 `503 Service Unavailable` 状态码和指定消息。

### 4. 代码示例

```go
package main

import (
    "net/http"
    "time"
)

// 模拟一个慢处理 Handler
func slowHandler(w http.ResponseWriter, r *http.Request) {
    time.Sleep(10 * time.Second) // 模拟 10 秒耗时
    w.Write([]byte("处理完成"))
}

func main() {
    // 限制 slowHandler 最多执行 5 秒，超时返回“处理超时”
    http.Handle("/slow", http.TimeoutHandler(http.HandlerFunc(slowHandler), 5*time.Second, "处理超时"))
    http.ListenAndServe(":8080", nil)
}
```

------

## 五、`http.FileServer`：静态文件处理器

### 1. 定义

```go
func FileServer(root FileSystem) Handler
```

返回一个 `Handler`，基于指定的文件系统 `root` 提供静态文件服务。

### 2. 核心依赖：`FileSystem` 接口

```go
// 来自 net/http/fs.go
type FileSystem interface {
    Open(name string) (File, error) // 打开指定路径的文件
}

// http.Dir 实现了 FileSystem 接口，用于访问本地文件系统
type Dir string

func (d Dir) Open(name string) (File, error) {
    // 拼接本地文件路径并打开
    return os.Open(filepath.Join(string(d), name))
}
```

### 3. 标准库源码（简化版）

```go
// 来自 net/http/fs.go
func FileServer(root FileSystem) Handler {
    return &fileHandler{root: root}
}

type fileHandler struct {
    root FileSystem
}

func (f *fileHandler) ServeHTTP(w ResponseWriter, r *Request) {
    // 清理 URL 路径，防止路径遍历攻击
    path := cleanPath(r.URL.Path)
    // 打开文件
    file, err := f.root.Open(path)
    if err != nil {
        NotFoundHandler().ServeHTTP(w, r)
        return
    }
    defer file.Close()
    // 输出文件内容
    http.ServeContent(w, r, file.Name(), time.Time{}, file)
}
```

### 4. 代码示例

#### 示例 1：基础静态文件服务

```go
package main

import "net/http"

func main() {
    // 直接将 wwwroot 目录下的文件暴露到根路径
    http.ListenAndServe(":8080", http.FileServer(http.Dir("wwwroot")))
}
```

#### 示例 2：配合 `StripPrefix` 隐藏路径

```go
package main

import "net/http"

func main() {
    // 访问 /files/index.html → 实际返回 wwwroot/index.html
    http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("wwwroot"))))
    http.ListenAndServe(":8080", nil)
}
```

------

## 六、核心总结

| 函数              | 核心作用              | 典型场景                           |
| :---------------- | :-------------------- | :--------------------------------- |
| `NotFoundHandler` | 返回 404 响应         | 处理未匹配路由                     |
| `RedirectHandler` | 重定向到指定 URL      | 旧地址迁移、链接跳转               |
| `StripPrefix`     | 移除 URL 前缀         | 静态文件服务、路由美化             |
| `TimeoutHandler`  | 限制 Handler 执行时间 | 防止慢请求占用资源                 |
| `FileServer`      | 提供静态文件服务      | 托管 HTML、CSS、JS、图片等静态资源 |