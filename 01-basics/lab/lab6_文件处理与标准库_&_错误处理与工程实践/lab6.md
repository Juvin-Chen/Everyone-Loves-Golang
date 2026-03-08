# 🧪 实验 6: 日志分析系统与工程化实战

**实验目标**：

1. **工程化**：学会使用 `go mod` 管理依赖，划分目录结构（将逻辑拆分为包）。
2. **文件 I/O**：高效处理大文件读取与写入。
3. **正则实战**：从杂乱的文本中提取关键信息。
4. **健壮性**：使用 `defer` 确保资源释放，使用 `error` 优雅处理异常。

------

## 🛠️ 实验前置：初始化工程环境

这一次，我们不能只写一个 `main.go` 了。请严格按照以下步骤构建目录：

1. **创建项目目录**：

   ```go
   mkdir log-analyzer
   cd log-analyzer
   ```

2. **初始化 Go Module**：


   ```go
go mod init github.com/yourname/log-analyzer
   ```

   *(这里的 `yourname` 可以随便写，比如 `juvin`)*

3. **创建目录结构**： 请手动创建以下文件夹和文件：

   ```go
   log-analyzer/
   ├── go.mod        <-- 依赖管理文件 (自动生成)
   ├── main.go       <-- 入口文件
   ├── data/         <-- 存放测试数据
   │   └── server.log
   ├── pkg/          <-- 存放库代码 (Package)
   │   ├── fileop/   <-- 文件操作包
   │   │   └── file_handler.go
   │   └── analyzer/ <-- 日志分析包
   │       └── regex_parser.go
   └── results/      <-- 存放输出结果
   ```

------

## 🧪 任务 1：生成模拟日志 (Data Generation)

在 `data/server.log` 中手动填入（或复制）以下模拟的脏数据。注意里面有正常的日志，也有错误的格式，还有我们要提取的 IP 和 错误码。


```
2024-02-14 10:00:01 [INFO] User login success from IP: 192.168.1.10
2024-02-14 10:00:05 [ERROR] Database connection failed! ErrorCode: 500, IP: 10.0.0.5
This is a garbage line with no format.
2024-02-14 10:01:20 [WARN] High memory usage detected. IP: 127.0.0.1
2024-02-14 10:02:00 [ERROR] Timeout waiting for service. ErrorCode: 504, IP: 192.168.1.20
2024-02-14 10:03:45 [INFO] Logout success.
```

------

## 🧪 任务 2：封装文件操作包 (`pkg/fileop`)

**目标**：练习 `os` 和 `bufio`，以及 `defer`。

**在 `pkg/fileop/file_handler.go` 中编写代码：**

1. **定义包名**：文件第一行必须是 `package fileop`。
2. **编写函数 `ReadLines(path string) ([]string, error)`**：
   - 使用 `os.Open` 打开文件。
   - **关键点**：立即使用 `defer file.Close()`。
   - 使用 `bufio.Scanner` 逐行读取文件。
   - 将每一行 append 到字符串切片中。
   - 如果读取过程中发生错误，返回 `nil, err`。
3. **编写函数 `WriteToFile(path string, content string) error`**：
   - 使用 `os.WriteFile` 或 `bufio.Writer` 将分析结果写入文件。
   - 权限推荐使用 `0644`。

------

## 🧪 任务 3：封装日志分析包 (`pkg/analyzer`)

**目标**：练习 `regexp` 正则提取。

**在 `pkg/analyzer/regex_parser.go` 中编写代码：**

1. **定义包名**：`package analyzer`。

2. **定义结构体**：


   ```go
type LogEntry struct {
    Level string // INFO, ERROR, WARN
    IP    string
    Msg   string
}
   ```

3. **编写函数 `ParseLog(line string) (\*LogEntry, error)`**：

   - **正则需求**：你需要编写一个正则表达式来匹配类似 `[ERROR] ... IP: 10.0.0.1` 的结构。
     - *提示正则*：`\[(INFO|ERROR|WARN)\]` 匹配日志级别。
     - *提示正则*：`IP:\s*(\d{1,3}(?:\.\d{1,3}){3})` 匹配 IP 地址。
   - **逻辑**：
     - 如果这一行不符合日志格式（比如那行 garbage line），返回 `nil, errors.New("invalid log format")`。
     - 如果匹配成功，提取 Level 和 IP，构建 `LogEntry` 结构体指针并返回。

------

## 🧪 任务 4：主程序组装 (`main.go`)

**目标**：练习跨包调用 (Import local packages) 和错误处理流程。

**在 `main.go` 中编写代码：**

1. **Import 包**：


   ```go
import (
    "fmt"
    "github.com/yourname/log-analyzer/pkg/fileop"   // 引入你写的包
    "github.com/yourname/log-analyzer/pkg/analyzer"
)
   ```

2. **流程逻辑**：

   - 调用 `fileop.ReadLines` 读取 `data/server.log`。
   - **错误处理**：如果读取失败，打印错误并退出 (`return`)。
   - 遍历读取到的每一行：
     - 调用 `analyzer.ParseLog`。
     - **错误处理**：如果返回 `err != nil`（说明是垃圾数据），打印 "跳过无效行"。
     - 如果成功，判断日志级别：
       - 如果是 `ERROR`，统计到一个 `errorCount` 变量中，并将该日志的 IP 记录下来。
   - 最后，调用 `fileop.WriteToFile` 将分析报告写入 `results/report.txt`。
     - 报告内容示例："分析完成。发现错误日志 2 条。涉及 IP：10.0.0.5, 192.168.1.20"。

------

## 🚨 避坑指南 (C++ 选手必看)

1. **包的可见性 (Public/Private)**：
   - 在 C++ 里有 `public/private` 关键字。
   - 在 Go 里，**首字母大写 = Public (能被其他包调用)**，**首字母小写 = Private (只能在本包内使用)**。
   - 所以你的函数名必须是 `ReadLines` 而不是 `readLines`，结构体字段必须是 `IP` 而不是 `ip`，否则 `main` 包访问不到！
2. **文件路径问题**：
   - 在 VS Code 中运行 `main.go` 时，默认的工作目录通常是项目根目录。所以文件路径写 `data/server.log` 是对的。
   - 如果你进入 `pkg/fileop` 目录去运行测试，路径就会找不到。
3. **Defer 的执行时机**：
   - `defer` 是在**函数返回前**执行的。如果你在 `for` 循环里打开文件并 `defer close`，这会导致所有文件直到函数结束才关闭，可能会耗尽文件句柄。但在这个实验里，我们在函数级打开文件，是安全的。

------

## 🚀 进阶挑战 (Optional)

如果你觉得上面的太简单，试试这个：

1. **并发分析**： 如果日志文件有 100万行，逐行正则太慢了。试着结合 **Stage 5** 的知识，启动 5 个 Goroutine 并发调用 `ParseLog`，并将结果发送到一个 `Channel` 进行汇总。
2. **JSON 输出**： 查阅 `encoding/json` 标准库，尝试把分析结果序列化成 JSON 格式保存，而不是纯文本。