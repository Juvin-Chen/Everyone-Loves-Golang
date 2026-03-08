package main

import "fmt"

/* 综合实战项目：简易库存管理系统,这是一个模拟电商后台的小程序。你需要结合结构体（定义商品）、Map（存储数据）、指针（修改状态）和 函数。 */
type Product struct {
	id    int
	name  string
	price float64
	stock int // 库存数量
}

func project() {
	/*
		note:
		只是声明了一个 nil map，直接赋值会panic，必须用 m := make(...)或者字面量初始化。
		var inventory map[int]*Product
		inventory = make(map[int]*Product)
		字面量初始化（直接写值）
		m1 := map[string]int{"Alice": 1, "Bob": 2}
		非字面量初始化（用make函数）
		m2 := make(map[string]int)
	*/

	iPhone := &Product{101, "iPhone", 5000.00, 20}
	Laptop := &Product{102, "Laptop", 20.00, 5}
	Bag := &Product{103, "Bag", 65.00, 2}
	Basketball := &Product{104, "Basketball", 90.00, 1}
	inventory := map[int]*Product{
		101: iPhone,
		102: Laptop,
		103: Bag,
		104: Basketball,
	}
	showInventory(inventory)
	addProduct(inventory, 105, "Cap", 18.00, 2)
	buyProduct(inventory, 105, 2)
	buyProduct(inventory, 105, 1)
	addProduct(inventory, 105, "Cap", 18.00, 2) // 再添加两个库存
	showInventory(inventory)
}

// 模块功能实现
func showInventory(m map[int]*Product) {
	fmt.Println("=== 商品库存清单 ===")
	for _, p := range m {
		fmt.Printf("商品ID: %d\n", p.id)
		fmt.Printf("名称:   %s\n", p.name)
		fmt.Printf("价格:   %.2f\n", p.price)
		fmt.Printf("库存:   %d\n", p.stock)
		fmt.Println("-------------------")
	}
}

func addProduct(m map[int]*Product, id int, name string, price float64, stock int) {
	product, ok := m[id]
	if ok {
		fmt.Printf("商品%s已存在，正在为该商品添加库存\n", name)
		product.stock += stock
		fmt.Printf("已添加完成，该商品现在库存数量为%d\n", product.stock)
	} else {
		fmt.Printf("商品%s不存在，正在创建该新商品\n", name)
		newproduct := &Product{id, name, price, stock}
		m[id] = newproduct
	}
}

func buyProduct(m map[int]*Product, id int, count int) {
	product, ok := m[id]
	if ok {
		if product.stock < count {
			fmt.Printf("%s商品库存不足\n", product.name)
		} else {
			product.stock -= count
			fmt.Printf("购买%s商品成功，共购买%d件商品，花费 %.2f 元，该商品还剩余%d件\n", product.name, count, float64(count)*product.price, product.stock)
		}
	} else {
		fmt.Println("商品不存在,无法购买")
	}
}
