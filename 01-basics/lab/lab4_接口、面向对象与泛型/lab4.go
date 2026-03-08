package main

import "fmt"

/*
建立一个Go与其他语言有很大不同的认知：
在同一个包（package）内，顺序完全不重要。你就像在一个大房间里，所有东西（函数、结构体、变量）都堆在一起，你想拿哪个就拿哪个，不用管谁先谁后。
*/

type TaskFunc func()

func main() {
	// Go语言当中函数是一等公民，可以赋值给变量，可以作为参数传递，可以作为返回值
	// 创建任务列表
	tasks := []TaskFunc{task1, task2, task3, task4, task5, task6}
	for _, task := range tasks {
		task()
		fmt.Println()
	}
}
