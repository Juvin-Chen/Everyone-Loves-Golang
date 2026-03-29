/*
在真实的校园共享雨伞平台中，最常见的场景绝不是一次性把几千把伞全部查出来（那样服务器会卡死，手机端也显示不下），而是：
条件筛选 (Where)：比如“只看‘北区食堂’里‘可借用 (available)’的伞”。
1.排序 (Order)：比如“把最新投放的伞排在最前面”。
2.分页 (Limit & Offset)：比如“手机屏幕一页只显示 10 把伞，向下滑动再加载下一页”。
3.在 GORM 中，这些复杂的操作可以通过“链式调用”非常优雅地写出来，就像拼积木一样 db.Where(...).Order(...).Limit(...).Find(...)。
*/

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 这段代码会先批量塞入几把测试伞，然后演示这三种进阶查询。

func test3_Select_main3() {
	dsn := "root:root@tcp(127.0.0.1:3306)/test_gorm_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	db.AutoMigrate(&Umbrella{})
	// 清空 umbrellas 表（重置自增ID，删除所有数据）
	db.Exec("TRUNCATE TABLE umbrellas")

	// 0. 准备测试数据：一次性批量插入 5 把不同状态和位置的伞
	umbrellas := []Umbrella{
		{SerialNumber: "U-101", Status: "available", Location: "图书馆"},
		{SerialNumber: "U-102", Status: "borrowed", Location: "图书馆"},
		{SerialNumber: "U-103", Status: "available", Location: "北区一食堂"},
		{SerialNumber: "U-104", Status: "available", Location: "北区一食堂"},
		{SerialNumber: "U-105", Status: "maintenance", Location: "南大门"}, // maintenance: 维修中
	}
	db.Create(umbrellas)
	fmt.Println("测试数据插入完毕")

	// 1.条件筛选
	fmt.Println("1. 条件筛选")
	var availableUmbrellas []Umbrella
	db.Where("location=? and status = ?", "北区一食堂", "available").Find(&availableUmbrellas)
	for _, u := range availableUmbrellas {
		fmt.Printf("找到符合条件的伞: 编号[%s]\n", u.SerialNumber)
	}

	// 2.排序
	fmt.Println("2.排序")
	var sortedUmbrellas []Umbrella
	// // desc 是降序，asc 是升序
	db.Order("id desc").Find(&sortedUmbrellas)
	fmt.Printf("最新放入系统的伞是: 编号[%s]\n", sortedUmbrellas[0].SerialNumber)

	// 3. 分页 (Limit & Offset)
	fmt.Println("3. 分页查询")
	var pagedUmbrellas []Umbrella
	pagesize := 2   // 每页显示 2 条
	pageNumber := 2 // 当前看第 2 页

	// 核心数学公式：跳过的数据量 = (当前页码 - 1) * 每页数量
	offset := (pageNumber - 1) * pagesize

	// Limit() 限制查几条（我只要几条数据！），Offset() 决定从第几条开始查
	db.Limit(pagesize).Offset(offset).Find(&pagedUmbrellas)

	for _, u := range pagedUmbrellas {
		fmt.Printf("第 %d 页的伞有: 编号[%s] (位置: %s)\n", pageNumber, u.SerialNumber, u.Location)
	}
}

/*
Q & A
Q:在未来的开发中，如果我们有一张表存了 100 万条用户的借阅记录。如果不加 Limit 限制，直接用 db.Find(&records) 把所有数据全查出来塞进切片里，你的 Go 程序可能会面临什么灾难性后果？
A:如果一次性查出 100 万条数据，主要会引发两场“灾难”：
	1.内存撑爆 (OOM - Out Of Memory)：100 万个结构体实例会被全部塞进 Go 程序的内存里。服务器内存一旦耗尽，Go 进程会直接被操作系统强制杀死（Kill），你的整个后端服务就宕机了。
	2.网络与磁盘 I/O 阻塞：数据库要把这 100 万条数据从磁盘读出来，再通过网络传输给你的 Go 程序，这会占用极大的带宽和时间，导致其他用户的正常请求全部卡死超时。
这就是为什么在企业级开发中，只要是查列表，强制要求必须加分页 (Limit)。
*/
