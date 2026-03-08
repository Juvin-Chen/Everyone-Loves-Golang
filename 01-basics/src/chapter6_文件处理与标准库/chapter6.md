# 📘 第 6 章: 文件处理与标准库 (File I/O & StdLib)

本章聚焦于 Go 语言中处理 IO 操作的核心能力，以及如何使用正则表达式处理文本数据。这是后端开发中最基础也最常用的技能。

## 1. 文件操作 (File Operations)

Go 提供了两个核心包来处理文件：`os`（底层操作）和 `bufio`（带缓冲的高效操作）。

### 1.1 简单读写 (Small Files)

适用于一次性读取或写入小文件。底层封装了 `Open`、`Read/Write`、`Close`。

- **写入**: `os.WriteFile(name, data, perm)`
  - `perm` (权限): 通常用 `0644` (所有者读写，组/其他人只读)。
- **读取**: `os.ReadFile(name)`
  - 返回 `[]byte` 切片，需要转 `string` 使用。
- **清理**: `os.Remove(name)`

### 1.2 缓冲读写 (Large Files)

**面试考点**：处理 GB 级大文件时，必须使用 `bufio`，否则频繁的系统调用 (System Call) 会导致性能极其低下。

- **写入 (`bufio.NewWriter`)**:
  - 数据先写入内存缓冲区。
  - **关键点**: 必须调用 `writer.Flush()`，否则最后一部分数据会丢失在内存里，没写进磁盘。
- **读取 (`bufio.NewScanner`)**:
  - 适用于**逐行读取** (Log 分析等场景)。
  - `scanner.Scan()`: 读一行，返回 `bool`。
  - `scanner.Text()`: 获取内容。

### 1.3 追加模式 (Append)

使用 `os.OpenFile` 并指定 Flag：

- `os.O_APPEND`: 追加写入。
- `os.O_CREATE`: 不存在则创建。
- `os.O_WRONLY`: 只写模式。
- **组合 Flag**: 使用位运算 `|` 连接，例如 `os.O_APPEND|os.O_CREATE|os.O_WRONLY`。

------

## 2. 正则表达式 (Regular Expressions)

Go 的 `regexp` 包提供了强大的文本匹配与处理能力。

### 2.1 核心函数

- **`regexp.MustCompile(pattern)`**:
  - **特点**: 如果正则语法错误，程序直接 **Panic (崩溃)**。
  - **场景**: 适合在全局变量或 `init` 中使用（硬编码的正则），因为这是程序员的锅。
- **`regexp.Compile(pattern)`**:
  - **特点**: 返回 `error`，不崩溃。
  - **场景**: 适合处理用户动态输入的正则。

### 2.2 常用操作

- **Check**: `MatchString` (返回 bool)
- **Find**: `FindAllString` (提取所有匹配项)
- **Replace**: `ReplaceAllString` (脱敏、清理数据)
- **Split**: `Split` (按正则切割字符串)

### 2.3 避坑指南

1. **反引号**: 写正则字符串时建议用 **` `** (反引号)，避免 `` 的二次转义地狱。
2. **性能**: 正则引擎开销较大，简单的字符串包含判断建议用 `strings.Contains` 代替。