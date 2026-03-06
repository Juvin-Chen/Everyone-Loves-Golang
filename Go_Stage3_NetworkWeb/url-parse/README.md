在 Java 中，网络编程依赖 `java.net` 包下的 `InetAddress`、`URL`、`Socket` 等类。**在 Go 语言中，我们主要使用 `net` 和 `net/url` 标准库。**

> **关于 Socket 的深层理解：** Socket 是应用层和传输层之间的桥梁。通信必须有两端：`Socket(IP, Port, 协议)` 组成的三元组代表一个端点。 在服务端，`ServerSocket` 就像公司的**总机接线员**（只负责监听端口），当有连接进来时（`accept()`），分配一个新的 `Socket`（相当于**员工分机**）去和客户端专门通信。