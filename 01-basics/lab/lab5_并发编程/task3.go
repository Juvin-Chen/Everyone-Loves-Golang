package main

import (
	"fmt"
	"time"
)

// 不可靠的服务器，Select 的“超时控制”

func rpcCall(c chan string) {
	// 单位可选：time.Millisecond（毫秒）、time.Second（秒）、time.Minute（分钟）
	time.Sleep(2 * time.Second)
	c <- "Response Data"
}

func task3() {
	fmt.Println("不可靠的服务器！")
	fmt.Println("网络请求如果卡死了怎么办？Select 是最好的解药。")
	fmt.Println("对服务器分别进行两种情况的测试：\n")

	var chA chan string
	chA = make(chan string)
	go rpcCall(chA)

	fmt.Println("测试情况 A：超时设为 1 秒（小于 RPC 耗时）")
	select {
	case res := <-chA:
		fmt.Printf("Received: %s\n", res)
	case <-time.After(1 * time.Second): // 监听到超时（1秒到了）
		fmt.Println("Timeout! Request canceled")
	}

	fmt.Println()

	fmt.Println("测试情况 B：超时设为 3 秒（大于 RPC 耗时）")
	chB := make(chan string)
	go rpcCall(chB)
	select {
	case res := <-chB:
		fmt.Printf("Received: %s\n", res)
	case <-time.After(3 * time.Second):
		fmt.Println("Timeout! Request canceled")
	}

	// 这个场景下为什么不需要额外加 WaitGroup 等待？
	// select 本身就是 “阻塞等待”，主协程执行到 select 时，会停在这一行不动，直到case中的情况发生
}
