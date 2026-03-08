package main

import "fmt"

// 任务1：数组倍增器
func task1() {
	arr := [5]int{10, 20, 30, 40, 50} // 只要带容量大小的都是数组，因为数组的大小算是它的一个特性
	doubleArrayByVal(arr)
	fmt.Println("1.错误演示，按值传递的方式：")
	for _, value := range arr {
		fmt.Printf("%d ", value)
	}
	fmt.Println("\n2.正确演示，按指针传递的方式：")
	doubleArrayByPtr(&arr)
	for _, value := range arr {
		fmt.Printf("%d ", value)
	}
}

// 错误示范：值传递
func doubleArrayByVal(arr [5]int) {
	for _, v := range arr {
		v *= 2
	}
}

// 正确示范：指针传递
func doubleArrayByPtr(arr *[5]int) {
	// 要求：在函数内部，使用指针解引用或直接索引的方式，将原数组每个元素乘以 2。
	for i := range arr {
		arr[i] *= 2
	}
}

/*
思考题 (写在注释里)：
- 在 `doubleArrayByPtr` 中，访问元素时写 `(*arr)[i]` 和 `arr[i]` 有区别吗？Go 编译器在这里做了什么？
- 答：解引用遍历的本质是 (*arr)[i]，但 Go 允许简化成 arr[i]，Go语言会自动解引用。
*/
