// 日志分析系统入口
// 目前最新的这种写法给这个程序添加了一个要求：并发

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Juvin-Chen/log-analyzer/pkg/analyzer"
	"github.com/Juvin-Chen/log-analyzer/pkg/fileop"
)

// 定义 JSON 报告的结构体
type Report struct {
	TotalErrors int      `json:"total_errors"`
	ErrorIPs    []string `json:"error_ips"`
}

// 这就是典型的后端日常工作流：读取文件 -> 数据清洗与提取 -> 统计聚合 -> 生成报告。
func main() {
	fmt.Println("🚀 开始执行日志分析任务...")

	// 1. 读取日志文件
	lines, err := fileop.ReadLines("data/server.log")
	if err != nil {
		fmt.Printf("❌ 读取日志文件失败: %v\n", err)
		return
	}

	// 2. 初始化统计变量
	var errorCount int
	var errorIPs []string // 声明一个字符串切片，用来存放所有收集到的错误 IP

	// 3. 遍历分析每一行
	for _, line := range lines {
		if strings.TrimSpace(line) == "" { // 容错处理：忽略纯空行，防止无意义的报错
			continue
		}

		entry, err := analyzer.ParseLog(line)
		if err != nil {
			fmt.Printf("⚠️ 跳过无效行: %s\n", line)
			continue
		}

		// 4. 业务逻辑判断：只统计 ERROR 级别的日志
		if entry.Level == "ERROR" {
			errorCount++
			errorIPs = append(errorIPs, entry.IP) // 把提取到的 IP 塞进切片里
		}
	}

	// 将统计数据组装到结构体中
	reportData := Report{
		TotalErrors: errorCount,
		ErrorIPs:    errorIPs,
	}

	// 将结构体序列化为 JSON 格式的 []byte
	// 使用 MarshalIndent 而不是 Marshal，生成的 JSON 带有缩进和换行，更易读
	jsonData, err := json.MarshalIndent(reportData, "", "    ")
	if err != nil {
		fmt.Printf("❌ JSON 序列化失败: %v\n", err)
		return
	}

	// 4. 调用修改后的 WriteToFile 写入文件
	err = fileop.WriteToFile("results/report.json", jsonData)
	if err != nil {
		fmt.Printf("❌ 写入 JSON 报告文件失败: %v\n", err)
		return
	}

	fmt.Println("🎉 任务圆满完成！报告已成功生成到 results/report.json")
	fmt.Println("📄 生成的 JSON 内容预览:")
	// string(jsonData) 可以将字节切片临时转换为字符串打印到终端
	fmt.Println(string(jsonData))
}
