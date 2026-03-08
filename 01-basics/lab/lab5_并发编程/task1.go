package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 并发下载器模拟，WaitGroup 的“点名”艺术

func download(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Start downloading [%s]...\n", url)

	randomMs := 1000 + rand.Intn(1000)
	sleepDuration := time.Duration(randomMs) * time.Millisecond // 单位转换
	time.Sleep(sleepDuration)

	fmt.Printf("Finish [%s]\n", url)
}

func task1() {
	// 初始化随机数种子（Go 1.20 及以上版本可省略这行）
	// rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup

	fmt.Println("并发下载器模拟:")
	website1 := "https://IDEA.com"
	website2 := "https://Cursor.com"
	website3 := "https://Trae.com"
	website4 := "https://VSCode.com"

	urls := []string{website1, website2, website3, website4}
	for _, url := range urls {
		wg.Add(1)
		go download(url, &wg)
	}

	wg.Wait()
	fmt.Println("All downloads completed")
}
