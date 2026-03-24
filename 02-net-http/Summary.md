# net/http 标准库 Summary（00-07）

## 目录

1. [整体认知与运行入口](#1-整体认知与运行入口)
2. [Handler、ServeMux 与路由分发（01）](#2-handlerservemux-与路由分发01)
3. [HTTP 响应写回（02）](#3-http-响应写回02)
4. [请求行信息解析（03）](#4-请求行信息解析03)
5. [请求体处理：表单与 JSON（04）](#5-请求体处理表单与-json04)
6. [请求头读取与应用（05）](#6-请求头读取与应用05)
7. [中间件与链式调用（06）](#7-中间件与链式调用06)
8. [Context：取消、超时、请求级数据（07）](#8-context取消超时请求级数据07)
9. [一页 API 速查表](#9-一页-api-速查表)
10. [常见坑与最佳实践](#10-常见坑与最佳实践)

---

## 1. 整体认知与运行入口

`00-测试.go` 是练习入口，通过切换 `main` 中调用的测试函数来验证某个知识点。

推荐方式：每次只跑一条主线，避免多个 `ListenAndServe` 同时占端口。

典型入口：

```go
func main() {
	testHttpServer()      // 01
	// testResponseWriter() // 02
}
```

---

## 2. Handler、ServeMux 与路由分发（01）

核心结论：

- `http.Handler` 是处理请求的抽象（执行者）。
- `http.ServeMux` 是路由分发器（调度员）。
- `http.HandleFunc` 本质是给 `http.DefaultServeMux` 注册路由。
- `http.ListenAndServe(addr, nil)` 的 `nil` 表示使用默认路由器。

关键调用：

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
})
http.ListenAndServe(":8080", nil)
```

自定义路由器：

```go
mux := http.NewServeMux()
mux.HandleFunc("/about", aboutHandler)
http.ListenAndServe(":8080", mux)
```

自定义 Handler：

```go
type myHandler struct{}

func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello web"))
}
```


### 【等价关系（高频面试点）】

下面两种写法 **100% 等价**：

1. 简化写法

```
http.ListenAndServe(addr, handler)
```

2. 底层真实写法（手动创建 Server 对象）

```
(&http.Server{Addr: addr, Handler: handler}).ListenAndServe()
```

---

## 3. HTTP 响应写回（02）

响应结构：状态行 + 响应头 + 响应体。

核心规则：

- 先设置头和状态码，再写响应体。
- 若不显式 `WriteHeader`，首次 `Write/Fprintf` 时默认发送 `200 OK`。
- `w.Header()` 返回 `http.Header`，底层是 `map[string][]string`。

关键调用：

```go
w.Header().Set("Content-Type", "text/plain; charset=utf-8")
w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
w.WriteHeader(http.StatusOK)
fmt.Fprintf(w, "Hello, World!")
```

`Set` 与 `Add`：

- `Set(k, v)`：覆盖该 key 的已有值。
- `Add(k, v)`：追加值，适合多值头（如多个 `Set-Cookie`）。

---

## 4. 请求行信息解析（03）

本节抓三件事：方法、路径、查询参数。

关键字段：

- `r.Method`
- `r.URL.Path`
- `r.URL.Query()`（返回 `url.Values`）

关键调用：

```go
method := r.Method
path := r.URL.Path
query := r.URL.Query()
name := query.Get("name")      // 取第一个值
hobbies := query["hobby"]      // 取全部值 []string
```

注意点：

- 查询参数会自动 URL 解码。
- 一个 key 可能对应多个值，所以是 `[]string`。
- 不建议一个处理器内部靠 `switch path` 手搓所有路由，优先分开注册路由。

---

## 5. 请求体处理：表单与 JSON（04）

### 5.1 POST 表单（`application/x-www-form-urlencoded`）

核心流程：

1. 校验方法（通常只接受 POST）。
2. `r.ParseForm()` 显式解析并处理错误。
3. 用 `FormValue/PostFormValue` 取值。

关键调用：

```go
if r.Method != http.MethodPost {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
if err := r.ParseForm(); err != nil {
	w.WriteHeader(http.StatusBadRequest)
	return
}
username := r.PostFormValue("username")
password := r.PostFormValue("password")
```

`FormValue` vs `PostFormValue`：

- `FormValue`：会看 URL 参数 + Body 表单参数。
- `PostFormValue`：只看 Body 表单参数，更不易被 URL 参数干扰。

### 5.2 JSON 请求体（`application/json`）

关键调用：

```go
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var req LoginRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	w.WriteHeader(http.StatusBadRequest)
	return
}
defer r.Body.Close()
```

建议补上：校验 `Content-Type` 是否为 `application/json`，错误时返回 `415`。

---

## 6. 请求头读取与应用（05）

核心目标：读懂客户端“附加信息”。

关键调用：

```go
for key, values := range r.Header {
	fmt.Fprintf(w, "%s: %v\n", key, values)
}

ua := r.Header.Get("User-Agent")
acceptAll := r.Header.Values("Accept")
host := r.Host
length := r.ContentLength
```

实践点：

- 可基于 `User-Agent` 做移动/桌面差异响应。
- 可从 `X-Forwarded-For`、`X-Real-IP` 推断客户端 IP。
- 任何请求头都可被伪造，安全判断要结合可信代理和鉴权机制。

---

## 7. 中间件与链式调用（06）

定义：中间件是

```go
func(http.Handler) http.Handler
```

日志中间件模板：

```go
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s cost=%v", r.URL.Path, time.Since(start))
	})
}
```

认证中间件模板：

```go
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

链式调用建议：从后往前包装，保证“传入顺序 = 执行顺序”。

```go
func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
```

执行模型（洋葱模型）：前置 A -> 前置 B -> 业务 -> 后置 B -> 后置 A。

---

## 8. Context：取消、超时、请求级数据（07）

Context 的两大能力：

- 控制流程：取消、超时、截止时间。
- 传递数据：请求范围内的轻量数据（如 requestID）。

### 8.1 手动取消

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
go longTask(ctx)
cancel()
```

### 8.2 超时取消

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
```

### 8.3 WithValue 传递请求级数据

```go
type contextKey string
const requestIDKey contextKey = "requestID"

ctx := context.WithValue(r.Context(), requestIDKey, requestID)
next.ServeHTTP(w, r.WithContext(ctx))
```

最佳实践：

- `key` 用自定义类型，避免冲突。
- 只存请求范围、轻量、不可变数据。
- 派生出 cancel/timeout 的 context 必须 `defer cancel()` 释放资源。

---

## 9. 一页 API 速查表

服务启动与路由：

- `http.HandleFunc(pattern, handler)`
- `mux := http.NewServeMux()`
- `mux.HandleFunc(pattern, handler)`
- `http.ListenAndServe(addr, handlerOrNil)`
- `server := http.Server{Addr, Handler}; server.ListenAndServe()`

请求读取：

- `r.Method`
- `r.URL.Path`
- `r.URL.Query().Get("k")`
- `r.ParseForm()`
- `r.FormValue("k")`
- `r.PostFormValue("k")`
- `json.NewDecoder(r.Body).Decode(&v)`
- `r.Header.Get("User-Agent")`
- `r.Context()`

响应写回：

- `w.Header().Set(k, v)`
- `w.WriteHeader(code)`
- `w.Write([]byte("..."))`
- `fmt.Fprintf(w, "...", args...)`

---

## 10. 常见坑与最佳实践

常见坑：

1. 写响应体后再改响应头，修改不会生效。
2. 忘记 `return`，导致错误响应后继续执行正常逻辑。
3. 中间件链顺序写反，导致执行顺序与预期相反。
4. `WithTimeout/WithCancel` 后忘记 `cancel()`，有资源泄露风险。
5. 用 `FormValue` 读取登录参数时被 URL 参数污染。

最佳实践：

1. 显式处理方法与输入校验：`405/400/415` 分层返回。
2. 在中间件统一做日志、鉴权、恢复 panic。
3. 在响应中统一 `Content-Type` 与错误格式。
4. 对 JSON API 统一定义请求/响应结构体。
5. 生产环境的客户端 IP 获取要基于可信代理配置。

