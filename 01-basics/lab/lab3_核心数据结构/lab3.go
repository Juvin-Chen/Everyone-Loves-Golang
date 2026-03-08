package main

import (
	"fmt"
)

type TaskFunc func()

func main() {
	tasks := []TaskFunc{task1, task2, task3}

	for _, task := range tasks {
		task()
		fmt.Println()
	}

	fmt.Println("lab3_Project:")
	project()
}
