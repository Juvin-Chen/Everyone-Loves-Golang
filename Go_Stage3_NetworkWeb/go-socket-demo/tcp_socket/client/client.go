// 客户端 Client
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// 客户端发起 TCP 连接请求
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接服务器失败：", err)
		return
	}
	defer conn.Close()

	// 发送消息给服务端
	conn.Write([]byte("你好，服务端！\n"))

	// 接收服务器的回复
	reply, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("服务器回复", reply)
}
