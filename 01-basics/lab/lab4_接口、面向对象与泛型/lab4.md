# 🧪 实验 4: 接口抽象、组合与泛型编程实战 (Hard Mode)

**实验目标**：彻底重构你的大脑。从 C++ 的“继承+模板”思维，转变为 Go 的“组合+接口+泛型”思维。

------

## 第一部分：基础概念扫盲 (Tasks 1-3)

### 🧪 任务 1：接口的“多态”与“类型断言” (Basic Polymorphism)

**实验目的**：理解 Go 的隐式接口，以及如何从 `interface{}` 把数据“抠”出来。

**📝 场景：万能电子设备充电站**

1. **定义接口**：
   - 定义 `USB` 接口，包含方法 `Start()`, `Stop()`。
2. **实现结构体**：
   - `Phone` (字段 Name)，实现 USB 接口。
   - `Camera` (字段 Name)，实现 USB 接口。
3. **异构集合**：
   - 在 `main` 中创建一个 `[]USB` 切片，放入一个 `Phone` 和一个 `Camera`。
4. **核心逻辑 (断言)**：
   - 编写一个函数 `Computer(devices []USB)`。
   - 遍历切片，调用 `Start()`。
   - **关键点**：在遍历中，使用 `if device, ok := v.(*Phone); ok { ... }` 判断：
     - 如果是 `Phone`，额外打印一条：“检测到手机，开启快速充电模式”。
     - 如果是其他设备，只打印普通连接信息。
   - 最后调用 `Stop()`。

------

### 🧪 任务 2：泛型数据结构 (Generic Data Structures)

**实验目的**：C++ 有 `std::stack`，Go 没有内置栈。利用泛型，自己造一个！

**📝 场景：通用泛型栈**

1. **定义泛型结构体**：
   - `Stack[T any]`，内部包含一个切片 `elements []T`。
2. **实现方法**：
   - `Push(v T)`: 入栈。
   - `Pop() (T, bool)`: 出栈。如果栈空，返回零值和 false；否则返回栈顶元素和 true。
   - `Peek() T`: 查看栈顶元素但不弹出。
3. **测试**：
   - 创建一个 `Stack[int]`，推入 1, 2, 3，验证 Pop 顺序。
   - 创建一个 `Stack[string]`，推入 "Go", "Is", "Good"，验证 Pop 顺序。
   - **思考**：如果不使用泛型，你要为 int 和 string 写两遍代码，感受一下泛型的爽点。

------

### 🧪 任务 3：组合优于继承 (Composition)

**实验目的**：用“匿名嵌入”模拟继承，理解字段提升（Promoted Fields）。

**📝 场景：智能家居系统**

1. **基类组件**：
   - 结构体 `Component`，字段 `ID int`, `Status string`。
   - 方法 `TurnOn()`: 打印 "Component [ID] is On"。
2. **具体设备**：
   - 结构体 `SmartLamp`（智能灯），**嵌入** `Component`，增加字段 `Brightness int`。
   - 结构体 `AirConditioner`（空调），**嵌入** `Component`，增加字段 `Temperature int`。
3. **方法遮蔽 (Shadowing)**：
   - 给 `AirConditioner` 重写 `TurnOn()` 方法：打印 "AirConditioner is cooling down..." (覆盖基类行为)。
4. **测试**：
   - 初始化一个灯和一个空调。
   - 直接访问 `lamp.ID`（验证字段提升）。
   - 调用 `lamp.TurnOn()` 和 `ac.TurnOn()`，观察区别。

------

## 第二部分：C++ 选手进阶陷阱 (Tasks 4-6)

**这部分是重头戏，专门治“懵懵懂懂”。请务必运行代码并阅读报错信息！**

### 🧪 任务 4：空接口与 Type Switch (The Mystery Box)

**实验目的**：`interface{}` 是万能容器，学会用 `switch type` 优雅地处理未知数据。这在解析 JSON 时非常常用。

**📝 场景：解包盲盒**

1. **定义函数**：
   - `func Unpack(box interface{})`。
2. **Type Switch 逻辑**：
   - 使用 `switch v := box.(type)` 语法。
   - `case int`: 打印 "是数字，它的平方是: v*v"。
   - `case string`: 打印 "是文本，内容是: v"。
   - `case []int`: 打印 "是整数切片，长度是: len(v)"。
   - `default`: 打印 "未知类型"。
3. **测试**：
   - 分别传入 `10`, `"Hello"`, `[]int{1,2,3}`, `true` 进行测试。

------

