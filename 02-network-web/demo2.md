# Go Web 编程：处理 HTTP 请求 (Handler)

**核心要点：** `http.Handle` 函数 / `http.HandleFunc` 函数

---

## 1. 创建 Web Server

使用 `http.ListenAndServe(addr, handler)` 启动服务：

* **第一个参数 (addr):** 网络地址。
  * 如果为 `""`，则默认监听所有网络接口的 `80` 端口。
* **第二个参数 (handler):** 处理请求的处理器。
  * 如果为 `nil`，那么就会使用默认的 `DefaultServeMux`。
  * `DefaultServeMux` 是一个 multiplexer（多路复用器/多路路由器），它本质上也是一个 `Handler`，实现了接口的 `ServeHTTP()` 方法。

---

## 2. http.ListenAndServe() 源代码解析

### (1) 顶层 ListenAndServe() 函数源码

该函数监听指定的 TCP 网络地址 `addr`，然后调用 `Serve` 方法处理传入的连接。

```go
// 参数说明：
// - addr: 监听地址（如 ":8080"）
// - handler: 处理请求的处理器，nil 表示使用默认的 DefaultServeMux（默认路由器）
func ListenAndServe(addr string, handler Handler) error {
    // 1. 创建 http.Server 结构体实例（这是 Web 服务的核心配置）
    server := &Server{Addr: addr, Handler: handler}
    // 2. 调用 Server 的 ListenAndServe 方法启动服务
    return server.ListenAndServe()
}
```

### (2) http.Server 结构体

`http.Server` 是一个 `struct`：

- `Addr` 字段表示网络地址（为 `""` 则是所有网络接口的 80 端口）。
- `Handler` 字段（如果为 `nil`，就是 `DefaultServeMux`）。

**总结：** 顶层函数本质上是调用了 `Server` 的 `ListenAndServe()` 函数。自己实例化 `http.Server` 结构体比直接调用 `http.ListenAndServe()` 更灵活。

------

## 3. HTTPS 支持 (SSL/TLS)

如果需要加 SSL/TLS 启动 HTTPS 服务，使用的是 `http.ListenAndServeTLS()` 或 `server.ListenAndServeTLS()`。

Go

```
http.ListenAndServeTLS(addr, certFile, keyFile, handler)
```

它比普通的启动多了两个参数：**公钥证书文件 (certFile)** 和 **私钥文件 (keyFile)**。

------

## 4. Handler 到底是什么？

`Handler` 在 Go 中是一个**接口（interface）**，它只定义了一个方法 `ServeHTTP()`。

Go

```
type Handler interface {
    // http.ResponseWriter 是一个接口（interface）
    // *http.Request 是一个结构体指针（struct pointer）
    ServeHTTP(ResponseWriter, *Request)
}
```

------

## 5. Q & A：多个 Handler 与路由的本质

### 核心问题：`myHandler` 算不算“路由”？

**严格来说，代码里写的 `myHandler` 不叫“路由”（Router），它只是一个纯粹的“处理器”（Handler）。**

以下是普通 Handler 和 路由器 (ServeMux) 的流转对比图：

代码段

```
graph LR
    %% 普通 Handler 流程
    subgraph 单一处理器
    req1[HTTP 请求] --> mh[MyHandler]
    end

    %% ServeMux 路由流程
    subgraph 多路路由器 
    req2[HTTP 请求] --> mux[DefaultServeMux]
    mux --> h1[Handler 1]
    mux --> h2[Handler 2]
    mux -.-> hn[...]
    mux --> h3[Handler 3]
    end
    
    style mh fill:#d9534f,stroke:#d43f3a,color:#fff
    style mux fill:#d9534f,stroke:#d43f3a,color:#fff
    style h1 fill:#f0ad4e,stroke:#eea236,color:#fff
    style h2 fill:#f0ad4e,stroke:#eea236,color:#fff
    style hn fill:#f0ad4e,stroke:#eea236,color:#fff
    style h3 fill:#f0ad4e,stroke:#eea236,color:#fff
```

- **上半部分（单一处理器）：** 对应代码中的 `demo2_2()`。所有的请求（不管访问 `/login` 还是 `/pay`），全部一股脑交给了 `myHandler`。它没有做任何“分发”工作，所有请求都返回 "Hello web"。**这不叫路由。**
- **下半部分（多路路由器）：** `DefaultServeMux` 才是真正的路由（多路复用器 Multiplexer）。它的工作是“看人下菜碟”：如果请求是 `/a`，它就交给 `Handler 1`；如果是 `/b`，就交给 `Handler 2`。

> **💡 总结：** 路由（`ServeMux`）的本质，其实也是一个特殊的 `Handler`（因为它也实现了 `ServeHTTP` 方法）。只不过它的 `ServeHTTP` 内部逻辑是：**根据 URL 路径，把请求转发给其他具体的 Handler。**

------

## 6. 完整代码演示

```go
package main

import "net/http"

type myHandler struct{}

// 自己定义一个handler，实现 ServeHTTP 方法
func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello web"))
}

func demo2() {
    // 实例化 http.Server 结构体
    server := http.Server{
        Addr:    "localhost:8080",
        Handler: nil, // 传入 nil，默认使用 DefaultServeMux
    }
    server.ListenAndServe()
    
    // 上面的写法更灵活一些，与下面直接调用的写法等价：
    // http.ListenAndServe("localhost:8080", nil)
}

func demo2_2() {
    mh := myHandler{}
    
    // 这种自定义处理器的方式，不同路径都会走到这个 Handler，都会把 hello web 输出
    server := http.Server{
        Addr:    "localhost:8080",
        Handler: &mh, // 传入自定义的 Handler
    }
    server.ListenAndServe()
}
```