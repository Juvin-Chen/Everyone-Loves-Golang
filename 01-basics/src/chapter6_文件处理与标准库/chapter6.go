package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	// demo1_file()
	// demo2_regex()
}

/*------------------------------------------------------------------*/
/*
demo1 文件操作
掌握文件的读写、创建删除，以及最重要的——如何高效处理大文件。

核心库介绍
• os：底层操作（创建、打开、删除），最常用。
• bufio：带缓冲的 I/O，读写大文件必用，能减少磁盘操作次数。
• path/filepath：处理路径（自动兼容 Windows 的 \ 和 Linux 的 /）。
*/

// 1. 简单读写（适合小文件）
// 对应 os.WriteFile 和 os.ReadFile
func basicFileOps() {
	fileName := "test_simple.txt"
	content := []byte("Hello, Go File System!\nThis is a simple write.")

	// 写入文件
	// 0644 是权限代码：所有者读写，其他人只读
	err := os.WriteFile(fileName, content, 0644)
	if err != nil {
		fmt.Println("写入失败", err)
		return
	}
	fmt.Println("1. 文件写入成功")

	// 读取文件
	readData, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}

	fmt.Println("file content:", readData)

	// 清理文件
	os.Remove(fileName)
}

// 2. 缓冲读写（适合大文件/逐行处理）
// 对应 bufio.NewScanner 和 bufio.NewWriter
func bufferedFileOps() {
	fileName := "test_buffered.txt"

	// --- 写入 ---
	// os.Create 创建文件，如果存在则清空
	file, err := os.Create(fileName)
	if err != nil {
		// 触发panic，程序崩溃，打印错误信息
		// 这里用panic是兜底——文件都打不开，后续操作没法做，只能退出
		panic(err) // panic参考chapter7错误处理
	}

	// 重要：函数结束前关闭文件，释放资源
	defer file.Close() // defer参考chapter7

	writer := bufio.NewWriter(file)
	writer.WriteString("Line 1: High Performance\n")
	writer.WriteString("Line 2: Buffered I/O\n")

	// 必须 Flush！否则数据可能还在内存里，没真正写入磁盘
	writer.Flush()
	fmt.Println("3. 缓冲写入完成")

	// --- 逐行读取 ---
	// 打开文件用于读取
	fileRead, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer fileRead.Close()

	// 创建扫描器
	scanner := bufio.NewScanner(fileRead)

	fmt.Println("4. 开始逐行读取:")
	lineCount := 0
	// scanner.Scan() 每次读一行，读完返回 false
	for scanner.Scan() {
		lineCount++
		// scanner.Text() 获取当前行内容
		fmt.Printf("第 %d 行: %s\n", lineCount, scanner.Text())
	}
}

// 3. 追加模式 (Append)
// 对应 os.OpenFile 和 os.O_APPEND
func appendToFile() {
	fileName := "test_buffered.txt"

	// os.O_APPEND: 追加内容
	// os.O_WRONLY: 只写模式
	// os.O_CREATE: 如果不存在则创建
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("打开失败")
		return
	}

	defer file.Close()

	// 这是一种写在一行式的简洁写法
	if _, err := file.WriteString("Line3:This is appended!\n"); err != nil {
		fmt.Println("追加失败", err)
	}
	fmt.Println("5. 追加内容成功")
}

func demo1_file() {
	fmt.Println("=== 文件操作 Demo ===")
	basicFileOps()
	bufferedFileOps()
	appendToFile()
	fmt.Println()
}

/*------------------------------------------------------------------*/
/*
demo2 正则表达式
1. 学会从文本中提取特定格式的数据（如手机号、邮箱、IP地址）。

2. 核心概念
• 包名：regexp。
• Compile vs MustCompile：
    ◦ Compile：返回错误，适合用户输入正则的场景。
    ◦ MustCompile：如果正则写错了直接崩溃（Panic），适合在代码里写死正则的场景（后端最常用）。
• 反引号：写正则字符串时，建议用反引号  ` `，这样就不需要对反斜杠 \ 进行二次转义。

3. 详细介绍
Go 语言中的 regexp 包
Go 语言的标准库提供了 regexp 包，用于处理正则表达式。以下是 regexp 包中常用的函数和方法：
（1）Compile 和 MustCompile
用于编译正则表达式。Compile 返回一个 *Regexp 对象和一个错误，而 MustCompile 在编译失败时会直接 panic。

（2）MatchString
检查字符串是否匹配正则表达式。

（3）FindString 和 FindAllString
用于查找匹配的字符串。FindString 返回第一个匹配项，FindAllString 返回所有匹配项。

（4）ReplaceAllString
用于替换匹配的字符串。

（5）Split
根据正则表达式分割字符串。

4.正则表达式的基本语法
以下是一些常用的正则表达式语法：

.：匹配任意单个字符（除了换行符）。
*：匹配前面的字符 0 次或多次。
+：匹配前面的字符 1 次或多次。
?：匹配前面的字符 0 次或 1 次。
\d：匹配数字字符（等价于 [0-9]）。
\w：匹配字母、数字或下划线（等价于 [a-zA-Z0-9_]）。
\s：匹配空白字符（包括空格、制表符、换行符等）。
[]：匹配括号内的任意一个字符（例如 [abc] 匹配 a、b 或 c）。
^：匹配字符串的开头。
$：匹配字符串的结尾。
*/

func demo2_regex() {
	//1. 检查字符串是否匹配正则表达式
	pattern1 := `^[a-zA-Z0-9]+$`
	regex1 := regexp.MustCompile(pattern1) // 这是一个指针对象

	str1 := "Hello123"
	if regex1.MatchString(str1) {
		fmt.Println("字符串匹配正则表达式")
	} else {
		fmt.Println("字符串不匹配正则表达式")
	}

	//2. 查找匹配的字符串
	pattern2 := `\d+`
	regex2 := regexp.MustCompile(pattern2)

	str2 := "我有 3 个苹果和 5 个香蕉"
	matches := regex2.FindAllString(str2, -1)
	fmt.Println("找到的数字：", matches)

	//3. 替换匹配的字符串
	pattern3 := `\s+`
	regex3 := regexp.MustCompile(pattern3)

	str3 := "Hello    World"
	result := regex3.ReplaceAllString(str3, " ")
	fmt.Println("替换后的字符串：", result)

	//4. 分割字符串
	pattern4 := `,`
	regex4 := regexp.MustCompile(pattern4)

	str4 := "apple,banana,orange"
	parts := regex4.Split(str4, -1)
	fmt.Println("分割后的字符串：", parts)
}

/*
1.性能问题
正则表达式的匹配和替换操作可能会消耗较多资源，尤其是在处理大量数据时。建议在性能敏感的场景下谨慎使用。

2.转义字符
在 Go 语言中，正则表达式中的反斜杠 \ 需要写成 \\，因为反斜杠在字符串中也是转义字符。

3.错误处理
使用 Compile 函数时，务必检查返回的错误，以避免程序崩溃。
*/
