package main

import (
	"fmt"
)

// task1:万能电子设备充电站

type USB interface {
	Start()
	Stop()
}

type Phone struct {
	Name string
}

type Camera struct {
	Name string
}

func (p *Phone) Start() { fmt.Println(p.Name, "开始充电") }

func (p *Phone) Stop() { fmt.Println(p.Name, "结束充电") }

func (c *Camera) Start() { fmt.Println(c.Name, "开始充电") }

func (c *Camera) Stop() { fmt.Println(c.Name, "结束充电") }

// 核心规则：
// 1. 当方法用「值接收者」实现接口时：值类型（Phone{}）和指针类型（&Phone{}）都能赋值给接口（因为指针可自动解引用调用值接收者方法）；
// 2. 当方法用「指针接收者」实现接口时：只有指针类型（&Phone{}）能赋值给接口（值类型无法调用指针接收者方法，无法传递地址）；
// 3. 如果选择值接收者实现接口，既兼容值传递（task1_copy），也兼容指针传递（task1_ptr），更易理解。

func task1() {
	fmt.Println("万能电子设备充电站：")
	var usb []USB = []USB{&Phone{"iPhone 17 pro"}, &Camera{"dj pocket3"}} // 上面的接口是用指针类型实现的
	Computer(usb)
}

func Computer(devices []USB) {
	for _, device := range devices {
		device.Start()
		if ptr, ok := device.(*Phone); ok && ptr != nil {
			fmt.Println("检测到手机，开启快速充电模式")
		}
		device.Stop()
	}
}
