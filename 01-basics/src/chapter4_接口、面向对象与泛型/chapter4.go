package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	// demo0()
}

/*------------------------------------------------------------------*/
/*
demo0 截止目前已学的知识，对指针等变量创建方式的补充，也是一个总结
*/
func demo0() {
	// 方式1：&取已有变量的地址
	v1 := Vertex{1, 2}
	p1 := &v1
	fmt.Println("p1指向的值：", *p1) // {1 2}

	// 方式2：&取字面量的地址，本质上是先创建对象再取地址
	p2 := &Vertex{3, 4}
	fmt.Println("p2指向的值：", *p2) // {3 4}

	// 上述方式优先栈，下述方式优先堆

	// 方式3：new()创建零值指针
	p3 := new(Vertex)
	p3.X = 5
	p3.Y = 6
	fmt.Println("p3指向的值：", *p3) // {5 6}

	// 方式4：new()创建基本类型指针
	numPtr := new(int)
	*numPtr = 10
	fmt.Println("numPtr指向的值：", *numPtr) // 10
}

/*------------------------------------------------------------------*/
/*
demo1 方法 （methods）
- 方法只是一个带 接收者 (Receiver) 的函数。
- 需要重点理解 值接收者 和 指针接收者 的区别。
*/

type Vertex struct { // vertex：点，顶点，端点
	X, Y float64
}

// 1. 值接收者 (Value Receiver)
// 语法：func (接收者变量 接收者类型) 方法名(参数列表) 返回值列表
// 特点：调用时会复制一份 Vertex，方法内修改不会影响原对象
func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// 2. 指针接收者 (Pointer Receiver) -- 推荐常用
// 特点：直接操作内存地址，方法内修改会影响原对象，且避免大对象复制，效率更高
func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func demo1() {
	v := Vertex{3, 4}

	// 调用值接收者方法
	fmt.Println("Before Scale:", v.Abs())

	// 调用指针接收者方法
	// Go 语法糖：虽然 v 是值，但 Go 会自动解释为 (&v).Scale(10)
	v.Scale(10) //  Go 最贴心的「语法糖」—— 自动转换，不用你手动写&取地址

	fmt.Println("After Scale:", v.Abs()) // 50 (原对象被修改了)
}

/*------------------------------------------------------------------*/
/*
demo2 接口
Go 的接口是隐式实现的。
只要一个类型实现了接口要求的所有方法，它就自动实现了该接口，不需要 implements 关键字
*/

// 定义一个接口：只要能“叫(Speak)”，就是一种 Animal
type Animal interface {
	Speak() string
}

type Dog struct {
	Name string
}

// Dog 实现了 Speak 方法 -> Dog 自动实现了 Animal 接口
func (d Dog) Speak() string {
	return "Woof"
}

type Cat struct {
	Name string
}

// Cat 也实现了 Speak 方法
func (c Cat) Speak() string {
	return "Meow!"
}

// 空接口 interface{}：没有任何方法要求，所以任何类型都实现了空接口
// 类似于 Java 的 Object 或 C++ 的 void*
/*
它是一个特殊的接口，没有要求实现任何方法；
Go 里的所有类型（int、string、数组、结构体、指针等）都「自动实现」了这个空接口；
因此，任何类型的变量都可以赋值给 interface{} 类型的参数 / 变量。
你可以把 interface{} 理解成一个「万能容器」—— 不管你往里面装什么类型的东西，它都能装下，还能记住这个东西原本的类型和值。
*/
func printAny(v interface{}) {
	fmt.Println("Type %T, Value:%v\n", v, v)
}

func demo2() {
	// 多态的使用
	var animals []Animal
	animals = append(animals, Dog{"Buddy"})
	animals = append(animals, Cat{"Kitty"})

	for _, a := range animals {
		fmt.Println(a.Speak())
	}

	// 空接口的使用
	printAny(100)
	printAny("hello")
	printAny(1.09)
	printAny(Dog{"旺旺"})
}

/*------------------------------------------------------------------*/
/*
demo3 Go语言类型转换
类型转换用于将一种数据类型的变量转换为另外一种类型的变量。
Go 语言类型转换基本格式如下：
type_name(expression)

将整型转换为浮点型：
var a int = 10
var b float64 = float64(a)

字符串类型转换
将一个字符串转换成另一个类型，可以使用以下语法：
var str string = "10"
var num int
num, _ = strconv.Atoi(str)
注意，strconv.Atoi 函数返回两个值，第一个是转换后的整型值，第二个是可能发生的错误，我们可以使用空白标识符 _ 来忽略这个错误
*/

func demo3() {
	// 1.字符串转换为整数
	str := "123"
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("转换错误：", err)
	} else {
		fmt.Println("字符串 '%s' 转换为整数为：%d\n", str, num)
	}

	// 2.整数转换为字符串
	num2 := 123
	str2 := strconv.Itoa(num2)
	fmt.Println("整数 %d  转换为字符串为：'%s'\n", num2, str2)

	// 3.字符串转换为浮点数
	str3 := "3.14"
	num3, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("转换错误:", err)
	} else {
		fmt.Printf("字符串 '%s' 转为浮点型为：%f\n", str3, num3)
	}

	// 4.浮点数转换为字符串
	num4 := 3.14
	str4 := strconv.FormatFloat(num4, 'f', 2, 64) //
	fmt.Printf("浮点数 %f 转为字符串为：'%s'\n", num4, str4)
}

