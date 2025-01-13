package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckImport(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	// 创建测试用的 .golangci.yml 文件
	configContent := []byte(`check_import:
  - "github.com/forbidden/pkg"
  - "internal/forbidden"
`)
	if err := os.WriteFile(filepath.Join(tmpDir, ".golangci.yml"), configContent, 0644); err != nil {
		t.Fatalf("创建配置文件失败: %v", err)
	}

	// 创建测试用的 Go 文件
	testFileContent := []byte(`package test

import (
	"fmt"
	"github.com/forbidden/pkg/subpkg"
	"internal/forbidden/utils"
)

func main() {
	fmt.Println("test")
}
`)
	if err := os.WriteFile(filepath.Join(tmpDir, "test.go"), testFileContent, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 切换到临时目录
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("切换目录失败: %v", err)
	}
	defer os.Chdir(oldDir)

	// 读取配置
	config, err := readConfig()
	if err != nil {
		t.Fatalf("读取配置失败: %v", err)
	}

	// 检查文件
	errs := checkFile("test.go", config.CheckImport)
	if len(errs) != 2 {
		t.Errorf("期望发现 2 个错误，实际发现 %d 个", len(errs))
	}
}

func TestNoConfigFile(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	// 切换到临时目录
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("切换目录失败: %v", err)
	}
	defer os.Chdir(oldDir)

	// 读取配置
	_, err = readConfig()
	if !os.IsNotExist(err) {
		t.Errorf("期望获得文件不存在错误，但得到: %v", err)
	}
}
