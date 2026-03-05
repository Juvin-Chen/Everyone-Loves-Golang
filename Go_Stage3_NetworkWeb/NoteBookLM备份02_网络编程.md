# NoteBook02 网络编程

承接NoteBook备份01

## 第五章：Go 语言网络编程实战 (对标 Java)

在 Java 中，网络编程依赖 `java.net` 包下的 `InetAddress`、`URL`、`Socket` 等类。**在 Go 语言中，我们主要使用 `net` 和 `net/url` 标准库。**

> **关于 Socket 的深层理解：** Socket 是应用层和传输层之间的桥梁。通信必须有两端：`Socket(IP, Port, 协议)` 组成的三元组代表一个端点。 在服务端，`ServerSocket` 就像公司的**总机接线员**（只负责监听端口），当有连接进来时（`accept()`），分配一个新的 `Socket`（相当于**员工分机**）去和客户端专门通信。

### 5.1 解析 URL (对应 Java 的 `URL` 类)

```go
package main

import (
    "fmt"
    "net/url"
)

func main() {
    // 对应 Java 中的 new URL(...)
    rawUrl := "https://www.itbaizhan.com/search.html?kw=java"
    parsedUrl, err := url.Parse(rawUrl)
    if err != nil {
        panic(err)
    }

    fmt.Println("协议 (Protocol):", parsedUrl.Scheme)       // https
    fmt.Println("主机名 (Host):", parsedUrl.Host)           // www.itbaizhan.com
    fmt.Println("路径 (Path):", parsedUrl.Path)           // /search.html
    fmt.Println("参数部分 (Query):", parsedUrl.RawQuery)     // kw=java

    // 获取具体的参数值
    queryParams := parsedUrl.Query()
    fmt.Println("kw 的值:", queryParams.Get("kw"))        // java
}
```

### 5.2 TCP 编程：服务端与客户端 (对应 Java `ServerSocket` 与 `Socket`)

**Go 服务端实现 (Server)** 在 Go 中，我们不需要像 Java 那样通过手动分配多线程（`extends Thread`）来处理多客户端并发。Go 原生提供轻量级的 `goroutine` 进行高并发处理，极其简洁。

```go
package main

import (
    "bufio"
    "fmt"
    "net"
)

func main() {
    // 1. 对应 Java: ServerSocket serverSocket = new ServerSocket(8888);
    listener, err := net.Listen("tcp", ":8888")
    if err != nil {
        fmt.Println("监听失败:", err)
        return
    }
    defer listener.Close()
    fmt.Println("服务端启动，等待监听 8888 端口...")

    for {
        // 2. 对应 Java: Socket socket = serverSocket.accept(); (阻塞等待)
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("接收连接失败:", err)
            continue
        }
        fmt.Println("有客户端连接了:", conn.RemoteAddr())

        // 3. 启动 Goroutine 处理该客户端的读写（替代 Java 的多线程机制）
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)
    for {
        // 4. 对应 Java: br.readLine() 获取客户端消息
        msg, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("客户端断开连接:", conn.RemoteAddr())
            return
        }
        fmt.Print("客户端说: ", msg)

        // 回复客户端 (对应 Java: pw.println(str); pw.flush();)
        reply := fmt.Sprintf("服务器已收到: %s", msg)
        conn.Write([]byte(reply))
    }
}
```

**Go 客户端实现 (Client)**

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func main() {
    // 1. 对应 Java: Socket socket = new Socket("127.0.0.1", 8888);
    conn, err := net.Dial("tcp", "127.0.0.1:8888")
    if err != nil {
        fmt.Println("连接服务端失败:", err)
        return
    }
    defer conn.Close()
    fmt.Println("客户端启动，连接服务端成功！")

    inputReader := bufio.NewReader(os.Stdin)
    serverReader := bufio.NewReader(conn)

    for {
        fmt.Print("请输入发送内容 (exit退出): ")
        // 读取键盘输入
        msg, _ := inputReader.ReadString('\n')

        // 对应 Java: pw.println(msg); pw.flush();
        _, err = conn.Write([]byte(msg))
        if msg == "exit\n" || msg == "exit\r\n" {
            break
        }

        // 等待接收服务端回复
        reply, _ := serverReader.ReadString('\n')
        fmt.Print("服务端返回: ", reply)
    }
}
```

### 5.3 UDP 编程：基本数据类型与对象的传输 (对标 Java `DatagramSocket`)

UDP 不需要提前建连，在 Go 中使用 `net.ListenUDP` 和 `net.DialUDP`。

**传输自定义对象 (使用 JSON 序列化代替 Java 的 Serializable)** 在 Java 中，传对象必须要实现 `Serializable` 接口并使用 `ObjectOutputStream`。 在现代 Web 尤其是 Go 语言开发中，**传递结构体（对象）最标准、跨语言的做法是将其序列化为 JSON 字节数组。**

**UDP 客户端 (发送方)**

```go
package main

import (
    "encoding/json"
    "fmt"
    "net"
)

// 对应 Java 中的 Person 类，注意 Go 中的字段需要大写才能被 json 包导出
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // UDP 服务端地址
    serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
    // 本地客户端分配随机 UDP 端口
    conn, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        fmt.Println("建立UDP连接失败:", err)
        return
    }
    defer conn.Close()

    // 1. 实例化对象并进行 JSON 序列化
    p := Person{Name: "Oldlu", Age: 18}
    // json.Marshal 替代了 Java 的 ObjectOutputStream.writeObject()
    data, err := json.Marshal(p)
    if err != nil {
        fmt.Println("序列化失败:", err)
        return
    }

    // 2. 发送 UDP 数据报文 (对应 Java的 datagramSocket.send(dp))
    _, err = conn.Write(data)
    if err != nil {
        fmt.Println("发送失败:", err)
    } else {
        fmt.Println("对象发送成功!")
    }
}
```

**UDP 服务端 (接收方)**

```go
package main

import (
    "encoding/json"
    "fmt"
    "net"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // 1. 监听本地 9999 端口
    addr, _ := net.ResolveUDPAddr("udp", ":9999")
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Println("监听失败:", err)
        return
    }
    defer conn.Close()
    fmt.Println("UDP服务端启动，等待数据...")

    // 2. 对应 Java 的 byte[] b = new byte; DatagramPacket dp = new DatagramPacket...
    buf := make([]byte, 1024)

    for {
        // 阻塞接收数据包
        n, clientAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("接收数据错误:", err)
            continue
        }

        // 3. JSON 反序列化 (替代 Java 的 ObjectInputStream.readObject())
        var p Person
        err = json.Unmarshal(buf[:n], &p)
        if err != nil {
            fmt.Println("反序列化失败:", err)
            continue
        }

        fmt.Printf("收到来自 %v 的对象数据: Name=%s, Age=%d\n", clientAddr, p.Name, p.Age)
    }
}
```