/*------------------------------------------------------------------*/
/*
demo4 接口类型转换
接口类型转换有两种情况：类型断言和类型转换。

1.类型断言
类型断言用于将接口类型转换为指定类型，其语法为：
value.(type)  或者  value.(T)
- 其中 value 是接口类型的变量，type 或 T 是要转换成的类型。
- 如果类型断言成功，它将返回转换后的值和一个布尔值，表示转换是否成功。
- 其实就是相当于在做一个类型判断

2.类型转换
类型转换用于将一个接口类型的值转换为另一个接口类型，其语法为：
T(value)
T 是目标接口类型，value 是要转换的值。
在类型转换中，我们必须保证要转换的值和目标接口类型之间是兼容的，否则编译器会报错。
*/

type Writer interface {
	// []byte是「字节切片」，Go 处理所有数据（文本 / 二进制）的通用容器，可修改
	Write([]byte) (int, error)
}

// 实现 Writer 接口的结构体 StringWriter
type StringWriter struct {
	str string
}

// 实现 Write 方法
func (sw *StringWriter) Write(data []byte) (int, error) {
	// string本身是不可改的
	// string的「不可改」：是「字符串内容不可改」，不是「字符串变量不可重新赋值」；
	// sw.str += ...：是创建新字符串，然后把sw.str这个变量指向新字符串 —— 没有修改任何原字符串，完全符合string不可改的规则。
	sw.str += string(data)
	return len(data), nil
}

func function_write_show_demo() {
	// 1. 创建StringWriter指针，赋值给Writer接口
	var w Writer = &StringWriter{}

	// 2. 准备要写入的字符串，转成[]byte
	inputStr := "你好，Go！"                               // 解释一下，中英文所占用的字节数量不一样
	inputBytes := []byte(inputStr)                     // string转[]byte
	fmt.Println("要写入的[]byte长度（字节数）：", len(inputBytes)) // 输出：12（"你好，Go！"的UTF-8字节数）

	// 3. 调用Write方法，传入[]byte
	n, err := w.Write(inputBytes)
	if err != nil {
		fmt.Println("写入失败：", err)
		return
	}

	// 4. 查看结果
	fmt.Println("成功写入的字节数：", n)                               // 输出：12
	fmt.Println("StringWriter的str字段：", w.(*StringWriter).str) // 输出：你好，Go！
}

func demo4() {
	var i interface{} = "hello world"
	str, ok := i.(string)
	if ok {
		fmt.Printf("'%s' is a string\n", str)
	} else {
		fmt.Println("conversion failed")
	}
	/*
		以上实例中，我们定义了一个接口类型变量 i，并将它赋值为字符串 "Hello, World"。然后，我们使用类型断言将 i 转换为字符串类型，并将转换后的值赋值给变量 str。
		最后，我们使用 ok 变量检查类型转换是否成功，如果成功，我们打印转换后的字符串；否则，我们打印转换失败的消息。
	*/

	// 创建一个 StringWriter 实例并赋值给 Writer 接口变量
	var w Writer = &StringWriter{}

	// 将 Writer 接口类型转换为 StringWriter 类型
	// 不安全断言（不带 ok），v := 接口.(类型)，断言失败时程序 panic（崩溃），100% 确定接口里的类型
	sw := w.(*StringWriter)

	// 修改 StringWriter 的字段
	sw.str = "Hello, World"

	// 打印 StringWriter 的字段值
	fmt.Println(sw.str)

	fmt.Println("调用函数展示write方法的demo:")
	function_write_show_demo()
}

/*
解析：
1.定义接口和结构体：
	-Writer 接口定义了 Write 方法。
	-StringWriter 结构体实现了 Write 方法。
2.类型转换：
	-将 StringWriter 实例赋值给 Writer 接口变量 w。
	-使用 w.(*StringWriter) 将 Writer 接口类型转换为 StringWriter 类型。
3.访问字段：
	-修改 StringWriter 的字段 str，并打印其值。
*/

/*------------------------------------------------------------------*/
/*
demo5 Go 语言泛型
核心目标： 学会使用 [T any] 定义泛型函数和结构体，理解类型参数和类型约束，告别为 int 写一遍、为 float 又写一遍的重复劳动。
1. 为什么需要泛型？
在没有泛型之前，如果我们想写一个比较大小的函数，对于 int 和 float64 我们必须写两个函数（如 MaxInt 和 MaxFloat）。 泛型允许我们在函数名后面加一个 [...]，里面定义一个“占位符”（通常叫 T），告诉 Go：“这个类型暂时不确定，等调用的时候再定”。
2. 核心语法
• 类型参数：func 函数名[T 约束](参数 T) ...
• 常用约束：
    ◦ any：任何类型都可以（也就是 interface{} 的别名）。
    ◦ comparable：只接受能用 == 或 != 比较的类型（如数字、字符串，但不能是切片/Map）。
    ◦ 联合约束：使用 | 指定几种特定的类型（如 int | float64）。

类型参数命名约定
通常使用大写字母：T、K、V、E 等
T：表示 Type（类型）
K：表示 Key（键）
V：表示 Value（值）
E：表示 Element（元素）
*/

