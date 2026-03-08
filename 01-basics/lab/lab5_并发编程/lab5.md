# 🧪 实验 5: Go 并发编程实战 (Concurrency)

**实验目标**：从“线性思维”转变为“并发思维”。掌握 Goroutine 的生命周期管理、Channel 的正确通信姿势、以及如何避免死锁和数据竞争。

------

## 第一部分：基础肌肉记忆 (Tasks 1-3)

### 🧪 任务 1：WaitGroup 的“点名”艺术 (Synchronizing Goroutines)

**痛点**：新手最容易犯的错就是主协程跑完了，子协程还在半路被杀掉了。

**📝 场景：并发下载器模拟**

1. **定义任务**：编写一个函数 `download(url string)`，打印 "Start downloading [url]..."，随机睡眠 1-2 秒，然后打印 "Finish [url]"。
2. **并发执行**：
   - 在 `main` 中定义一个 URL 切片（3-5 个地址）。
   - 遍历切片，为每个 URL 启动一个 Goroutine 去下载。
3. **同步要求**：
   - 使用 `sync.WaitGroup` 确保所有下载任务完成后，主程序才打印 "All downloads completed"。
   - **关键点**：请在循环中正确地传递 loop variable（循环变量）给协程，避免所有协程都下载同一个 URL 的经典坑。

------

### 🧪 任务 2：Channel 的“接力赛” (Unbuffered Channels)

**痛点**：理解无缓冲通道的“阻塞”特性——不见不散。

**📝 场景：乒乓球比赛 (Ping-Pong)**

1. **定义通道**：创建一个无缓冲的 `chan string`。
2. **选手 A**：
   - 启动一个协程，循环 5 次。
   - 每次向通道发送 "Ping"。
   - **发送后**，打印 "A sent Ping"。
3. **选手 B**：
   - 启动另一个协程，循环 5 次。
   - 每次从通道接收数据。
   - **接收后**，打印 "B received: [data]"。
4. **观察顺序**：
   - 运行程序，观察打印日志。是 A 全部发完 B 再收？还是 A 发一个，B 收一个？
   - *思考*：为什么无缓冲通道能起到“同步”的作用？

------

### 🧪 任务 3：Select 的“超时控制” (Timeout Pattern)

**痛点**：网络请求如果卡死了怎么办？Select 是最好的解药。

**📝 场景：不可靠的服务器**

1. **模拟接口**：
   - 编写函数 `rpcCall(c chan string)`。
   - 在函数内睡眠 2 秒（模拟慢速网络），然后发送 "Response Data" 到通道。
2. **主控逻辑**：
   - 在 `main` 中创建通道并启动 `rpcCall`。
   - 使用 `select` 监听该通道。
   - 添加一个 `time.After(1 * time.Second)` 的超时 case。
3. **测试**：
   - 情况 A：超时设为 1 秒（小于 RPC 耗时），预期输出 "Timeout! Request canceled"。
   - 情况 B：超时设为 3 秒（大于 RPC 耗时），预期输出 "Received: Response Data"。

------

## 第二部分：避坑与进阶 (Tasks 4-5)

### 🧪 任务 4：数据竞争与 Mutex (Data Race Safety)

**痛点**：多个协程改同一个变量，结果不对了。

**📝 场景：银行账户抢红包**

1. **不安全版**：
   - 定义全局变量 `balance = 0`。
   - 启动 1000 个协程，每个协程执行 `balance = balance + 1`。
   - 等待所有协程结束（用 WaitGroup），打印 `balance`。
   - *预期*：结果往往小于 1000（数据竞争）。
2. **安全版**：
   - 引入 `sync.Mutex`。
   - 在累加操作前后加锁 `mu.Lock()` 和解锁 `mu.Unlock()`。
   - 再次运行，确认结果稳定为 1000。

------

### 🧪 任务 5：Context 的“一键止停” (Cancellation)

**痛点**：如何优雅地停止一堆正在干活的协程？

**📝 场景：搜索雷达**

1. **雷达协程**：
   - 编写函数 `search(ctx context.Context, name string)`。
   - 使用 `for + select` 循环：
     - `case <-ctx.Done()`: 打印 "[name] 收到停止指令，正在关闭..."，然后 return。
     - `default`: 打印 "[name] 正在搜索..."，睡眠 500ms。
