// 服务端 Server
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// 作用：创建 TCP 监听器，占用本机 8888 端口，等待客户端连接；
	// 参数："tcp"= 用 TCP 协议，":8888"= 监听本机 8888 端口；
	// 返回：listener= 监听器（负责等连接），err= 启动失败的错误（比如端口被占）。
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("启动监听失败:", err)
		return
	}

	// 设置“总机关机时，自动拔掉接线器”（释放端口资源）
	defer listener.Close()
	fmt.Println("服务端已启动，正在监听 8888 端口...")

	// 该版本支持多个客户端同时连接
	for {
		// 子步骤1：等待客户打电话（阻塞式）
		// listener.Accept() → 调用监听器的Accept方法，相当于“接线员等电话”
		// 特性：阻塞 → 如果没有客户端连接，程序会停在这里，直到有客户端打进来
		// 返回2个值：conn（单个客户的通话线路，socket连接）、err（连接失败的错误）
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// 子步骤3：派一个“专属接线员”处理这个客户（并发），支持多客户端同时连接
		go handleClient(conn)
	}
}

// 专门处理单个客户端，参数是conn（客户端的通话线路，net.Conn类型）
func handleClient(conn net.Conn) {
	defer conn.Close()

	// 1.给这个通话线路配一个“听筒”（带缓冲的读取器）
	// bufio.NewReader(conn) → 基于conn创建读取器，方便“按行读取”客户端发的消息
	// 因为客户端发的消息是按行结尾的（\n），用这个工具能精准读到一行消息
	reader := bufio.NewReader(conn)

	// 2.听客户说话（读取客户端发的消息）
	// reader.ReadString('\n') → 从读取器里读数据，直到遇到换行符\n为止
	// 返回2个值：msg（读到的消息内容）、_（错误信息，这里用_忽略，新手先不处理）
	msg, _ := reader.ReadString('\n')

	// 3.把客户说的话记下来（打印到控制台）
	fmt.Print("收到客户端消息：", msg)

	// 4.回复客户端，conn.Write() 只能接收 byte 切片（网络传输只认二进制字节）
	conn.Write([]byte("你好，客户端，已收到你的消息！\n"))
}
