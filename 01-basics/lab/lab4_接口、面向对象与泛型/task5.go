package main

import "fmt"

// 接口的“非空”幻觉,制造一个“不等于 nil 的 nil”
/*
在 Go 中，接口不是简单的指针，而是包含两个部分的结构体：
1.类型（Type）：表示接口持有的具体类型（例如 *Cat）
2.值（Value）：指向实际数据的指针（例如 *Cat 指针）
*/

type Animal interface {
}

type Cat struct {
}

func task5() {
	fmt.Println("接口的“非空”幻觉:")

	// 情况 1: 指针是 nil
	var c *Cat = nil
	fmt.Println("c == nil?", c == nil) // true

	// 情况 2: 接口持有 nil 指针（陷阱！）
	var a Animal = c
	fmt.Println("a == nil?", a == nil) // false (但内部值是 nil)

	// 修正：让接口真正等于 nil
	a = nil
	fmt.Println("a == nil after fix?", a == nil) // true
}
