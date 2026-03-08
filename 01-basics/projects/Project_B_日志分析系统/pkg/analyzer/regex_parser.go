// 封装日志分析包
// Package analyzer 提供日志解析功能，支持从日志行中提取级别、IP等核心信息
// 核心函数为 ParseLog，可解析格式为 "... [LEVEL] ... IP: xxx.xxx.xxx.xxx" 的日志行

package analyzer

import (
	"errors"
	"regexp"
	"strings"
)

// LogEntry 定义了日志的数据结构
type LogEntry struct {
	Level string // INFO, ERROR, WARN
	IP    string // IP 地址
	Msg   string // 日志的具体信息
}

// 预编译正则表达式，避免重复编译
var logRegex = regexp.MustCompile(`\[(INFO|ERROR|WARN)\]\s+(.*?)(?:,\s*)?IP:\s*(\d{1,3}(?:\.\d{1,3}){3})`)

// ParseLog 解析单行日志，成功返回结构体指针，失败返回 error
func ParseLog(line string) (*LogEntry, error) {
	// 返回一个切片，包含完整匹配项以及所有 () 分组捕获的内容
	matches := logRegex.FindStringSubmatch(line)

	// 如果没有匹配成功，或者捕获的分组数量不是4个
	if len(matches) != 4 {
		return nil, errors.New("invalid log format")
	}

	// 提取数据并组装返回
	// matches[0] 是整行匹配到的字符串
	// matches[1] 是第一个括号 (INFO|ERROR|WARN)
	// matches[2] 是第二个括号 (.*?) 提取出的 Msg
	// matches[3] 是第三个括号提取出的 IP
	return &LogEntry{
		Level: matches[1],
		Msg:   strings.TrimSpace(matches[2]), // 去除头尾可能多余的空格
		IP:    matches[3],
	}, nil
}
