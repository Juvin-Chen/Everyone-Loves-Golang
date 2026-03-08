package main

import (
	"context"
	"fmt"
	"time"
)

// 搜索雷达，Context 的“一键止停”
// Context 就是 “控制协程生命周期” 的工具

func search(ctx context.Context, name string) {
	for {
		select {
		// case1：收到停止指令（ctx.Done()通道关闭）
		case <-ctx.Done():
			fmt.Printf("[%s] 收到停止指令，正在关闭...\n", name)
			return
		default:
			fmt.Printf("[%s] 正在搜索...\n", name)
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func task5() {
	fmt.Println("搜索雷达，Context 的“一键止停”")
	// context.WithCancel 创建上下文：做一个 “总开关”（ctx 是信号通道，cancel 是开关按钮）
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动 3 个雷达协程
	go search(ctx, "Radar-A")
	go search(ctx, "Radar-B")
	go search(ctx, "Radar-C")

	fmt.Println("主控室：所有雷达启动，开始搜索...")
	time.Sleep(2 * time.Second)

	fmt.Println("\n主控室：2秒已到，发送停止指令！")
	cancel() // 调用 cancel()：按下总开关，给所有雷达发 “停止指令”

	time.Sleep(1 * time.Second)
	fmt.Println("\n主控室：所有雷达已关闭，任务结束！")

	// 这个程序的功能简单讲就是：让多个一直干活的协程，在你指定的时间点 “一键全部停止”，而不是让它们无限循环下去
}
