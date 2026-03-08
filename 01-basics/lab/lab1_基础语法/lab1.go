package main

import "fmt"

var tasks []func()

func main() {
	// 设置一个任务调度器~
	// 把所有任务函数添加到切片中
	tasks = append(tasks, task1)
	tasks = append(tasks, task2)
	tasks = append(tasks, task3)
	tasks = append(tasks, task4)
	tasks = append(tasks, task5)
	tasks = append(tasks, task6)

	for i, task := range tasks {
		fmt.Printf("\n--- 执行第 %d 个任务 ---\n", i+1)
		task()
	}

	fmt.Println("\n 所有任务已全部完成！")
}

// task1:个人名片生成器
func task1() {
	var name string = "Juvin"
	var age int = 20
	is_active := true
	var salary float64 = 20000.00

	// "\t"制表符打印，\t会自动跳到(0,8,16...)的位置，原sprintf里面的空格 like" | "，会原封不动照搬
	fmt.Println("方式1-使用制表符打印：")
	fmt.Printf(" Name\t| Age\t| IsActive\t| Salary\n")
	infoStr := fmt.Sprintf("%q\t| %d\t| %t   \t| %.2f", name, age, is_active, salary)
	fmt.Println(infoStr)

	fmt.Print("\n------------------------------------------\n\n")

	// 固定格数
	fmt.Println("方式2-使用固定宽度打印：")
	fmt.Printf("%-8s | %-4s | %-10s | %-10s\n", " Name", "Age", "IsActive", "Salary")
	// 内容：用 %-Ns 左对齐，控制每列宽度
	infoStr = fmt.Sprintf("%-8q | %-4d | %-10t | %-10.2f", name, age, is_active, salary)
	fmt.Println(infoStr)
}

// task2:超市收银系统
func task2() {
	// go禁止不同变量的隐式计算如int*float
	var price float64 = 19.9
	count := 5
	var total float64
	total = price * float64(count)
	fmt.Printf("购买了%d件商品,单价%.1f元，总共花费%.2f元\n", count, price, total)
	//额外作业：声明一个string变量不赋值，直接打印
	var str string
	fmt.Println("var nil string: ", str)
}

// task3:权限管理模拟
func task3() {
	//在系统开发中常用二进制位表示权限
	const (
		Readable = 1 << iota
		Writable
		Executable
	)
	var userPerm int = Readable | Writable
	if userPerm&Executable == Executable { // & 仅是提取执行这个权限
		fmt.Println("Permission Allowed")
	} else {
		fmt.Println("Permission Denied")
		userPerm |= Executable
	}
	fmt.Println("Final permission value:", userPerm) //应该是7
}

// task4:硬盘单位转换器
func task4() {
	// 编写一个能够打印存储单位字节数的程序,int64是64位有符号整数，最大值是 2^63 - 1
	// const本身是无类型常量，等到被使用的时候才会确定类型
	const (
		_ = 1 << (iota * 10)
		KB
		MB
		GB
		TB
	)
	println("1GB对应的字节数: ", GB) //自动推导为int
}

// task5 : version2 , ! 字节转存储单位
type ByteSize int

// 为 ByteSize 这个类型的 Println 定义的函数，打印的时候会自动格式化
func (b ByteSize) String() string {
	div := 1
	var suffix string // 注意：Go 语言里的字符串，永远不可能是 nil，它是值类型，用""表示空字符串
	switch {
	case b < (1 << 10):
		return fmt.Sprintf("%d B", b)
	case b < (1 << 20):
		div = (1 << 10)
		suffix = "KB"
	case b < (1 << 30):
		div = (1 << 20)
		suffix = "MB"
	case b < (1 << 40):
		div = (1 << 30)
		suffix = "GB"
	default:
		div = (1 << 40)
		suffix = "TB"
	}

	value := float64(b) / float64(div)
	return fmt.Sprintf("%.2f %s", value, suffix)
}

func task5() {
	// use ByteSize 类型自动格式化成 “1.00 GB”,提前用了点接口））
	const (
		_ ByteSize = 1 << (iota * 10)
		KB
		MB
		GB
		TB
	)

	var size ByteSize = 1024 * 1024 // 1MB
	fmt.Println("Size : ", size)

	var bigSize ByteSize = 1024 * 1024 * 1024 * 1024 // 1TB
	fmt.Println("Big Size:", bigSize)

	var tinySize ByteSize = 1500 // 约 1.46 KB
	fmt.Println("Tiny Size:", tinySize)

}

// task6:变量交换与匿名变量

func fakeFetch() (int, string) {
	return 200, "Success"
}

func task6() {

	//这里不能用var，_不能用于作为定义的变量名称
	num, _ := fakeFetch()
	println("忽略提示信息，仅打印状态码：", num)

	x := 10
	y := 20
	println("before swap:", x, y)
	x, y = y, x
	println("after swap:", x, y)
}
