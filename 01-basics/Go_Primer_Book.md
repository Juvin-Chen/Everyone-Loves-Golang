# 📘 Go Primer Book：Go 语言入门书

## 📑 目录 (Table of Contents)
- [第 1 章：包机制、程序入口与工程化](#第-1-章包机制程序入口与工程化)
- [第 2 章：基础语法、变量与数据类型](#第-2-章基础语法变量与数据类型)
- [第 3 章：流程控制与函数](#第-3-章流程控制与函数)
- [第 4 章：核心数据结构](#第-4-章核心数据结构)
- [第 5 章：指针、结构体与内存逃逸](#第-5-章指针结构体与内存逃逸)
- [第 6 章：方法与接口](#第-6-章方法与接口)
- [第 7 章：泛型编程](#第-7-章泛型编程)
- [第 8 章：并发编程](#第-8-章并发编程)
- [第 9 章：Defer 与错误处理](#第-9-章defer-与错误处理)
- [第 10 章：标准库实战](#第-10-章标准库实战)

------

## 第 1 章：包机制、程序入口与工程化

### 1.1 包 (Package) 与 Main 函数规则

Go 的代码组织依赖于包。`main` 包是唯一的“特殊款”，它是程序的启动开关。

- **规则1：** `func main()` 必须定义在 `package main` 中才能编译为可执行文件。
- **规则2：** 同一个目录（包）下只能有一个 `main` 函数。
- **可见性：** 大写字母开头的函数/变量对外导出（公有），小写私有。

```go
// 代码示例 1.1：可执行程序的极简结构
package main // 必须声明为 main 包

import "fmt" // 引入内置的 fmt 工具箱

// 首字母大写，外部包可见
func PublicFunc() {
	fmt.Println("I am public")
}

// 程序的唯一入口
func main() {
	fmt.Println("Hello, 世界") //
}
```

### 1.2 Go Modules 依赖管理

Go 1.16+ 默认的依赖管理工具，取代了传统的 `GOPATH`，解决了版本冲突和路径限制问题。

- **初始化：** `go mod init <模块名>`
- **整理依赖：** `go mod tidy`（自动下载缺失包，删除无用包）

```go
# 代码示例 1.2：从零建项目
mkdir myproject && cd myproject
go mod init github.com/myname/myproject # 生成 go.mod 文件
go get github.com/gin-gonic/gin@v1.9.0  # 下载特定版本的第三方库
go mod tidy                             # 写完代码后整理依赖
```

------

## 第 2 章：基础语法、变量与数据类型

### 2.1 变量与常量声明

Go 提供了灵活的变量声明方式。常量则通过 `const` 定义，结合 `iota` 可实现优雅的枚举。

```go
// 代码示例 2.1：变量、常量与 iota 技巧
package main
import "fmt"

// 1. 全局变量因式分解写法
var (
	globalA int    = 100
	globalB string = "Go"
)

// 2. 常量与 iota 行计数器
const (
	Unknown = iota // 0
	Female         // 1 (继承 iota，自动+1)
	Male           // 2
)

// 3. iota 与位运算结合 (左移)
const (
    Read   = 1 << iota // 1 << 0 = 1 (二进制 001)
    Write              // 1 << 1 = 2 (二进制 010)
    Execute            // 1 << 2 = 4 (二进制 100)
)

func main() {
	// 短变量声明 (仅限函数内)
	score := 95.5
	x, y := 1, 2
	x, y = y, x // 优雅交换

	fmt.Printf("Read: %d, Write: %d\n", Read, Write) // 输出 1, 2
}
```

### 2.2 字符串、Byte 与 Rune

Go 严格区分物理存储和逻辑字符。字符串本质是只读的 UTF-8 字节切片。

- `byte` (uint8)：存储物理字节（英文占1字节）。
- `rune` (int32)：存储逻辑字符（中文占3字节，但转为 rune 后算1个字）。

```go
// 代码示例 2.2：处理中文字符串防乱码
package main
import "fmt"

func main() {
	s := "Hi中"
	fmt.Println("字节长度 len():", len(s)) // 5 (H=1, i=1, 中=3)

	// 错误遍历：按字节拆分会导致中文乱码
	for i := 0; i < len(s); i++ {
		fmt.Printf("%X ", s[i]) // 打印底层 16 进制字节
	}
	fmt.Println()

	// 正确遍历：按 rune (Unicode 码点) 遍历
	for index, char := range s {
		// %c 打印单个字符
		fmt.Printf("位置:%d 字符:%c\n", index, char) //
	}

	// 强转切片计算真实字数
	runeSlice := []rune(s)
	fmt.Println("真实字数:", len(runeSlice)) // 3
}
```

------

## 第 3 章：流程控制与函数

### 3.1 强大的 Switch 与 Label 跳转

Go 的 `switch` 默认匹配即停（不需要 break），而且支持变量/逻辑判断。`Label` 配合 `break/continue` 是跳出多层循环的神器。

```go
// 代码示例 3.1：Switch 与 Label
package main
import "fmt"

func main() {
	// 1. 无表达式 Switch (替代复杂 if-else)
	score := 85
	switch {
	case score >= 90:
		fmt.Println("A")
	case score >= 80:
		fmt.Println("B") // 匹配即停，不会往下穿透
	default:
		fmt.Println("C")
	}

	// 2. 配合 Label 跳出双层循环
OuterLoop:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break OuterLoop // 直接终结双层循环
			}
			fmt.Printf("i=%d, j=%d\n", i, j)
		}
	}
}
```

### 3.2 闭包 (Closures) 与 一等公民

函数可以作为参数传递。闭包能捕获外部变量，将其从栈逃逸到堆上，持久保存状态。

```go
// 代码示例 3.2：闭包实现轻量级状态管理
package main
import "fmt"

// 返回一个 func() int 类型的闭包函数
func createCounter() func() int {
	x := 0 // 这个局部变量被闭包捕获，生命周期延长
	return func() int {
		x++
		return x
	}
}

func main() {
	counterA := createCounter()
	fmt.Println(counterA()) // 1
	fmt.Println(counterA()) // 2

	counterB := createCounter() // 独立的背包
	fmt.Println(counterB()) // 1
}
```

------

## 第 4 章：核心数据结构

### 4.1 数组 vs 切片 (Slice)

数组长度固定且是值传递（全量拷贝）；切片是动态的“数组窗口”，底层指向原数组，是引用传递。

```go
// 代码示例 4.1：切片的截取与扩容
package main
import "fmt"

func main() {
	// 1. 数组 (定长，传值)
	arr :=int{1, 2, 3, 4, 5}

	// 2. 切片 (变长，底层引用)
	// make([]type, len, cap)
	s1 := make([]int, 0, 5)

	// 3. 从数组截取切片 (左闭右开)
	s2 := arr[1:4] //
	s2 = 999    // 修改切片会直接修改底层数组 arr
	fmt.Println("原数组被修改:", arr) //

	// 4. append 追加 (如果超出 cap，会自动扩容)
	s1 = append(s1, 10)
	s1 = append(s1, 20, 30)
	fmt.Printf("s1 len=%d cap=%d val=%v\n", len(s1), cap(s1), s1)
}
```

### 4.2 集合 (Map)

哈希键值对，必须初始化后才能使用。

```go
// 代码示例 4.2：Map 的安全操作
package main
import "fmt"

func main() {
	// 必须用 make 初始化，否则 panic
	dict := make(map[string]int)
	dict["Go"] = 100

	// 安全读取：ok-idiom 模式
	val, ok := dict["Java"]
	if ok {
		fmt.Println("Java score:", val)
	} else {
		fmt.Println("Java 不存在，默认零值为:", val) // 0
	}

	// 删除元素
	delete(dict, "Go")
}
```

------

## 第 5 章：指针、结构体与内存逃逸

### 5.1 指针与结构体组合 (模拟继承)

Go 没有 `class`，通过匿名结构体嵌入实现属性和方法的复用（组合优于继承）。

```go
// 代码示例 5.1：结构体指针与组合
package main
import "fmt"

type User struct {
	Name string
	Age  int
}

func (u *User) SayHello() {
	fmt.Printf("Hi, I am %s\n", u.Name)
}

// 组合：Admin 自动拥有了 User 的所有属性和方法
type Admin struct {
	User  // 匿名嵌入
	Level string
}

func main() {
	// 创建指针实例
	admin := &Admin{
		User:  User{Name: "Alice", Age: 25},
		Level: "Super",
	}

	// 直接访问内层属性和方法
	fmt.Println(admin.Name) // 等同于 admin.User.Name
	admin.SayHello()
}
```

### 5.2 内存逃逸 (Escape Analysis)

变量在栈（Stack）还是堆（Heap）上，完全由编译器根据**生命周期是否超出函数**来决定，与 `new` 关键字无关。

```go
// 代码示例 5.2：触发内存逃逸
package main

// 1. 分配在栈上：生命周期只在函数内
func stackAlloc() int {
	x := 10 // x 留在栈上，函数返回后销毁
	return x
}

// 2. 逃逸到堆上：返回了局部变量的指针，外部还在用它
func heapAlloc() *int {
	y := 20  // y 必须逃逸到堆上，否则返回野指针
	return &y
}
```

------

## 第 6 章：方法与接口

### 6.1 值接收者 vs 指针接收者

方法是带接收者的函数。指针接收者可以直接修改原对象，避免拷贝。

```go
// 代码示例 6.1：指针接收者的高效修改
package main
import "fmt"

type Vertex struct {
	X, Y float64
}

// 指针接收者：直接操作原内存地址
func (v *Vertex) Scale(f float64) {
	v.X *= f
	v.Y *= f
}

func main() {
	v := Vertex{3, 4}
	v.Scale(10) // 语法糖：Go 自动转为 (&v).Scale(10)
	fmt.Println(v) // {30 40}
}
```

### 6.2 接口 (鸭子类型) 与类型断言

只要实现了接口所有方法，就自动隐式实现了该接口。空接口 `interface{}` 可装载任何类型，使用时需通过类型断言还原。

```go
// 代码示例 6.2：多态与类型断言
package main
import "fmt"

type Animal interface {
	Speak() string
}

type Dog struct{}
func (d Dog) Speak() string { return "Woof!" } // 隐式实现

func main() {
	var a Animal = Dog{}
	fmt.Println(a.Speak())

	// 空接口万能容器
	var i interface{} = "hello"

	// 安全类型断言
	str, ok := i.(string)
	if ok {
		fmt.Println("断言成功，长度为:", len(str))
	}

    // Type Switch
    switch v := i.(type) {
    case string:
        fmt.Println("字符串:", v)
    case int:
        fmt.Println("整数:", v)
    }
}
```

------

## 第 7 章：泛型编程

### 7.1 泛型函数与泛型结构体 (Go 1.18+)

泛型允许使用类型参数 `[T constraint]` 编写通用代码。

```go
// 代码示例 7.1：泛型切片查找与泛型栈
package main
import "fmt"

// 1. comparable 约束：T 必须支持 == 判断
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// 2. 泛型结构体
type Stack[T any] struct {
	elements []T
}
func (s *Stack[T]) Push(v T) {
	s.elements = append(s.elements, v)
}

func main() {
	fmt.Println(Contains([]int{1, 2, 3}, 2))       // true
	fmt.Println(Contains([]string{"A", "B"}, "C")) // false

	intStack := Stack[int]{}
	intStack.Push(100)
}
```

------

## 第 8 章：并发编程

### 8.1 Goroutine 与 Channel

“不要通过共享内存来通信，而要通过通信来共享内存”。Channel 负责在协程间安全传递数据。

```go
// 代码示例 8.1：无缓冲 Channel 与同步
package main
import "fmt"

func sum(s []int, c chan int) {
	total := 0
	for _, v := range s {
		total += v
	}
	c <- total // 将结果发进管道 (阻塞直到被拿走)
}

func main() {
	s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int) // 无缓冲，"不见不散"

	go sum(s[:3], c) // 前半段交由子协程 1
	go sum(s[3:], c) // 后半段交由子协程 2

	x, y := <-c, <-c // 从管道取数据 (主协程死等)
	fmt.Println(x, y, x+y)
}
```

### 8.2 Select 多路复用与 WaitGroup

`select` 用于监听多个 Channel；`WaitGroup` 像班长点名，用于等待所有子任务完成。

```go
// 代码示例 8.2：Select 超时与 WaitGroup 点名
package main
import (
	"fmt"
	"sync"
	"time"
)

func main() {
    // 1. WaitGroup 点名模型
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1) // 名单 +1
		go func(id int) {
			defer wg.Done() // 干完活划掉名字
			fmt.Printf("Worker %d finished\n", id)
		}(i)
	}
	wg.Wait() // 班长死等所有任务完成

    // 2. Select 超时控制
	c := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		c <- "result"
	}()

	select {
	case res := <-c:
		fmt.Println("收到:", res)
	case <-time.After(1 * time.Second): // 1秒超时控制
		fmt.Println("超时退出！")
	}
}
```

------

## 第 9 章：Defer 与错误处理

### 9.1 延迟执行与显式异常检查

`defer` 入栈（后进先出），常用于释放锁、关闭文件。Go 通过多返回值显式处理 `error`。

```go
// 代码示例 9.1：Defer 栈与 Error 接口
package main
import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零") // 创建错误对象
	}
	return a / b, nil
}

func main() {
	defer fmt.Println("1. 最后执行 (倒数第一)")
	defer fmt.Println("2. 倒数第二执行")

	res, err := divide(10, 0)
	if err != nil { // 经典的 Go 风格异常检查
		fmt.Println("发生错误:", err)
		return
	}
	fmt.Println("结果:", res)
}
```

------

## 第 10 章：标准库实战

### 10.1 文件操作 (os & bufio)

对于大文件，必须使用 `bufio` 带缓冲读写，提升磁盘 I/O 效率。

```go
// 代码示例 10.1：高效写入与逐行读取大文件
package main
import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fileName := "test.txt"

	// 1. 缓冲写入
	file, _ := os.Create(fileName)
	writer := bufio.NewWriter(file)
	writer.WriteString("Line 1\nLine 2\n")
	writer.Flush() // 必须 Flush，把内存数据刷入磁盘！
	file.Close()

	// 2. 逐行读取 (防 OOM)
	readFile, _ := os.Open(fileName)
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)
	for scanner.Scan() { // 每次取一行
		fmt.Println("读取:", scanner.Text())
	}
	os.Remove(fileName) // 清理文件
}
```

### 10.2 字符串拼接与正则表达式 (regexp)

`strings.Builder` 是高性能拼接首选。硬编码正则使用 `MustCompile`。

```go
// 代码示例 10.2：Builder 拼接与正则提取
package main
import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	// 1. 高性能字符串拼接
	var sb strings.Builder
	sb.Grow(32) // 预分配容量提速
	sb.WriteString("User: ")
	sb.WriteString("Alice")
	fmt.Println(sb.String())

	// 2. 正则表达式提取数据
	// 使用反引号 ` ` 避免转义字符困扰
	regex := regexp.MustCompile(`\d+`)
	text := "Apple: 10, Banana: 25"
	matches := regex.FindAllString(text, -1)
	fmt.Println("提取出的数字:", matches) //
}
```

------

> First stage finished