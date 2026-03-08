package main

import (
	"errors"
	"fmt"
)

func main() {
	demo1_defer()
	demo2_error()
}

/*------------------------------------------------------------------*/
/*
demo1 Defer(延迟执行)
defer 语句会将函数推迟到当前包围的函数返回之后执行。常用于关闭文件、解锁 Mutex 等
作用
1.用来收尾、释放资源：关闭文件、解锁、关闭数据库连接、关闭网络连接等。
2.保证一定会执行，哪怕函数中间报错、return 了。
*/

func demo1_defer() {
	fmt.Println("1 开始")

	// defer 会被压入栈中 (LIFO - 后进先出)
	defer fmt.Println("2.这是第一个defer (最后执行)")
	defer fmt.Println("3.这是第二个defer (倒数第二执行)")

	fmt.Println("4. 结束")

	// 实际输出顺序：1 -> 4 -> 3 -> 2
}

/*------------------------------------------------------------------*/
/*
demo2 错误处理 (Error Handling)
Go 不使用 try-catch，而是通过多返回值返回 error 对象。习惯写法是 if err != nil。
*/
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

func demo2_error() {
	// 正常情况
	if result, err := divide(10, 2); err == nil {
		fmt.Println("10 / 2 =", result)
	}

	// 错误情况
	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("Error occurred:", err) // 处理错误
	} else {
		fmt.Println("Result:", result)
	}
}

/*------------------------------------------------------------------*/
/*
demo3 go工程实践
Go Modules
Go Modules 是 Go 语言的官方依赖管理工具，自 Go 1.11 版本开始引入，在 Go 1.16 版本成为默认的依赖管理模式。
Go Modules 解决了 Go 语言长期以来在依赖管理方面的痛点，为开发者提供了版本控制、依赖隔离和可重复构建等核心功能。
Go Modules 是一组相关 Go 包的集合，它们被版本化并作为一个独立的单元进行管理。每个模块都有一个明确的版本标识，允许开发者在项目中精确指定所需依赖的版本。

核心概念解析
1.模块（Module）：包含 go.mod 文件的目录树，该文件定义了模块的路径、Go 版本要求和依赖关系。
2.版本（Version）：遵循语义化版本控制（Semantic Versioning）的标识符，格式为 vMAJOR.MINOR.PATCH。
3.依赖图（Dependency Graph）：模块及其所有传递依赖的层次结构，Go 工具会自动解析和维护。

为什么需要 Go Modules？
传统 GOPATH 的问题
在 Go Modules 出现之前，Go 使用 GOPATH 模式，存在以下局限性：
1.工作空间限制：所有项目必须放在 GOPATH 目录下
2.版本管理困难：无法精确控制依赖版本
3.依赖冲突：多个项目可能使用同一依赖的不同版本
4.可重复构建挑战：难以确保不同环境下的构建一致性

给核心结论（记死这一句）
Go Modules = 你的 Go 项目专属「依赖管家」：
1.帮你管项目里要用的第三方包（比如想做 Web 开发要装的 gin 框架）；
2.解决 “不同项目用同一包的不同版本不冲突”“不用把项目塞到 GOPATH 里” 的问题；
3.核心就干 3 件事：记清单（go.mod）、验真假（go.sum）、装 / 更 / 删包（go get/tidy 等命令）。

⭐ 详细用法等可查看docs/关于go语言模块工程化实践的文档
*/
