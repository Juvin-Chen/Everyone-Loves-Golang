package main

import "fmt"

// 方法集的严格约束，编译器报错挑战

type Mover interface {
	Move()
}

type Dog struct {
	Name string
}

func (d *Dog) Move() {
	fmt.Printf("Dog[%s] is moving\n", d.Name)
}

func task6() {
	fmt.Println("编译器报错挑战:")
	var m Mover = &Dog{"poppy"}
	m.Move()
	/*
		方法调用 vs 接口赋值（关键区别） /  对比task3
		方法调用
		(值类型调用指针接收者方法)	lamp.TurnOn()	 自动转换为 (&lamp).TurnOn()	运行时Go 会自动取地址
		接口赋值
		(值类型赋值给接口)	var m Mover = Dog{} / var usb []USB = []USB{Phone{}} （错误！）	 不能自动转换	编译时Go 需要确定类型是否实现接口
	*/
}

/*
1.思考：为什么指针接收者的方法，不能被值类型实现？
提示：接口里的值如果是不可寻址的，怎么调用修改状态的指针方法？

答：指针接收者的方法不能被值类型实现，因为值类型（非指针）在编译时无法被识别为该方法的接收者类型，而接口需要在编译时确定类型是否实现接口。

Go 语言的设计哲学：在编译时确保安全，在运行时提供便利。
*/
