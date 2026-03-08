package fileop

import (
	"bufio"
	"os"
)

// 封装文件操作包

// ReadLines 读取文件所有行，返回字符串切片
// 1.使用 os.Open 打开文件，并通过 defer 确保资源释放
// 2.使用 bufio.Scanner 逐行读取，自动处理换行符
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// 检查扫描过程中是否发生错误（如 I/O 错误）
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// WriteToFile 将字节切片写入指定文件，处理底层数据最通用的类型是字节切片 []byte
func WriteToFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		return err
	}

	// 刷新缓冲区，确保所有数据写入磁盘
	return writer.Flush()
}

// 其实 Go 标准库提供了一个极简的封装 os.WriteFile(path, data, 0644)，一行代码就能替代上面这一整段。目前是手动使用 bufio
