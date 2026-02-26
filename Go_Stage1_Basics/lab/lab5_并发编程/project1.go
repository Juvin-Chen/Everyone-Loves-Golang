package main

import (
	"fmt"
)

// 并发日志处理流水线
// Producer（生产者）- Consumer（消费者） 模型，Channel 串联

func project1() {
	fmt.Println("并发日志处理器：")
	ch := make(chan string)
	go func() {
		for i :=0 ;i<5;i++{
			str := "Log"+
			ch<-str
		}
	}()
}
