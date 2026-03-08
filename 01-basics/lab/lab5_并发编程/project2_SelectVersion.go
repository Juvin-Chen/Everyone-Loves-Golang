package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 该版本添加了紧急停止的功能，并且加了些emoji~
// worker 函数现在多接了两个参数：
// ctx：用来听广播（听有没有人按紧急按钮）
// cancel：给工人配发的“红色紧急按钮”，谁遇到致命错误谁就按

func worker_2(ctx context.Context, id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup, cancel context.CancelFunc) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done(): // 只要有人调用了 cancel()，这里就会收到信号
			fmt.Printf("🔴 Worker [%d] 听到全厂警报，立刻丢下手中的活，紧急撤离！\n", id)
			return // 直接 return，结束协程

		case job, ok := <-jobs:
			if !ok {
				// jobs 通道被老板正常关闭了，且没有剩余任务，正常下班
				return
			}

			fmt.Printf("Worker [%d] 正在处理任务 [%s]\n", id, job)

			// 🚨 模拟遇到致命错误（遇到 URL_7 就爆炸）
			if job == "URL_7" {
				fmt.Printf("💣 Worker [%d] 发现服务器异常 (URL_7)！立马按下紧急停止按钮！\n", id)
				cancel() // 按下按钮，触发全厂广播！
				continue // 放弃当前任务
			}

			// 模拟正常的下载耗时
			time.Sleep(500 * time.Millisecond)

			// 干完活，交差
			results <- fmt.Sprintf("【处理结果】: %s (由 Worker %d 完成)", job, id)
		}
	}
}

func project_2() {
	const numJobs = 20
	const numWorkers = 3

	jobs := make(chan string, numJobs)
	results := make(chan string, numJobs)
	var wg sync.WaitGroup

	// 🚨 1. 核心组装：创建一个带“取消功能”的 context (背景上下文)
	// ctx 是广播喇叭，cancel 是紧急停止按钮
	ctx, cancel := context.WithCancel(context.Background())

	// 好习惯：无论如何，main 函数退出前确保调用一下 cancel，释放底层资源
	defer cancel()

	fmt.Println("--- 启动带紧急制动的 Worker Pool ---")
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker_2(ctx, w, jobs, results, &wg, cancel)
	}

	// 老规矩：关灯大爷
	go func() {
		wg.Wait()
		close(results)
	}()

	// 把任务塞进通道
	for j := 1; j <= numJobs; j++ {
		jobs <- fmt.Sprintf("URL_%d", j)
	}
	close(jobs) // 任务发完了，关闭传送带

	fmt.Println("--- 开始收集结果 ---")
	for res := range results {
		fmt.Println("✅ 收集到:", res)
	}

	fmt.Println("主程序结束运行！")
}
