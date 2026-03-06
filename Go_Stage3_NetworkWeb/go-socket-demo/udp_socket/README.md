### UDP 编程：基本数据类型与对象的传输 

UDP 不需要提前建连，在 Go 中使用 `net.ListenUDP` 和 `net.DialUDP`。

**传输自定义对象 (使用 JSON 序列化代替 Java 的 Serializable)** 在 Java 中，传对象必须要实现 `Serializable` 接口并使用 `ObjectOutputStream`。 在现代 Web 尤其是 Go 语言开发中，**传递结构体（对象）最标准、跨语言的做法是将其序列化为 JSON 字节数组。**