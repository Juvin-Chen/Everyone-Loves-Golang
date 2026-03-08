package main

import (
	"fmt"
	"sync"
	"time"
)

// 限制并发的爬虫 (Worker Pool Pattern)
// 控制并发数量（Semaphore/Worker Pool），防止把服务器搞崩。
// Worker Pool 模型
// 这里用for range 就足够合适，select 的强项在于同时监听多个通道（比如既要接收任务，又要监听一个“取消”信号，或者设置超时时间）。但在标准的 Worker Pool 模型中，Worker 们的职责很单一：只需要死死盯住 jobs。

func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker [%d] 正在处理任务 [%s]\n", id, job)
		time.Sleep(1 * time.Second) // 模拟耗时

		// 假设有的任务可能会失败，从而不产生结果，也不影响整体运行
		// 可以说是边处理边打印，具有连续效果
		results <- fmt.Sprintf("【处理结果】: %s (由 Worker %d 完成)", job, id)
	}
}

func project2() {
	const numJobs = 20
	const numWorkers = 3

	jobs := make(chan string, numJobs)
	results := make(chan string, numJobs)

	// 创建 WaitGroup 监工
	var wg sync.WaitGroup

	fmt.Println("--- 启动 Worker Pool ---")
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1) // 每雇佣一个 Worker，监工的计数器加 1
		go worker(w, jobs, results, &wg)
	}

	// 核心重构：启动一个专属的“关灯协程”，因为需要等待所有工人把任务完成之后results才能关闭
	go func() {
		wg.Wait()
		close(results)
	}()

	// 把任务塞进通道
	for j := 1; j <= numJobs; j++ {
		jobs <- fmt.Sprintf("URL_%d", j)
	}
	// 任务发放完毕，关闭 jobs 通道，这会让 Worker 们的 range 循环最终能够结束
	close(jobs)

	fmt.Println("--- 开始收集结果 ---")
	// 未知数量收集法：一直 range，直到 results 通道被“关灯协程”关闭
	for res := range results {
		fmt.Println("收集到:", res)
	}

	fmt.Println("所有结果收集完毕，主程序完美退出！")
}
