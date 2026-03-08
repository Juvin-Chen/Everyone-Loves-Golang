// 测试file_handler中读取文件功能是否正确实现
// 1.正常日志文件：验证核心功能；
// 2.空文件：验证边界条件；
// 3.文件不存在：验证错误处理；
// 4.含空行文件：验证空行过滤逻辑；
// 5.大文件：验证性能和批量读取正确性。

package fileop

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestReadLines_NormalCase 测试正常日志文件读取（含多行有效内容）
func TestReadLines_NormalCase(t *testing.T) {
	// 1. 创建临时测试文件
	tempDir := t.TempDir() // 测试结束后自动删除临时目录
	tempFile := filepath.Join(tempDir, "normal.log")
	// 写入模拟日志内容
	testContent := `2024-02-14 10:02:00 [ERROR] Timeout waiting for service. ErrorCode: 504, IP: 192.168.1.20
					2024-02-14 10:03:00 [INFO] Service started successfully. IP: 192.168.1.21
					2024-02-14 10:04:00 [WARN] Low memory detected. IP: 192.168.1.22`
	if err := os.WriteFile(tempFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建正常测试文件失败: %v", err)
	}

	// 2. 调用ReadLines读取文件
	lines, err := ReadLines(tempFile)
	if err != nil {
		t.Fatalf("ReadLines执行失败: %v", err)
	}

	// 3. 验证结果
	expectedLines := []string{
		"2024-02-14 10:02:00 [ERROR] Timeout waiting for service. ErrorCode: 504, IP: 192.168.1.20",
		"2024-02-14 10:03:00 [INFO] Service started successfully. IP: 192.168.1.21",
		"2024-02-14 10:04:00 [WARN] Low memory detected. IP: 192.168.1.22",
	}

	if len(lines) != len(expectedLines) {
		t.Errorf("行数不符: 预期%d行, 实际%d行", len(expectedLines), len(lines))
	}
	for i := range expectedLines {
		if lines[i] != expectedLines[i] {
			t.Errorf("第%d行内容不符: 预期[%s], 实际[%s]", i+1, expectedLines[i], lines[i])
		}
	}
}

// TestReadLines_EmptyFile 测试空文件场景
func TestReadLines_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "empty.log")
	// 创建空文件
	if err := os.WriteFile(tempFile, []byte(""), 0644); err != nil {
		t.Fatalf("创建空测试文件失败: %v", err)
	}

	lines, err := ReadLines(tempFile)
	if err != nil {
		t.Fatalf("ReadLines读取空文件失败: %v", err)
	}

	if len(lines) != 0 {
		t.Errorf("空文件应返回空切片: 实际返回%d行", len(lines))
	}
}

// TestReadLines_FileNotExist 测试文件不存在场景
func TestReadLines_FileNotExist(t *testing.T) {
	// 故意用不存在的文件路径
	nonExistFile := filepath.Join(t.TempDir(), "not_exist.log")

	lines, err := ReadLines(nonExistFile)
	// 验证是否返回错误
	if err == nil {
		t.Error("文件不存在时应返回错误，但未返回")
	}
	// 验证返回的切片应为nil
	if lines != nil {
		t.Error("文件不存在时返回的切片应为nil")
	}
}

// TestReadLines_WithEmptyLines 测试含空行的日志文件（空行应被过滤）
func TestReadLines_WithEmptyLines(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "with_empty.log")
	// 含空行的测试内容
	testContent := `
					2024-02-14 10:02:00 [ERROR] Timeout waiting for service. ErrorCode: 504, IP: 192.168.1.20

					2024-02-14 10:03:00 [INFO] Service started successfully. IP: 192.168.1.21
					`
	if err := os.WriteFile(tempFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建含空行测试文件失败: %v", err)
	}

	lines, err := ReadLines(tempFile)
	if err != nil {
		t.Fatalf("ReadLines读取含空行文件失败: %v", err)
	}

	expectedLines := []string{
		"2024-02-14 10:02:00 [ERROR] Timeout waiting for service. ErrorCode: 504, IP: 192.168.1.20",
		"2024-02-14 10:03:00 [INFO] Service started successfully. IP: 192.168.1.21",
	}

	if len(lines) != len(expectedLines) {
		t.Errorf("含空行文件过滤后行数不符: 预期%d行, 实际%d行", len(expectedLines), len(lines))
	}
}

// TestReadLines_LargeFile 测试大文件读取（验证性能和正确性）
func TestReadLines_LargeFile(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "large.log")

	// 生成1000行测试日志
	var largeContent string
	for i := 0; i < 1000; i++ {
		largeContent += string(fmt.Sprintf("2024-02-14 10:%02d:00 [INFO] Request %d processed. IP: 192.168.1.%d\n", i/60, i, i%255))
	}
	if err := os.WriteFile(tempFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("创建大测试文件失败: %v", err)
	}

	lines, err := ReadLines(tempFile)
	if err != nil {
		t.Fatalf("ReadLines读取大文件失败: %v", err)
	}

	if len(lines) != 1000 {
		t.Errorf("大文件读取行数不符: 预期1000行, 实际%d行", len(lines))
	}
	// 验证最后一行内容
	lastLine := lines[999]
	expectedLastLine := "2024-02-14 10:16:00 [INFO] Request 999 processed. IP: 192.168.1.239"
	if lastLine != expectedLastLine {
		t.Errorf("大文件最后一行内容不符: 预期[%s], 实际[%s]", expectedLastLine, lastLine)
	}
}
