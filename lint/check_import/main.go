package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory>\n", os.Args[0])
		os.Exit(1)
	}

	dir := os.Args[1]

	// 读取特定文件内容作为 forbiddenPrefix
	forbiddenPrefix, err := readForbiddenPrefix(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := checkDir(dir, forbiddenPrefix); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// 读取特定文件内容
func readForbiddenPrefix(dir string) (string, error) {
	filePath := filepath.Join(dir, "app/Godeps/GoDeps.json")

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("GoDeps.json file does not exist: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read GoDeps.json file: %v", err)
	}

	var data struct {
		ImportPath string `json:"ImportPath"`
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return "", fmt.Errorf("failed to parse GoDeps.json file: %v", err)
	}

	return data.ImportPath, nil
}

// 遍历指定目录及其子目录中的所有 Go 文件
func checkDir(dir, forbiddenPrefix string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			err = checkFile(path, forbiddenPrefix)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		return nil
	})
}

// 检查文件是否引入了不被允许的包
func checkFile(path, forbiddenPrefix string) error {
	currentPackagePath, err := getPackagePath(path)
	if err != nil {
		return fmt.Errorf("failed to get package path for %s: %v", path, err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %v", path, err)
	}

	for _, importSpec := range file.Imports {
		importPath := strings.Trim(importSpec.Path.Value, "\"")
		// 检查是否是允许的包
		if !isAllowedPackage(currentPackagePath, importPath, forbiddenPrefix) {
			return fmt.Errorf("file %s imports forbidden package %s", path, importPath)
		}
	}
	return nil
}

// 获取当前文件的包路径
func getPackagePath(filename string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("failed to parse file %s: %v", filename, err)
	}
	return file.Name.Name, nil
}

func isAllowedPackage(currentPackagePath, importPath, forbiddenPrefix string) bool {
	// 检查是否是禁止的路径前缀
	if strings.HasPrefix(importPath, forbiddenPrefix) {
		// 检查是否在 commlib 下
		if strings.HasPrefix(currentPackagePath, "commlib") {
			return true
		}
		// 检查是否在当前包的包路径下（去除禁止路径前缀后的相对路径比较）
		currentPathWithoutForbidden := strings.TrimPrefix(currentPackagePath, forbiddenPrefix+"/")
		importPathWithoutFirbidden := strings.TrimPrefix(importPath, forbiddenPrefix+"/")
		commonPrefix := longestCommonPrefix(currentPathWithoutForbidden, importPathWithoutFirbidden)
		return commonPrefix == currentPathWithoutForbidden
	}
	return true
}

// 获取两个字符串的最长公共前缀
func longestCommonPrefix(str1, str2 string) string {
	i := 0
	for i < len(str1) && i < len(str2) && str1[i] == str2[i] {
		i++
	}
	return str1[:i]
}
