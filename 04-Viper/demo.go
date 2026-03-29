/*
第一部分：Viper (配置管理)
在写项目时，我们绝对不能把数据库密码、端口号等敏感信息直接写死在 Go 代码里。我们通常会写在一个配置文件里（比如 YAML 格式），然后用 Viper 去读取。

1. 创建一个配置文件 config.yaml
这个文件一般放在项目的根目录下：
server:
  port: 8080
database:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "your_password"
  dbname: "campus_db"

2. 使用 Viper 读取配置
安装 Viper：go get github.com/spf13/viper

在 Go 代码中，我们可以把 YAML 的内容直接映射到一个结构体（Struct）中，这样在代码里调用配置就会非常方便且安全：
*/

package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

func main() {
	// 1. 设置配置文件的名字和路径
	viper.SetConfigName("config") // 文件名不需要带后缀
	viper.SetConfigType("yaml")   // 文件类型
	viper.AddConfigPath(".")      // 在当前目录查找配置文件

	// 2. 读取配置
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	// 3. 将配置解析到结构体中
	var appConfig Config
	if err := viper.Unmarshal(&appConfig); err != nil {
		panic(fmt.Errorf("解析配置失败: %w", err))
	}

	// 4.测试输出
	fmt.Printf("服务启动端口: %d\n", appConfig.Server.Port)
	fmt.Printf("准备连接数据库: %s\n", appConfig.Database.Dbname)
}
