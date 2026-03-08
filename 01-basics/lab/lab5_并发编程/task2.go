package main

import (
	"fmt"
	"sync"
)

// 乒乓球比赛 (Ping-Pong)，Channel 的“接力赛”，实现双方选手的对打
// 理解无缓冲通道的“阻塞”特性——不见不散

// 分别写了用 1 / 2 个通道实现的方式，随便写写的，顺便加点emoji :)

func PlayerA(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done() // 协程结束标记完成
	for i := 0; i < 5; i++ {
		if i != 0 {
			data := <-ch
			fmt.Printf("A received: [%s]\n", data)
		}
		ch <- "Ping"
		fmt.Println("A sent Ping🏓") // 发送后打印
	}
	data := <-ch
	fmt.Printf("A received: [%s] (Final)\n", data)

	close(ch) // 接完最后一次再关闭通道
}

func PlayerB(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		data := <-ch
		fmt.Printf("B received: [%s]\n", data)
		ch <- "Pong"
		fmt.Println("B sent Pong🏓")
	}
}

func task2() {
	var wg1 sync.WaitGroup
	wg1.Add(2)
	fmt.Println("乒乓球比赛🏆 :")
	fmt.Println("选手A和选手B正在场上进行乒乓球🏓比赛：")

	fmt.Println("\n实现方式1：调用Player函数，只用一个channel")
	fmt.Println("\n第一轮比赛正式开始🎉")
	ch1 := make(chan string)

	go PlayerA(ch1, &wg1)
	go PlayerB(ch1, &wg1)

	wg1.Wait()

	fmt.Println("\n\n实现方式2：匿名函数实现，使用两个channel")
	fmt.Println("\n第二轮比赛正式开始🎉")
	var wg2 sync.WaitGroup
	wg2.Add(2)
	pingch := make(chan string)
	pongch := make(chan string)

	go func() {
		defer wg2.Done()
		for i := 0; i < 5; i++ {
			pingch <- "Ping" // 1. A 先发球
			fmt.Println("A sent Ping🏓")
			data := <-pongch // 2. A 接 B 打回来的球
			fmt.Printf("A received: [%s]\n", data)
		}
		close(pingch) // close()在Go语言里面也不是一定要写，Go 的垃圾回收器（GC）会自动把它们清理掉，不会造成内存泄漏。
	}()
	go func() {
		defer wg2.Done()
		for i := 0; i < 5; i++ {
			data := <-pingch
			fmt.Printf("B received: [%s]\n", data)
			pongch <- "Pong"
			fmt.Println("B sent Pong🏓")
		}
		close(pongch)
	}()

	wg2.Wait()
}

/*
运行程序，观察打印日志。是 A 全部发完 B 再收？还是 A 发一个，B 收一个？
答：是 A 发一个，B 收一个。

思考：为什么无缓冲通道能起到“同步”的作用？
答：因为无缓冲通道的容量为 0，它不存储任何数据。发送方和接收方必须在通道的两端同时准备好，就像接力赛交接棒一样，必须两个人的手同时握住接力棒，交接才能完成。
少任何一方，先到的一方就只能原地干等（阻塞）。
*/