2. **主控室**：
   - 使用 `context.WithCancel` 创建一个上下文。
   - 启动 3 个雷达协程（"Radar-A", "Radar-B", "Radar-C"），都传入这个 `ctx`。
   - 主程序睡眠 2 秒后，调用 `cancel()` 函数。
   - 睡眠 1 秒观察输出。
3. **结果**：你应该看到调用 `cancel()` 后，三个雷达几乎同时停止工作。

------

## 第三部分：综合实战项目库 (Projects)

为了让你彻底通透，我设计了 **3 个不同模式** 的项目。你可以全做，也可以选 2 个重点突破。

### 🏗️ Project 1: 并发日志处理流水线 (Pipeline Pattern)

**难度**：⭐⭐⭐ **核心点**：Producer-Consumer 模型，Channel 串联。

**需求**：

1. **Generator (生产者)**：
   - 启动一个协程，不断生成模拟日志（字符串），例如 "Log 1", "Log 2", "ERROR: Disk Full", "Log 4"...
   - 发送到通道 `chData`。
2. **Processor (处理者)**：
   - 启动一个协程，从 `chData` 读取日志。
   - 过滤：如果是包含 "ERROR" 的日志，发送到通道 `chErr`；否则丢弃或打印 "Pass"。
3. **Saver (消费者)**：
   - 启动一个协程，从 `chErr` 读取错误日志。
   - 模拟写入文件（打印 "Writing to file: [Error Msg]"）。
4. **退出机制**：
   - 生产者发送 100 条后关闭 `chData`。
   - 处理者读完 `chData` 后关闭 `chErr`。
   - 消费者读完 `chErr` 后通知主程序结束。

------

### 🏗️ Project 2: 限制并发的爬虫 (Worker Pool Pattern)

**难度**：⭐⭐⭐⭐ **核心点**：控制并发数量（Semaphore/Worker Pool），防止把服务器搞崩。

**需求**：

1. **任务池**：
   - 定义一个包含 20 个 URL 的切片（模拟任务）。
   - 创建一个 `jobs` 通道，把这 20 个 URL 塞进去。
2. **工人池 (Worker Pool)**：
   - 启动 **3 个** 工人协程（Worker）。
   - 每个 Worker 从 `jobs` 通道抢任务。
   - 模拟下载：打印 "Worker [ID] started job [URL]"，睡眠 1 秒，打印 "Worker [ID] finished"。
3. **结果收集**：
   - Worker 处理完后，把结果发送到 `results` 通道。
4. **观察**：
   - 你会发现虽然有 20 个任务，但任何时刻只有 3 个任务在同时进行。

------

### 🏗️ Project 3: 简单的并发聊天服务器 (Chat Room)

**难度**：⭐⭐⭐⭐⭐ **核心点**：Select, Map 共享状态, 广播机制。

**需求（模拟版，不涉及真实网络 Socket）**：

1. **Broadcaster (广播台)**：
   - 维护一个 `clients` 映射（`map[User]chan string`），记录所有在线用户和他们的接收通道。
   - 拥有一个 `messages` 通道，用于接收某人发来的消息。
   - 启动一个协程监听 `messages`：一旦有消息，遍历 `clients`，把消息发给所有人的通道。
2. **User (用户)**：
   - 结构体包含 `Name` 和 `UserChannel`。
   - 启动协程模拟用户行为：每隔随机时间，向 `messages` 通道发送一条 "[Name]: Hello!"。
3. **Join/Leave**：
   - 另外两个通道 `entering` 和 `leaving` 用于处理用户上线下线。
   - 广播台需要用 `select` 同时处理：新消息广播、用户加入、用户退出。

------

### 💡 建议

1. **先做 Task 1-3**：这是基础中的基础，如果连 WaitGroup 和 Channel 阻塞都不熟，做项目会非常痛苦。
2. **Project 2 是经典面试题**：Worker Pool 模式在实际工作中用得最多（比如控制数据库并发查询数），强烈推荐优先做这个。
3. **死锁了别怕**：遇到 `fatal error: all goroutines are asleep - deadlock!` 是好事，说明你开始触碰到 Channel 的边界了。检查一下是不是“有发无收”或者“有收无发”。