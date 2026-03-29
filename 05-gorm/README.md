# GORM 
## 一、核心定义
Go 语言**生态最主流的 ORM 框架**，将数据库表映射为 Go 结构体，**全自动封装 SQL**，无需手写原生查询，专注业务开发。

## 二、核心特性
1. 兼容全主流数据库：MySQL、PostgreSQL、SQLite、SQL Server
2. 极简链式 API：一行代码实现增删改查
3. 自动迁移：结构体变更自动同步表结构
4. 强大关联：一对一/一对多/多对多关联查询
5. 必备能力：事务、软删除、分页、钩子函数
6. 无侵入、轻量级，性能接近原生 SQL

## 三、核心优势
- **类型安全**：编译期校验，杜绝 SQL 语法错误
- **开发高效**：减少 80% 数据库层冗余代码
- **上手零成本**：API 设计简洁，文档完善
- **框架标配**：Gin/Echo 等 Web 框架首选数据层工具

## 四、快速入门
### 安装
```bash
go get gorm.io/gorm
go get gorm.io/driver/mysql
```

### 极简示例
```go
package main

import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

// 结构体映射数据库表
type User struct {
  gorm.Model // 内置ID/创建时间/更新时间/删除时间
  Name string
  Age  int
}

func main() {
  // 连接数据库
  dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  // 自动创建表
  db.AutoMigrate(&User{})
  // 新增数据
  db.Create(&User{Name: "test", Age: 18})
}
```

## 五、适用场景
Gin Web 项目、Go 后端服务、微服务数据层开发