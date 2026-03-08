# 🧪 实验 3.2：RPG 游戏底层数据与内存管理系统

**Experiment 3.2: RPG Data Structure & Memory Management**

**实验背景**： 我们需要为一款简单的 RPG 游戏构建后端逻辑。你需要处理角色的属性（结构体）、角色的固定装备槽（数组）、动态的背包系统（切片）、以及全服的怪物图鉴（Map）。

**核心要求**：

1. **代码位置**：新建 `lab3.2.go`，所有代码写在一个文件中。
2. **严格规范**：除非特定说明，尽量使用 `:=` 短变量声明；严格区分**值传递**和**指针传递**。
3. **零拷贝思想**：在修改角色状态时，禁止产生不必要的结构体副本。

------

### 📝 任务 3.1：定义核心模型 (Structs & Arrays)

**知识点**：`struct` 定义, `[N]Type` 数组的固定性

定义一个名为 `Character` 的结构体，包含以下字段：

1. `Name` (string): 角色名字。
2. `Level` (int): 等级。
3. `HP` (int): 当前血量。
4. `Equipment` ([3]string): **固定长度为 3 的数组**。
   - *设计意图*：分别代表 [0]武器, [1]盔甲, [2]饰品。这是定长的，不能用 append。
5. `Inventory` ([]string): **切片**。
   - *设计意图*：代表背包，长度不固定，捡到东西就往里放。

**编写辅助函数**：

- 编写一个函数 `createHero(name string) *Character`。
  - **注意**：返回的是**指针**。初始化时，HP 默认 100，Level 默认 1，Equipment 初始化为 `{"木棍", "布衣", "无"}` (利用数组字面量初始化)，Inventory 初始化为空切片。

------

### 📝 任务 3.2：状态修改与指针 (Pointers vs Values)

**知识点**：指针传参 `*T` vs 值传参 `T`

这是 Go 语言最容易踩坑的地方。请编写两个“受伤”函数来做对比：

1. **无效函数**：`func damageByValue(c Character, dmg int)`
   - 逻辑：`c.HP = c.HP - dmg`。
2. **有效函数**：`func damageByPtr(c *Character, dmg int)`
   - 逻辑：`c.HP = c.HP - dmg`。
   - **关键点**：Go 语言中访问结构体指针的字段也用 `.`，不需要像 C++ 那样用 `->`，但你要心里清楚这是指针。

**验证步骤（在 main 中）**：

- 创建角色。
- 调用 `damageByValue` 扣除 50 点血，打印 HP（预期：HP 没变，还是 100）。
- 调用 `damageByPtr` 扣除 50 点血，打印 HP（预期：HP 变为 50）。

------

### 📝 任务 3.3：动态背包管理 (Slices: Append & Delete)

**知识点**：切片扩容 `append`，切片截取（删除元素）

编写一个函数 `lootItem(c *Character, item string)`：

- 功能：将 `item` 添加到角色的 `Inventory` 切片中。
- **打印内存变化**：每次添加后，打印当前切片的 `len()` 和 `cap()`。
- **Loop 测试**：在 `main` 中循环调用 5 次，加入 5 个物品。观察 `cap` 是如何翻倍增长的。

编写一个函数 `usePotion(c *Character)`：

- 功能：遍历背包，找到第一个名称为 "Potion" 的物品。

- **逻辑**：

  1. 如果找到，角色 HP + 50。
  2. **核心难点**：将该物品从切片中**删除**。

  - *提示*：使用 `append(slice[:i], slice[i+1:]...)` 语法。

------

### 📝 任务 3.4：怪物图鉴系统 (Maps & Range)

**知识点**：Map 初始化，KV 读写，`range` 遍历

在 `main` 函数外部（或者 main 内部开头），创建一个全局的怪物数据库。

1. **初始化**：定义一个 `map[string]int`，Key 是怪物名字，Value 是怪物的攻击力。
   - 插入数据："Slime": 10, "Wolf": 30, "Dragon": 100。
2. **查找机制**：
   - 编写逻辑，模拟遇到一只 "Wolf"。
   - 从 Map 中查找 "Wolf" 的攻击力。
   - **防御性编程**：同时查找一只不存在的 "Ghost"，利用 Map 的 `value, ok := map[key]` 语法，如果 `ok` 为 false，打印“未知怪物”。
3. **图鉴遍历**：
   - 使用 `range` 遍历这个 Map，打印所有已知的怪物名字和攻击力。

------

### 📝 任务 3.5：综合战斗循环 (Integration Loop)

**知识点**：串联前面所有逻辑，加上流程控制

在 `main` 中把上面所有零件组装成一个 Demo：

1. **开局**：调用 `createHero` 创建你的主角。
2. **捡垃圾**：调用 `lootItem` 让主角捡到 "Potion", "Sword", "Potion", "Gold"。
3. **遭遇战**：
   - 定义一个怪物切片 `enemies := []string{"Slime", "Dragon", "Ghost"}`。
   - 使用 `range` 遍历怪物列表。
   - 在循环内，利用 **任务 3.4** 的 Map 查找怪物攻击力。
   - 如果怪物存在：
     - 打印 "遭遇 [Monster]！"。
     - 调用 `damageByPtr` 让主角扣血。
     - **判断**：如果 HP < 30，自动调用 `usePotion` 喝药回血（触发切片删除逻辑）。
   - 如果怪物不存在（Ghost）：打印 "也就是一阵风罢了..."。
4. **结局**：
   - 打印主角最终的 HP。
   - 使用 `range` 打印主角最终剩下的背包物品。
   - 使用 `for` 循环打印主角的固定装备数组 `Equipment`。