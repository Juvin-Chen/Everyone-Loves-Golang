package main

import (
	"fmt"
	"sync"
)

// 银行账户抢红包，数据竞争与 Mutex
// 多个协程改同一个变量，结果不对channel，Java风格是加锁🔒
// Go语言优先选择

var balance int

func task4() {
	fmt.Println("银行账户抢红包")
	fmt.Println("初始balance为0")
	fmt.Println("\n不安全版：")
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			balance += 1
		}()
	}
	wg.Wait()
	fmt.Println("🎉 （不安全版）所有协程执行完毕！balance金额为", balance)

	fmt.Println("\n安全版（加锁🔒）：")
	balance = 0
	var sm sync.Mutex
	wg.Add(1000) // 直接复用上面的wg，之前的任务中还多创建了但其实没必要
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			// 加锁：同一时间只有一个协程能改balance
			sm.Lock()
			defer sm.Unlock() // 解锁：其他协程可以继续改
			balance += 1
		}()
	}
	wg.Wait()
	fmt.Println("🎉 （安全版）所有协程执行完毕！balance金额为", balance)
}
