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
