// 03_plugins 展示gradle-parser的插件提取功能
package main

import (
	"fmt"
	"os"

	"github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
	// 硬编码配置参数，根据需要修改
	// MODIFY HERE: 更改以下参数
	filePath := "../sample_files/build.gradle" // 要解析的Gradle文件路径
	detectType := true                         // 是否检测项目类型

	// 提取插件信息
	fmt.Printf("从文件提取插件: %s\n", filePath)
	plugins, err := api.GetPlugins(filePath)
	if err != nil {
		fmt.Printf("提取插件失败: %v\n", err)
		os.Exit(1)
	}

	// 显示插件信息
	fmt.Println("\n=== 插件列表 ===")
	if len(plugins) == 0 {
		fmt.Println("未找到插件")
	} else {
		for i, plugin := range plugins {
			if plugin.Version != "" {
				fmt.Printf("[%d] %s (版本: %s)\n", i+1, plugin.ID, plugin.Version)
			} else {
				fmt.Printf("[%d] %s\n", i+1, plugin.ID)
			}
		}
	}

	// 如果需要检测项目类型
	if detectType {
		fmt.Println("\n=== 项目类型检测 ===")

		// 检测是否是Android项目
		if api.IsAndroidProject(plugins) {
			fmt.Println("✓ 这是一个Android项目")
		} else {
			fmt.Println("✗ 这不是一个Android项目")
		}

		// 检测是否是Kotlin项目
		if api.IsKotlinProject(plugins) {
			fmt.Println("✓ 这是一个Kotlin项目")
		} else {
			fmt.Println("✗ 这不是一个Kotlin项目")
		}

		// 检测是否是Spring Boot项目
		if api.IsSpringBootProject(plugins) {
			fmt.Println("✓ 这是一个Spring Boot项目")
		} else {
			fmt.Println("✗ 这不是一个Spring Boot项目")
		}
	}

	// 统计插件信息
	fmt.Println("\n=== 插件统计 ===")
	fmt.Printf("总插件数: %d\n", len(plugins))

	// 统计有版本信息的插件数量
	versionedCount := 0
	for _, plugin := range plugins {
		if plugin.Version != "" {
			versionedCount++
		}
	}
	fmt.Printf("有版本信息的插件数: %d\n", versionedCount)
	fmt.Printf("无版本信息的插件数: %d\n", len(plugins)-versionedCount)
}
