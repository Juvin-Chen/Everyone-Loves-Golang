# NoteBook02 网络编程

承接NoteBook备份01

## 第五章：Go 语言网络编程实战 (对标 Java)

在 Java 中，网络编程依赖 `java.net` 包下的 `InetAddress`、`URL`、`Socket` 等类。**在 Go 语言中，我们主要使用 `net` 和 `net/url` 标准库。**

> **关于 Socket 的深层理解：** Socket 是应用层和传输层之间的桥梁。通信必须有两端：`Socket(IP, Port, 协议)` 组成的三元组代表一个端点。 在服务端，`ServerSocket` 就像公司的**总机接线员**（只负责监听端口），当有连接进来时（`accept()`），分配一个新的 `Socket`（相当于**员工分机**）去和客户端专门通信。

### 5.1 解析 URL (对应 Java 的 `URL` 类)

```go
// 解析url
// 在 Go 语言中，我们主要使用 net 和 net/url 标准库
// url 格式：协议://（主机IP）服务器域名:端口号/路径?参数名=参数值。
package main

import (
	"fmt"
	"net/url"
)

// 把一个完整的 URL 字符串（比如百度搜索的 URL）解析成 Go 语言能识别的结构化对象，然后提取出「协议、主机名、路径、查询参数」这些关键信息。
func main() {
	// 这里用的就是教科书版的完整的 url 用于解析，实际上大部分是省略了路径和参数的极简版：https://www.baidu.com
	rawUrl := "https://www.baidu.com/s?wd=go%E8%AF%AD%E8%A8%80&rn=10&tn=baidu"
	// 解析原始 URL 字符串，返回*url.URL对象
	parseUrl, err := url.Parse(rawUrl)
	if err != nil {
		// Go的内置函数，调用后立即停止当前程序的执行，打印错误信息 + 调用栈（方便定位哪里错了）；
		panic(err)
	}

	fmt.Println("协议（Protocol）:", parseUrl.Scheme)
	fmt.Println("主机名（Host）:", parseUrl.Host)
	fmt.Println("路径（Path）:", parseUrl.Path)

	// 获取具体的参数值，把 URL 中 ? 后面的查询参数（wd=go%E8%AF%AD%E8%A8%80&rn=10&tn=baidu）解析成 url.Values 类型（本质是键值对的 map）
	queryParams := parseUrl.Query()
	fmt.Println("wd的值：", queryParams.Get("wd")) // 上面的 url ?后面是 wd
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

type Person struct {
	Name string `json:"Name"`
	Age  int    `json:Age`
}

func main() {
	// net 包提供的专门用于解析 UDP 地址的函数，得到 *net.UDPAddr 类型对象，主要目的就是把它转换成一个 Go 语言能看懂的专用地址对象 serverAddr
	// 这行代码的意义是找服务端
	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9999")

	// 本地客户端分配随机 UDP 端口
	// func DialUDP(network string, laddr, raddr *UDPAddr) 参数1 网络协议，参数2 本地客户端的地址，随机分配传参nil就可以，参数3 服务端
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("建立UDP连接失败", err)
		return
	}
	defer conn.Close()

	p := Person{Name: "Angelababy", Age: 37}
	// json.Marshal(p)：把 Person 结构体（Go 语言内部格式）转换成 JSON 字符串（字节切片 []byte）
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println("序列化失败", err)
		return
	}

	// 发送UDP数据报文
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("发送失败")
	} else {
		fmt.Println("对象发送成功")
	}
}
```

**UDP 服务端 (接收方)**

```go
// UDP 服务端程序，核心功能是监听 9999 端口，接收客户端发来的 JSON 格式数据，反序列化成 Person 结构体并打印。

package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Person struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
}

func main() {
	// 准备监听地址，监听本地 9999 端口，前面没有 IP，意思是“监听本机所有网卡上的 9999 端口”
	// 在计算机网络中，IP 为空或 0.0.0.0 代表“本机所有可用的网络接口”。
	// 这意味着：不管别人是从 127.0.0.1 访问你，还是从局域网 IP 192.168.1.5 访问你，只要端口是 9999，你都能收到。
	addr, _ := net.ResolveUDPAddr("udp", ":9999")

	// conn：*net.UDPConn 类型的 UDP 连接对象，后续收发数据都靠它
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("监听失败", err)
		return
	}

	defer conn.Close()
	fmt.Println("UDP服务端启动，等待数据...")

	// 创建接收数据的缓冲区，创建一个长度为 1024 字节的字节切片
	buf := make([]byte, 1024)

	for {
		// conn.ReadFromUDP(buf)：阻塞等待客户端发送 UDP 数据包
		// n：实际接收到的字节数（比如客户端发了 20 字节，n=20）；
		// clientAddr：*net.UDPAddr 类型，包含客户端的 IP 和端口（比如 127.0.0.1:54321）；
		// err：接收数据时的错误（比如连接中断）。
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("接收数据错误", err)
			continue
		}

		// json反序列化
		var p Person
		err = json.Unmarshal(buf[:n], &p)
		if err != nil {
			fmt.Println("反序列化失败：", err)
			continue
		}
		fmt.Printf("收到来自 %v 的对象数据：Name = %s, Age = %d\n", clientAddr, p.Name, p.Age)

	}
}
```



