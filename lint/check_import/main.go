package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	CheckImport []string `yaml:"check_import"`
}

func main() {
	// 读取配置文件
	config, err := readConfig()
	if err != nil {
		if os.IsNotExist(err) {
			// 如果配置文件不存在，直接返回成功
			os.Exit(0)
		}
		fmt.Printf("读取配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 扫描当前目录下的所有 Go 文件
	hasError := false
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			if errs := checkFile(path, config.CheckImport); len(errs) > 0 {
				hasError = true
				for _, err := range errs {
					fmt.Printf("%s: %v\n", path, err)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("扫描文件失败: %v\n", err)
		os.Exit(1)
	}

	if hasError {
		os.Exit(1)
	}
}

func readConfig() (*Config, error) {
	data, err := ioutil.ReadFile(".golangci.yml")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

func checkFile(path string, forbiddenImports []string) []error {
	if len(forbiddenImports) == 0 {
		return nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		return []error{fmt.Errorf("解析文件失败: %v", err)}
	}

	var errors []error
	for _, imp := range f.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")
		for _, forbidden := range forbiddenImports {
			if strings.HasPrefix(importPath, forbidden) {
				pos := fset.Position(imp.Pos())
				errors = append(errors, fmt.Errorf("第 %d 行: 禁止导入包 %s", pos.Line, importPath))
			}
		}
	}

	return errors
}
