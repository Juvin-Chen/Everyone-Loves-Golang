/*
ORM（对象关系映射）最大的魅力就是：你只需要写 Go 的结构体（Struct），GORM 会自动帮你在 MySQL 里建好对应的数据表。

1. 定义数据模型 (Model)
我们要定义一个雨伞（Umbrella）模型。GORM 提供了一个非常强大的特性叫 Struct Tags (结构体标签)，它能让你在代码里直接定义数据库表的细节（比如是不是唯一、默认值是什么）。
*/

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 定义雨伞模型 (GORM 会自动将其转换为名为 "umbrellas" 的表)
type Umbrella struct {
	gorm.Model // 嵌入 GORM 默认模型，自动包含: ID (主键), CreatedAt, UpdatedAt, DeletedAt

	// 通过 \`gorm:"..."\` 标签来定制数据库字段属性
	SerialNumber string `gorm:"type:varchar(50);uniqueIndex;not null;comment:雨伞唯一编号"`
	Status       string `gorm:"type:varchar(20);default:'available';comment:状态(available/borrowed/lost)"`
	Location     string `gorm:"type:varchar(100);comment:当前所在存放点"`
}

// 自定义规则
// 只要给结构体写 TableName() string 方法
// GORM 优先使用你 return 的表名
/*
func (Umbrella) TableName() string {
	return "umbrella"
}
*/

// 2. 连接 MySQL 并执行自动迁移
func test1_main() {
	// 1. 配置数据库连接 DSN (Data Source Name)
	// 请把 user, password, dbname 换成你本地 MySQL 实际的账号密码和建好的数据库名
	// 注意：需要提前在 MySQL 里 CREATE DATABASE test_gorm_db;
	dsn := "root:root@tcp(127.0.0.1:3306)/test_gorm_db?charset=utf8mb4&parseTime=True&loc=Local"

	// 2. 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err) // log.Fatalf 数据库连接不上，程序直接崩溃退出
	}
	fmt.Println("数据库连接成功")

	// 3.自动迁移
	// gorm 会自动检查 MySQL 里有没有 umbrellas 表，没有就创建，字段不对就修改
	err = db.AutoMigrate(&Umbrella{})
	if err != nil {
		log.Fatalf("自动建表失败: %v", err)
	}
	fmt.Println("数据表自动迁移完成！请去 MySQL 软件里查看一下。")
}
