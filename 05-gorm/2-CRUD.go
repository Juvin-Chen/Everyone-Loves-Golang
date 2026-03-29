// CRUD（Create, Read, Update, Delete）是所有后端业务的基石。在 GORM 里，你几乎不需要写一行 SQL 就能完成这些操作。
// 像之前写有关于 Java/C++ 代码的时候会需要手动去拼接 sql 语句，但 gorm 简化了很多

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func test_CRUD_main2() {
	// 更换为刚刚创建的新数据库 test_gorm_db
	dsn := "root:root@tcp(127.0.0.1:3306)/test_gorm_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动建表
	db.AutoMigrate(&Umbrella{})

	fmt.Println("================ 开始演示 CRUD ================")

	// 1.增 - 购入两把新伞
	fmt.Println("1.增加数据")
	u1 := Umbrella{SerialNumber: "UMB-001", Location: "北区一食堂"}
	u2 := Umbrella{SerialNumber: "UMB-002", Location: "图书馆正门"}

	// 注意：这里必须传结构体的指针 (&)
	db.Create(&u1)
	db.Create(&u2)
	fmt.Printf("成功插入两把伞，数据库自动分配的 ID 分别是: %d 和 %d\n", u1.ID, u2.ID)

	// 2.查 - 找伞
	fmt.Println("2.查询数据")
	var findUmbrella Umbrella

	// 用主键 ID 查 (找 ID 为 1 的那把伞)
	// 直接给数字 → db.First(obj, 1) → 默认查主键 ID
	// 对应下面一行（用条件查）前面加 Where → db.Where(条件).First(obj) → 查自定义条件，取第一条
	db.First(&findUmbrella, 1)
	fmt.Printf("通过 ID=1 查到的伞: 编号[%s], 状态[%s]\n", findUmbrella.SerialNumber, findUmbrella.Status)

	// 用条件查 (找编号为 UMB-002 的伞)
	var findUmbrella2 Umbrella
	db.Where("serial_number = ?", "UMB-002").First(&findUmbrella2)
	fmt.Printf("通过编号查到的伞: ID[%d], 位置[%s]\n", findUmbrella2.ID, findUmbrella2.Location)

	// 3. 改 - 有人借走了一把伞
	fmt.Println("3.更新数据")
	// 将 u1 (即 UMB-001) 的状态更新为 "borrowed"
	db.Model(&findUmbrella).Update("Status", "borrowed")
	fmt.Println("UMB-001 已被借走，状态更新完毕！")

	// 4. 删 - UMB-002 损坏报废了
	fmt.Println("4.删除数据")
	db.Delete(&findUmbrella2)
	fmt.Println("UMB-002 已被报废删除！")

	// ----------关于软删除----------
	// “软删除”就是为了保留历史痕迹，保证数据完整性。
	// 如果你写错了一条测试数据，真的想把它从数据库里彻底抹除，GORM 也提供了方法，只需要加上 Unscoped() ：
	// 这次是真删了，连灰都不剩
	// db.Unscoped().Delete(&findUmbrella2)
}
