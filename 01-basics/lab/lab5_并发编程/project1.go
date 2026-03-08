package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// 并发日志处理流水线
// Producer（生产者）- Consumer（消费者） 模型，Channel 串联

func project1() {
	fmt.Println("并发日志处理器：")
	chErr := make(chan string)
	chData := make(chan string)

	var wg sync.WaitGroup
	wg.Add(1) // 连锁反应，主协程只需要在最后一道关卡（消费者）设卡等待即可。

	// 生产者
	go func() {
		for i := 0; i < 3; i++ {
			str := "Log" + strconv.Itoa(i) + " get successfully!"
			chData <- str
		}
		for i := 0; i < 2; i++ {
			str := "Error" + strconv.Itoa(i)
			chData <- str
		}
		// 生产完毕，关闭 chData，否则处理者会一直死等
		close(chData)
	}()

	// 处理者
	go func() {
		// range可以自动持续读取直到通道被显式close()，所以通道必须要被手动关闭
		for msg := range chData {
			if strings.Contains(msg, "Err") {
				chErr <- msg
			} else {
				fmt.Println(msg + " Pass!")
			}
		}
		// 处理完毕，关闭 chErr，否则消费者会一直死等
		close(chErr)
	}()

	// 消费者
	go func() {
		defer wg.Done() // 通知主协程，文件写入工作全部完成

		file, err := os.OpenFile("./project1_ErrorFile.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("打开文件失败：%v\n", err)
			return
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		defer func() {
			if err := writer.Flush(); err != nil {
				fmt.Printf("刷新缓冲区失败：%v\n", err)
			}
		}()

		for errMsg := range chErr {
			if errMsg == "" {
				continue
			}

			_, writeErr := writer.WriteString(errMsg + "\n")
			if writeErr != nil {
				fmt.Printf("写入错误信息失败：%v\n", writeErr)
				continue
			}
		}
	}()

	wg.Wait()
	fmt.Println("并发日志处理器工作完成。")
}