// ----------------------
// 1. 基础泛型函数 (使用 any)
// ----------------------
func PrintSlice[T any](s []T) {
	for _, v := range s {
		fmt.Printf("%v", v)
	}
	fmt.Println()
}

func Swap[T any](a, b T) (T, T) {
	return b, a
}

// ----------------------
// 2. 约束：comparable
// ----------------------
// FindIndex 查找元素在切片中的位置
// T 必须是 comparable，因为函数体内用到了 v == target 判断
func FindIndex[T comparable](slice []T, target T) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}

// ----------------------
// 3. 联合约束 (Union Constraints)
// ----------------------

// Number 定义一个接口，包含所有数字类型
// 语法：类型1 | 类型2 | ...
type Number interface {
	int | int64 | float64
}

// Add 只能对 Number 类型的数字求和
// 如果传入 string 会报错，因为 string 不在 Number 约束里
func Add[T Number](a, b T) T {
	return a + b
}

// ----------------------
// 4. 泛型结构体 (Generic Structs)
// ----------------------

// Stack 是一个通用的栈
type Stack[T any] struct {
	elements []T
}

// Push 入栈
// 注意：方法接收者也要带 [T]
func (s *Stack[T]) Push(value T) {
	s.elements = append(s.elements, value)
}

// Pop 出栈
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T // 获取 T 类型的零值（比如 int是0，string是""）
		return zero, false
	}
	// 获取最后一个元素
	idx := len(s.elements) - 1
	val := s.elements[idx]
	// 切片缩容
	s.elements = s.elements[:idx]
	return val, true
}

func demo5() {
	fmt.Println("=== 1. 基础泛型 ===")
	// 显式指定类型：PrintSlice[int](...)
	PrintSlice[int]([]int{1, 2, 3})
	// 类型推断：Go 编译器能自动猜出 T 是 string，所以 [string] 可以省略
	PrintSlice([]string{"Hello", "Generics"})

	a, b := Swap("A", "B")
	fmt.Println("Swapped:", a, b)

	fmt.Println("\n=== 2. comparable 约束 ===")
	nums := []int{10, 20, 30, 40}
	fmt.Println("Index of 30:", FindIndex(nums, 30))

	strs := []string{"apple", "banana"}
	fmt.Println("Index of banana:", FindIndex(strs, "banana"))

	fmt.Println("\n=== 3. 联合约束 (Number) ===")
	fmt.Println("Int Sum:", Add(10, 20))
	fmt.Println("Float Sum:", Add(1.1, 2.2))
	// fmt.Println(Add("a", "b")) // 编译报错！string 不满足 Number 约束

	fmt.Println("\n=== 4. 泛型结构体 (Stack) ===")
	// 创建一个 int 类型的栈
	intStack := Stack[int]{}
	intStack.Push(100)
	intStack.Push(200)
	val, ok := intStack.Pop()
	fmt.Printf("Pop Int: %v (Success: %v)\n", val, ok)

	// 创建一个 string 类型的栈
	strStack := Stack[string]{}
	strStack.Push("Go")
	valStr, _ := strStack.Pop()
	fmt.Printf("Pop String: %v\n", valStr)
}

/*------------------------------------------------------------------*/
/*
/*
demo6 Go 结构体嵌入
在面向对象编程（OOP）中，继承是一种机制，允许一个类（子类）从另一个类（父类）继承属性和方法。通过继承，子类可以复用父类的代码，并且可以在不修改父类的情况下扩展或修改其行为。
Go 语言并不是一种传统的面向对象编程语言，它没有类和继承的概念。
Go 使用结构体（struct）和接口（interface）来实现类似的功能。
重点： Go 语言没有继承（Inheritance），只有 组合（Composition）。我们通过“匿名结构体嵌入”来实现类似继承的代码复用效果。
*/

type User struct {
	Name string
	Age  int
}

func (u *User) SayHello() {
	fmt.Printf("Hi,l'm %s, %d years old.", u.Name, u.Age)
}

// Admin 结构体 "嵌入" 了 User
// 这就是 Go 的组合方式。Admin 自动拥有了 User 的字段和方法。
type Admin struct {
	User  // 匿名嵌入，没有字段名，只有类型名
	Level string
}

func demo6() {
	admin := Admin{
		User: User{
			Name: "AdminOp",
			Age:  30,
		},
		Level: "Super",
	}

	// 1. 可以直接访问 User 的字段（语法糖）
	fmt.Println("Admin Name:", admin.Name) // 等同于 admin.User.Name

	// 2. 可以直接调用 User 的方法
	admin.SayHello()

	fmt.Println("Admin Level:", admin.Level)
}
