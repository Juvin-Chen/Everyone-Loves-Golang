/*
1. 为什么需要路由分组？
在开发中，API 通常按模块或版本分组，例如：
/api/v1/users
/api/v1/products
/api/v2/users

路由分组可以：
共享中间件（如认证、日志）
统一前缀路径
清晰划分功能模块
*/

package main

import "github.com/gin-gonic/gin"

// 2. 路由分组的基本用法
func test5_demo() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", listUsers)
		// v1.POST("/users", createUser)
		// v1.GET("/users/:id", getUser)
	}
	// 创建 v2 分组
	// v2 := r.Group("/api/v2")
	{
		// v2.GET("/users", listUsersV2)
	}
}

// 所有在 v1 分组内注册的路由，路径会自动加上 /api/v1 前缀。分组支持嵌套。

// 3. 为分组添加中间件
// 可以为分组单独添加中间件：
func test5_demo2() {
	// auth := r.Group("/admin")
	// 该组内所有路由都需要认证
	/*
		auth.Use(AuthMiddleware())
		{
			auth.GET("/dashboard", adminDashboard)
			auth.POST("/settings", updateSettings)
		}
	*/

	// 也可以在创建分组时直接传入中间件：
	// auth := r.Group("/admin", AuthMiddleware())
}

// 4. 项目结构建议
// 当项目变大时，将所有路由和处理器放在 main.go 会难以维护。推荐使用以下结构：
/*
myapp/
├── cmd/
│   └── server/
│       └── main.go           # 程序入口
├── internal/
│   ├── handler/              # 处理器（业务逻辑）
│   │   ├── user.go
│   │   └── product.go
│   ├── middleware/           # 中间件
│   │   ├── auth.go
│   │   └── logger.go
│   ├── model/                # 数据模型
│   │   └── user.go
│   ├── repository/           # 数据访问层（可选）
│   └── service/              # 服务层（可选）
├── pkg/                      # 可复用的公共代码
├── config/                   # 配置文件
├── go.mod
└── go.sum
*/
