package main

import (
	"fmt"
)

// 智能家居系统

// 基类组件
type Component struct {
	ID     int
	Status string
}

func (c *Component) TurnOn() {
	fmt.Printf("Component [%d] is On\n", c.ID)
}

/*接收者类型一致 → 外部方法必遮蔽；接收者类型不一致 → 看调用实例是值 / 指针，匹配对应接收者的方法。*/

// 普通字段：Component 只是一个普通成员，必须通过 lamp.component 访问其字段 / 方法；
// 匿名嵌入：Component 是嵌入结构体，其字段 / 方法会被 “提升” 到外部结构体，可直接通过 lamp 访问。

// 具体设备，注意只写一个Component才是匿名嵌入
// 智能灯
type SmartLamp struct {
	Component
	Brightness int
}

// 空调
type AirConditioner struct {
	Component
	Temperature int
}

// 方法遮蔽 (Shadowing)：给 `AirConditioner` 重写 `TurnOn()` 方法
func (a *AirConditioner) TurnOn() {
	fmt.Println("AirConditioner is cooling down...")
}

func task3() {
	fmt.Println("智能家居系统:")

	fmt.Println("正在创建一个灯和一个空调...")
	lamp := SmartLamp{Component{1001, "normal"}, 40}
	airc := AirConditioner{Component{1002, "normal"}, 26}
	fmt.Printf("灯的ID为%d\n", lamp.ID)

	fmt.Println("正在执行两个家具的turn on 操作")
	lamp.TurnOn() // 指针接收者方法可以被值类型调用，Go 会自动处理。
	airc.TurnOn()
}
