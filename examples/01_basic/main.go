// 01_basic 展示gradle-parser的基本使用方法
package main

import (
	"fmt"
	"os"

	"github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
	// 使用硬编码的文件路径，可以根据需要修改为您自己的Gradle文件路径
	// MODIFY HERE: 更改此路径以指向您要解析的Gradle文件
	filePath := "../sample_files/build.gradle"

	// 解析Gradle文件
	fmt.Printf("解析文件: %s\n", filePath)
	result, err := api.ParseFile(filePath)
	if err != nil {
		fmt.Printf("解析文件失败: %v\n", err)
		os.Exit(1)
	}

	// 输出解析结果概要
	fmt.Println("=== Gradle解析结果 ===")
	if result.Project != nil {
		fmt.Printf("项目名称: %s\n", result.Project.Name)
		fmt.Printf("项目组: %s\n", result.Project.Group)
		fmt.Printf("项目版本: %s\n", result.Project.Version)
		fmt.Printf("项目描述: %s\n", result.Project.Description)
	}

	// 输出依赖信息
	fmt.Println("\n=== 依赖 ===")
	deps := result.Project.Dependencies
	if len(deps) > 0 {
		for i, dep := range deps {
			fmt.Printf("[%d] %s:%s:%s (%s)\n", i+1, dep.Group, dep.Name, dep.Version, dep.Scope)
		}
	} else {
		fmt.Println("未找到依赖")
	}

	// 输出插件信息
	fmt.Println("\n=== 插件 ===")
	plugins := result.Project.Plugins
	if len(plugins) > 0 {
		for i, plugin := range plugins {
			if plugin.Version != "" {
				fmt.Printf("[%d] %s (版本: %s)\n", i+1, plugin.ID, plugin.Version)
			} else {
				fmt.Printf("[%d] %s\n", i+1, plugin.ID)
			}
		}
	} else {
		fmt.Println("未找到插件")
	}

	// 输出仓库信息
	fmt.Println("\n=== 仓库 ===")
	repos := result.Project.Repositories
	if len(repos) > 0 {
		for i, repo := range repos {
			if repo.URL != "" {
				fmt.Printf("[%d] %s (%s)\n", i+1, repo.Name, repo.URL)
			} else {
				fmt.Printf("[%d] %s\n", i+1, repo.Name)
			}
		}
	} else {
		fmt.Println("未找到仓库")
	}

	// 输出解析时间
	fmt.Printf("\n解析耗时: %s\n", result.ParseTime)

	// 如果有错误，输出错误信息
	if len(result.Errors) > 0 {
		fmt.Println("\n=== 解析错误 ===")
		for i, err := range result.Errors {
			fmt.Printf("[%d] %v\n", i+1, err)
		}
	}
}