### 🧪 任务 5：接口的“非空”幻觉 (The Nil Interface Trap)

**背景**：Go 的接口包含 `(type, value)`。只有当两者都为 nil，接口才等于 nil。

**📝 任务：制造一个“不等于 nil 的 nil”**

1. **定义**：结构体 `Cat`，接口 `Animal`。
2. **制造陷阱**：
   - `var c *Cat = nil` （明确的空指针）。
   - `var a Animal = c` （把空指针赋给接口）。
3. **观测**：
   - `fmt.Println("c == nil?", c == nil)`
   - `fmt.Println("a == nil?", a == nil)`  <-- **重点观察这里**
4. **修正**：
   - 如果我想让 `a` 真正变成 `nil`，应该怎么赋值？（直接 `var a Animal = nil`）。

------

### 🧪 任务 6：方法集的严格约束 (Pointer Receiver)

**背景**：这是 Go 编译器最严格的地方。

**📝 任务：编译器报错挑战**

1. **定义**：接口 `Mover` (方法 `Move()`)，结构体 `Dog`。
2. **实现**：
   - `func (d *Dog) Move() { ... }` (**指针接收者**)
3. **报错尝试**：
   - `var m Mover = Dog{}` (**尝试把 值 赋给接口**)
   - **预期**：编译报错！
4. **修正**：
   - `var m Mover = &Dog{}` (**把 指针 赋给接口**)
5. **思考**：为什么指针接收者的方法，不能被值类型实现？
   - *提示：接口里的值如果是不可寻址的，怎么调用修改状态的指针方法？*

------

## 第三部分：综合实战项目 (Project)

### 💳 Project: 泛型插件化支付网关 (The Generic Payment Gateway)

**目标**：结合 **接口**（定义规范）、**泛型**（处理结果）、**组合**（复用代码）实现一个模拟企业级的支付系统。

#### 📋 需求文档

**1. 核心接口 (`PaymentStrategy`)** 定义一个接口，规定所有支付插件必须具备的能力：

- `Pay(ctx string, amount float64) string`
  - 输入：上下文ID，金额。
  - 输出：支付结果描述。

**2. 基础组件 (`BasePayment`)**

- 我们发现所有支付方式都需要记录日志。
- 定义结构体 `BasePayment`，实现方法 `Log(msg string)`。
- **组合**：后续的支付实现都要 **嵌入** 这个 `BasePayment`。

**3. 具体插件实现**

- **AliPay**：嵌入 `BasePayment`。实现 `Pay`，逻辑：打印 "调用支付宝API..."，返回 "AliPay Success"。
- **WeChatPay**：嵌入 `BasePayment`。实现 `Pay`，逻辑：打印 "调用微信API..."，返回 "WeChatId Success"。

**4. 泛型结果容器 (`Result[T]`)**

- 支付可能会返回不同的元数据（有时是字符串ID，有时是结构体）。
- 定义泛型结构体 `Result[T any]`：
  - `Success bool`
  - `Data T` (支付成功后返回的具体数据)
  - `Timestamp int64`

**5. 支付网关 (`Gateway`)**

- 这是一个“处理器”。
- 编写函数 `ProcessPayment[T any](method PaymentStrategy, amount float64) Result[T]`。
- **逻辑**：
  1. 调用 `method.Pay()`。
  2. 利用泛型封装返回值（这里为了简化，你可以把 `Pay` 的返回值强行转为 T，或者让 Pay 直接返回 T，视你的设计而定。为了简单起见，建议让 `Pay` 返回 string，然后封装进 `Result[string]`）。
  3. 调用 `method.Log()` 记录流水。

**6. 主程序剧本** 在 `main` 中：
> 该项目在文档中仅给出大致需求与框架，建议根据自己的想法进行扩展（如添加用户类等），完成一个可运行的小型程序。

1. 初始化 `AliPay` 和 `WeChatPay` 对象。
2. 创建一个切片 `strategies := []PaymentStrategy{ali, wechat}`。
3. 遍历切片，多态调用 `Pay`。
4. 使用泛型容器 `Result` 包装结果并打印。

------

### 执行要求

1. **新建文件**：`lab4.go`。
2. **分块编写**：把 `Task 1` 到 `Task 6` 封装成 `task1()` 到 `task6()`。
3. **Project 独立**：`project` 放在单独Go_Stage1_Basics/projects目录下编写。
4. **关键点**：遇到 Task 4, 5, 6 的报错或陷阱时，**不要直接删掉**，用注释把报错原因写出来，证明你懂了。