# Viper 
## 1. 基本定义
Viper 是 **Go 语言最主流的配置管理库**，由 Go 生态知名开发者 spf13 打造，是 Gin、Echo 等 Web 框架的**标配配置工具**，专门解决 Go 项目中配置文件读取、解析、热更新等问题。

## 2. 核心优势
- 支持**多种配置格式**：YAML、JSON、TOML、INI 等（Gin 项目首选 YAML）
- 支持**多来源配置**：本地文件、环境变量、命令行参数、远程配置中心
- 支持**配置热重载**：修改配置文件无需重启项目，自动生效
- 极简 API：一行代码读取/监听配置项
- 零耦合：可无缝集成到任意 Go 项目中

## 3. 核心作用
统一管理项目配置（如服务端口、数据库连接、日志参数、第三方配置等），替代硬编码，让项目配置更灵活、易维护。

## 4. 快速入门
### 安装
```bash
go get github.com/spf13/viper
```

### 极简使用（读取 YAML 配置）
```go
package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")   // 配置文件名
	viper.SetConfigType("yaml")     // 配置文件类型
	viper.AddConfigPath("./")       // 配置文件路径

	// 读取配置
	_ = viper.ReadInConfig()

	// 获取配置项
	port := viper.GetString("server.port")
	fmt.Println("服务端口：", port)
}
```

## 5. 适用场景
- Gin Web 项目配置管理（最常用）
- Go 命令行工具配置
- 微服务/后端项目统一配置中心