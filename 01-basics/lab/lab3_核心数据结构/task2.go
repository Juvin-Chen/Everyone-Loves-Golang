package main

import "fmt"

// 任务2：动态数据流处理
func task2() {
	// 1.切片删除 (Delete)：
	fmt.Println("==切片删除==")
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} //没有大小，即为切片，切片是引用类型传递
	data = append(data[:4], 6, 7, 8, 9)
	fmt.Println("删除指定数字5后的切片：")
	for i := range data {
		fmt.Printf("%d ", data[i])
	}

	// 2.切片备份
	fmt.Println("\n==切片备份==")
	backup := make([]int, len(data), cap(data))
	copy(backup, data)
	for i := range data {
		fmt.Printf("%d ", data[i])
	}
	backup[0] = 999
	fmt.Println("\n判断数据是否备份正确:")
	for i := range backup {
		fmt.Printf("%d ", backup[i])
	}

	// 3.扩容观察
	fmt.Println("\n==扩容观察==")
	for i := 0; i < 5; i++ {
		data = append(data, 100)
		fmt.Printf("此时的len:%d,cap:%d\n", len(data), cap(data))
	}
}
