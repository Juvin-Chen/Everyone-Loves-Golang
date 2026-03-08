package main

import "fmt"

// 解包盲盒
// `interface{}` 是万能容器，学会用 `switch type` 优雅地处理未知数据。这在解析 JSON 时非常常用。
func task4() {
	fmt.Println("解包盲盒，测试开始：")
	Uppack(10)
	Uppack("Hello")
	Uppack([]int{1, 2, 3})
	Uppack(true)
}

func Uppack(box interface{}) {
	fmt.Println("解包盲盒:")

	// 提取空接口的实际类型（仅 type switch 可用），这里是一个独特用法.(type)
	switch v := box.(type) {
	case int:
		fmt.Println("是数字，它的平方是:", v*v)
	case string:
		fmt.Println("是文本，内容是:", v)
	case []int:
		fmt.Println("是整数切片，长度是:", len(v))
	default:
		fmt.Println("未知类型")
	}
}
