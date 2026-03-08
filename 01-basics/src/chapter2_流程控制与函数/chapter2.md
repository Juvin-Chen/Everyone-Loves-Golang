# 📘 第 2 章: Go 语言核心逻辑与流程控制

## 1. 条件语句 (Conditionals)

Go 的条件语句是对 C++ 的一次“现代化重构”，去掉了许多历史包袱（如 switch 的穿透陷阱）。

### 1.1 Switch 的彻底进化

Go 的 `switch` 是这一节的明星，它比 C++ 强大且安全得多。

| **特性**     | **C++ switch**                     | **Go switch**                                  |
| ------------ | ---------------------------------- | ---------------------------------------------- |
| **Break**    | 必须手动写，否则穿透 (Fallthrough) | **默认自动 Break**，匹配即停                   |
| **强制穿透** | 无                                 | 使用 `fallthrough` 关键字                      |
| **条件类型** | 仅限整数、枚举 (Integral types)    | **任意可比较类型** (String, Int, Interface...) |
| **表达式**   | 必须是编译期常量                   | 可以是运行时逻辑 (如 `case score >= 90`)       |

**代码实战：**

Go

```
// 1. 字符串匹配 (C++ 做不到)
color := "red"
switch color {
case "red":
    fmt.Println("Stop") // 自动 break
case "green":
    fmt.Println("Go")
}

// 2. 逻辑判断 (替代 if-else if 链)
score := 85
switch { // 默认 switch true
case score >= 90:
    fmt.Println("A")
case score >= 80:
    fmt.Println("B")
}
```

### 1.2 Select (预告)

- **概念**：IO 多路复用，专门用于监听 Channel 通道。
- **注意**：属于并发范畴，将在 **Stage 5** 深入学习。

------

## 2. 循环与标签 (Loops & Labels)

Go 只有一种循环关键字：`for`。它既是 `for` 又是 `while`。

### 2.1 核心机制

- **普通循环**：`for i := 0; i < n; i++ { ... }`
- **While 模式**：`for condition { ... }`
- **死循环**：`for { ... }`

### 2.2 Label (标签)：C++ 痛点的终结者

在 C++ 中跳出多层嵌套循环通常需要设置 `flag` 变量。Go 通过配合 `Label` 完美解决了这个问题。

- **break Label**：直接跳出指定的那一层**外层**循环。
- **continue Label**：直接结束内层，开始指定**外层**循环的下一轮迭代。

**代码实战：**

Go

```
Loop1: // 定义标签
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i == 1 && j == 1 {
                break Loop1 // 直接“瞬移”出最外层
            }
        }
    }
```

------

## 3. 函数：一等公民 (First-Class Citizen)

这是 Go 与 C++ 面向过程部分最大的区别。函数不再只是代码块，它是**值**，是**变量**。

### 3.1 基础特性

- **声明顺序**：无需前向声明 (Forward Declaration)。Go 编译器按包 (Package) 扫描，`main` 写在最上面也没问题。
- **多返回值**：`func foo() (int, error)` 是 Go 的标配。
- **参数简写**：`func add(x, y int)` 等同于 `x int, y int`。

### 3.2 闭包 (Closures) —— 自带“背包”的函数

> **💡 C++ 视角对比**
>
> - **Go 闭包** ≈ **C++ Lambda (带捕获列表)**
> - **区别**：Go 不需要手动写 `[=]` 或 `[&]`，编译器全自动管理。

核心机制：逃逸分析 (Escape Analysis)

你可能会问：“如果函数返回了闭包，闭包引用的局部变量 x 本该在栈上销毁，为什么还能用？”

- **C++**：如果你引用了栈上的局部变量，函数返回后就是**悬垂指针 (Dangling Pointer)**，程序崩溃。
- **Go**：编译器发现 `x` 被闭包带出去了，会自动把 `x` **从栈 (Stack) 搬运到堆 (Heap)** 上。GC (垃圾回收) 会负责后续的清理。

**代码实战：**

Go

```
func createCounter() func() int {
    x := 0 // 本该销毁的局部变量
    return func() int {
        x++ // 被捕获到堆上，拥有了“记忆”
        return x
    }
}
```

### 3.3 函数作为实参 (回调)

函数可以像 `int` 一样传递。这实现了逻辑的“外包”。

Go

```
// op 是一个函数类型的参数
func compute(a, b int, op func(int, int) int) int {
    return op(a, b)
}
```

------

## 4. 变量作用域 (Scope)

- **局部变量**：函数内定义，生命周期随函数栈帧（除非发生逃逸）。
- **全局变量**：函数外定义。
  - **大写开头** (e.g., `Var`)：可被外部包访问 (Public)。
  - **小写开头** (e.g., `var`)：仅当前包可见 (Private)。
- **Shadowing (遮蔽)**：局部变量会优先于同名全局变量。

