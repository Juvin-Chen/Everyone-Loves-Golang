# 🧪 实验 3：指针、切片与容器综合练习

**实验目标**：彻底搞懂 Go 语言中“数组是值类型”这一核心特性，以及如何利用指针打破这一限制。 **C++ 差异点**：在 C++ 中数组名退化为指针，传参默认是引用效果；在 Go 中，传数组等于完全拷贝。

------

### 🧪 任务 1：指针与数组的“手术刀” (Pointers vs Arrays)

**📝 任务：数组倍增器**

1. **定义数据**：
   - 在 `main` 函数中初始化一个整型数组 `arr := [5]int{10, 20, 30, 40, 50}`。
2. **错误示范 (按值传递)**：
   - 编写一个函数 `doubleArrayByVal(arr [5]int)`。
   - 在函数内部遍历数组，将每个元素乘以 2。
   - 在 `main` 中调用它，然后打印原数组。
   - **预期结果**：原数组纹丝不动。
3. **正确修正 (按指针传递)**：
   - 编写一个函数 `doubleArrayByPtr(arr *[5]int)`。
   - **难点要求**：在函数内部，使用指针解引用或直接索引的方式，将原数组每个元素乘以 2。
   - 在 `main` 中调用它（注意传入的是 `&arr`）。
   - **预期结果**：原数组数值翻倍。
4. **思考题 (写在注释里)**：
   - 在 `doubleArrayByPtr` 中，访问元素时写 `(*arr)[i]` 和 `arr[i]` 有区别吗？Go 编译器在这里做了什么？

------

### 🧪 任务 2：切片的“缝合手术” (Slice Manipulation)

**实验目的**：熟练掌握切片的截取、`append` 的内存分配、`copy` 的深拷贝，理解切片只是底层数组的“窗口”。

**📝 任务：动态数据流处理**

假设你有一组原始数据：`data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}`。

1. **切片删除 (Delete)**：
   - 编写代码，利用切片截取 (`data[a:b]`) 和 `append` 组合，将索引为 5 的元素（即数字 5）从切片中物理删除。
   - **提示公式**：`slice = append(slice[:i], slice[i+1:]...)`。
   - 打印删除后的切片。
2. **切片备份 (Backup & Isolate)**：
   - 这里的陷阱是：如果你直接 `backup := data`，它们共享底层数组。
   - **要求**：创建一个新的切片 `backup`，长度和容量与当前的 `data` 一致。
   - 使用 `copy()` 函数将 `data` 的内容完全复制过去。
   - **验证独立性**：修改 `backup[0] = 999`。
   - 打印 `data` 和 `backup`，证明修改备份没有影响原数据。
3. **扩容观察 (Capacity)**：
   - 对 `data` 进行连续 5 次 `append(data, 100)`。
   - 每次 `append` 后，打印 `len(data)` 和 `cap(data)`。
   - **观察**：容量 (`cap`) 是如何增长的？是每次 +1 还是成倍增长？

------

### 🧪 任务 3：选票统计器 (Map + Range)

**实验目的**：练习 Map 的查找、更新机制，以及 range 遍历的不确定性。

**📝 任务：即时计票系统**

1. **模拟输入**：
   - 定义一个字符串切片，模拟投票箱，里面包含重复的名字： `votes := []string{"Alice", "Bob", "Alice", "Charlie", "Alice", "Bob", "Charlie", "David", "Alice"}`。
2. **统计逻辑**：
   - 定义一个 Map `counter`，Key 是候选人名字，Value 是票数。
   - 使用 `for range` 遍历 `votes` 切片。
   - **逻辑**：如果名字在 Map 中不存在，初始化为 1；如果存在，票数 +1。
3. **输出结果**：
   - 使用 `for range` 遍历 Map，打印格式如下：`Candidate [Name]: [Count] votes`。
4. **进阶挑战 (Optional)**：
   - 在遍历 Map 的同时，找出票数最多的候选人，并在循环结束后宣布 **Winner**。

------

### 🛒 综合实战项目：简易库存管理系统 (Inventory System)

**项目代号**：`Stage3_DataStruct/project/inventory` **目标**：这是一个模拟电商后台的小程序。你需要结合 **结构体**（定义商品）、**Map**（存储数据）、**指针**（修改状态）和 **函数**。

#### 📋 需求文档

1. **数据结构定义**
   - 定义一个结构体 `Product`，模拟商品信息：
     - `ID (int)`: 商品唯一编号
     - `Name (string)`: 商品名称
     - `Price (float64)`: 单价
     - `Stock (int)`: 库存数量 (**核心被修改字段**)
2. **全局/主数据存储**
   - 在 `main` 函数中创建一个核心 Map，名为 `inventory`：
     - **Key**: 商品 ID (`int`)
     - **Value**: 指向商品的指针 (`*Product`)
   - **为什么用指针？** 因为我们需要在“购买”函数中直接修改库存，如果存的是结构体的值 (`Value`)，修改的只是副本。
   - 初始化至少 3 个商品存入 Map（例如：ID 101 是 "iPhone", 库存 5; ID 102 是 "Laptop", 库存 2...）。
3. **功能模块实现 (Functions)**
   - **A. `func showInventory(m map[int]\*Product)`**
     - 遍历 Map，打印当前所有商品的详细信息（ID, Name, Price, Stock）。
     - 注意：Map 遍历是无序的，这很正常。
   - **B. `func addProduct(m map[int]\*Product, id int, name string, price float64, stock int)`**
     - 检查 ID 是否已存在于 Map 中。
     - 如果存在，打印“错误：商品已存在”。
     - 如果不存在，创建一个新的 `Product` 并取其地址存入 Map。
   - **C. `func buyProduct(m map[int]\*Product, id int, count int)`**
     - 这是核心逻辑。
     - 根据 ID 查找商品。
     - 检查 1：商品是否存在？如果不，打印“商品不存在”。
     - 检查 2：库存是否充足 (`Stock >= count`)？如果不，打印“库存不足”。
     - 执行交易：如果可以购买，修改该结构体指针指向的 `Stock` 字段（减去购买数量），并打印“购买成功，花费 xx 元”。
4. **主程序流程 (Main Scenario)**
   - 请在 `main` 中按顺序模拟以下剧本：
     - **初始化**：创建 Map 并预存几个商品。
     - **展示**：调用 `showInventory` 查看初始状态。
     - **进货**：调用 `addProduct` 添加一个新商品（ID 103, "Headphone", ...）。
     - **交易**：调用 `buyProduct` 购买 ID 101 的商品 2 件（库存应减少）。
     - **错误测试**：调用 `buyProduct` 购买 ID 102 的商品 100 件（应提示库存不足）。
     - **最终核算**：再次调用 `showInventory`，确保 ID 101 的库存是正确的（初始值 - 2）。

------

### 💡 避坑指南 (Coding 前必读)

- **Map 初始化**：`var m map[int]*Product` 只是声明了一个 `nil map`，直接赋值会 `panic`。必须用 `m := make(...)` 或者字面量初始化。
- **结构体指针字面量**：存入 map 时，可以这样写：`m[101] = &Product{...}`。
- **切片越界**：在做实验 2 的删除操作时，小心 `append(s[:i], s[i+1:]...)` 中的索引不要越界。

🚀 **开始行动** 请在一个 `lab3.go` 文件中完成所有代码。为了整洁，你可以把前三个实验封装成 `func task1()`, `func task2()`, `func task3()`，并在 `main` 中逐个调用，最后运行 `project()`。