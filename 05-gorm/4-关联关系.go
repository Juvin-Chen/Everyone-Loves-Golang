/*关联关系 (Associations) —— GORM 的核心灵魂*/

/*
在我们的校园共享雨伞平台中，数据绝对不是孤立的。伞要有记录，记录要有借伞的人。
数据库最强大的地方就在于关系（Relational）。

先重点搞定最常用的一种：一对多 (One-to-Many)。

需求场景：

	一个学生（User）可以借多次伞，产生多条借阅记录（BorrowRecord）。
	我们要实现：只要查出这个学生，就能用一行代码，顺藤摸瓜把他名下的所有借阅记录全带出来。
*/
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 1. 定义学生模型
type Student struct {
	gorm.Model // 包含了ID
	Name       string
	// 【核心魔法】：定义一个切片，告诉 GORM 这个学生有多条借阅记录 (一对多)
	Records []BorrowRecord
}

// 2. 定义借阅记录模型
type BorrowRecord struct {
	gorm.Model
	// 【外键 (Foreign Key)】：GORM 默认会把 "模型名+ID" (StudentID) 当作外键，用来和 Student 表关联
	StudentID  int
	UmbrellaSN string // 借的哪把伞 (存伞的编号)
	Status     string // 状态: "borrowing" (借出中), "returned" (已还)
}

func test4_main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test_gorm_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	// 自动建表：这次建两张表，GORM 会自动在 BorrowRecords 表里加上 student_id 字段
	db.AutoMigrate(&Student{}, &BorrowRecord{})
	fmt.Println("--- 两张关联表创建完毕 ---")

	// 1. 模拟业务：学生注册，并借了两把伞
	fmt.Println("\n--- 1. 插入带有关联的数据 ---")

	// 我们可以在创建学生的同时，直接把他的借阅记录塞进去！GORM 会自动帮我们处理两张表的插入和 ID 绑定。
	newStudent := Student{
		Name: "zhangsan",
		Records: []BorrowRecord{
			{UmbrellaSN: "UMB-999", Status: "borrowing"},
			{UmbrellaSN: "UMB-888", Status: "returned"},
		},
	}
	db.Create(&newStudent)
	fmt.Println("zhangsan注册成功，并自动生成了两条借伞记录！")

	// 2. 模拟业务：查询学生个人主页 (带出历史记录)
	fmt.Println("\n--- 2. 关联查询 (Preload) ---")
	var findStudent Student

	// 【面试必考点】：Preload("Records")
	// 如果不加 Preload，查出来的 Student 里的 Records 切片是空的。
	// 加上 Preload("Records")，GORM 会自动去关联表里把这个学生的所有记录查出来拼装好。
	db.Preload("Records").Where("name=?", "zhangsan").First(&findStudent)

	fmt.Printf("查找到学生: %s\n", findStudent.Name)
	fmt.Printf("该学生一共有 %d 条借阅记录:\n", len(findStudent.Records))

	for i, record := range findStudent.Records {
		fmt.Printf("  记录 %d: 借了伞 [%s], 当前状态 [%s]\n", i+1, record.UmbrellaSN, record.Status)
	}
}
