// UDP 客户端

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
