package main

import "fmt"

type TaskFunc func()

func main() {
	// tasks := []TaskFunc{task1, task2, task3, task4, task5}
	// for _, task := range tasks {
	// 	fmt.Println()
	// 	task()
	// }
	project1()
	fmt.Println()
	project2()
}
