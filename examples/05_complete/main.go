// 05_complete 展示gradle-parser的完整功能和用法
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/scagogogo/gradle-parser/pkg/api"
	"github.com/scagogogo/gradle-parser/pkg/model"
	"github.com/scagogogo/gradle-parser/pkg/parser"
	"github.com/scagogogo/gradle-parser/pkg/util"
)

func main() {
	// 硬编码配置参数，根据需要修改
	// MODIFY HERE: 更改以下参数
	filePath := "../sample_files/build.gradle" // 要解析的单个Gradle文件路径
	projectDir := "../sample_files"            // 项目目录路径（如果要分析整个项目）
	useProjectMode := true                     // 是否分析整个项目（而不是单个文件）
	jsonOutput := false                        // 是否以JSON格式输出

	// 解析器选项
	skipComments := true      // 是否跳过注释
	collectRawContent := true // 是否收集原始内容
	parsePlugins := true      // 是否解析插件
	parseDependencies := true // 是否解析依赖
	parseRepositories := true // 是否解析仓库
	parseTasks := true        // 是否解析任务

	// 准备解析器选项
	options := api.DefaultOptions()
	options.SkipComments = skipComments
	options.CollectRawContent = collectRawContent
	options.ParsePlugins = parsePlugins
	options.ParseDependencies = parseDependencies
	options.ParseRepositories = parseRepositories
	options.ParseTasks = parseTasks

	// 创建定制解析器
	parser := api.NewParser(options)

	// 判断是解析单个文件还是整个项目
	if useProjectMode {
		// 解析整个项目
		analyzeProject(projectDir, parser, jsonOutput)
	} else {
		// 解析单个文件
		analyzeFile(filePath, parser, jsonOutput)
	}
}

// 解析单个文件
func analyzeFile(filePath string, parser parser.Parser, jsonOutput bool) {
	fmt.Printf("解析文件: %s\n", filePath)

	// 解析文件
	result, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("解析文件失败: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	if jsonOutput {
		// JSON格式输出
		printJsonResult(result)
	} else {
		// 常规格式输出
		printResult(result)
	}
}

// 解析整个项目
func analyzeProject(projectDir string, parser parser.Parser, jsonOutput bool) {
	fmt.Printf("分析项目: %s\n", projectDir)

	// 查找项目根目录
	rootDir, err := util.FindProjectRoot(projectDir)
	if err != nil {
		fmt.Printf("错误: 未找到Gradle项目, 请确保目录包含build.gradle或build.gradle.kts文件\n")
		os.Exit(1)
	}

	fmt.Printf("项目根目录: %s\n", rootDir)

	// 查找所有Gradle文件
	files, err := util.FindGradleFiles(rootDir)
	if err != nil {
		fmt.Printf("错误: 查找Gradle文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("找到%d个Gradle文件\n", len(files))

	// 分析每个文件
	for _, file := range files {
		// 显示相对路径
		relPath, err := filepath.Rel(rootDir, file)
		if err != nil {
			relPath = file
		}

		fmt.Printf("\n==== 分析文件: %s ====\n", relPath)

		// 解析文件
		result, err := parser.ParseFile(file)
		if err != nil {
			fmt.Printf("解析文件失败: %v\n", err)
			continue
		}

		// 输出结果
		if jsonOutput {
			// JSON格式输出
			printJsonResult(result)
		} else {
			// 常规格式输出
			printResult(result)
		}
	}
}

// 以常规格式打印结果
func printResult(result interface{}) {
	// 这里可以实现更复杂的格式化输出
	fmt.Println("=== 解析结果 ===")

	// 将 result 转为 *model.ParseResult
	if parseResult, ok := result.(*model.ParseResult); ok {
		if parseResult.Project != nil {
			// 打印项目信息
			fmt.Println("\n== 项目信息 ==")
			fmt.Printf("名称: %s\n", parseResult.Project.Name)
			fmt.Printf("组: %s\n", parseResult.Project.Group)
			fmt.Printf("版本: %s\n", parseResult.Project.Version)
			fmt.Printf("描述: %s\n", parseResult.Project.Description)

			// 打印依赖信息
			if len(parseResult.Project.Dependencies) > 0 {
				fmt.Printf("\n== 依赖 (%d个) ==\n", len(parseResult.Project.Dependencies))
				for i, dep := range parseResult.Project.Dependencies {
					fmt.Printf("[%d] %s:%s:%s (%s)\n", i+1, dep.Group, dep.Name, dep.Version, dep.Scope)
				}
			}

			// 打印插件信息
			if len(parseResult.Project.Plugins) > 0 {
				fmt.Printf("\n== 插件 (%d个) ==\n", len(parseResult.Project.Plugins))
				for i, plugin := range parseResult.Project.Plugins {
					if plugin.Version != "" {
						fmt.Printf("[%d] %s (版本: %s)\n", i+1, plugin.ID, plugin.Version)
					} else {
						fmt.Printf("[%d] %s\n", i+1, plugin.ID)
					}
				}
			}

			// 打印仓库信息
			if len(parseResult.Project.Repositories) > 0 {
				fmt.Printf("\n== 仓库 (%d个) ==\n", len(parseResult.Project.Repositories))
				for i, repo := range parseResult.Project.Repositories {
					if repo.URL != "" {
						fmt.Printf("[%d] %s (%s)\n", i+1, repo.Name, repo.URL)
					} else {
						fmt.Printf("[%d] %s\n", i+1, repo.Name)
					}
				}
			}
		}

		// 打印解析时间
		fmt.Printf("\n解析耗时: %s\n", parseResult.ParseTime)

		// 如果有错误，输出错误信息
		if len(parseResult.Errors) > 0 {
			fmt.Printf("\n== 解析错误 (%d个) ==\n", len(parseResult.Errors))
			for i, err := range parseResult.Errors {
				fmt.Printf("[%d] %v\n", i+1, err)
			}
		}
	} else {
		fmt.Printf("%+v\n", result)
	}
}

// 以JSON格式打印结果
func printJsonResult(result interface{}) {
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
