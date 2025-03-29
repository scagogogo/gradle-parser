// 02_dependencies 展示gradle-parser的依赖提取功能
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/scagogogo/gradle-parser/pkg/api"
	"github.com/scagogogo/gradle-parser/pkg/model"
)

func main() {
	// 硬编码配置参数，根据需要修改
	// MODIFY HERE: 更改以下参数
	filePath := "../sample_files/build.gradle" // 要解析的Gradle文件路径
	showScope := true                          // 是否按范围分组显示依赖
	filter := "org.springframework"            // 过滤依赖（包含指定字符串），留空表示显示所有依赖

	// 提取依赖信息
	fmt.Printf("从文件提取依赖: %s\n", filePath)
	dependencies, err := api.GetDependencies(filePath)
	if err != nil {
		fmt.Printf("提取依赖失败: %v\n", err)
		os.Exit(1)
	}

	// 如果设置了过滤器，过滤依赖
	if filter != "" {
		dependencies = filterDependencies(dependencies, filter)
		fmt.Printf("\n=== 过滤后的依赖（包含 '%s'）===\n", filter)
	} else {
		fmt.Println("\n=== 所有依赖 ===")
	}

	// 按范围分组显示
	if showScope {
		// 按范围分组
		dependencySets := api.DependenciesByScope(dependencies)
		displayDependencyByScope(dependencySets)
	} else {
		// 直接显示
		displayDependencies(dependencies)
	}

	// 统计依赖信息
	fmt.Println("\n=== 依赖统计 ===")
	printDependencyStats(dependencies)
}

// 过滤依赖
func filterDependencies(deps []*model.Dependency, filter string) []*model.Dependency {
	var filtered []*model.Dependency
	for _, dep := range deps {
		// 检查Group、Name或者Version是否包含过滤器字符串
		if strings.Contains(dep.Group, filter) ||
			strings.Contains(dep.Name, filter) ||
			strings.Contains(dep.Version, filter) {
			filtered = append(filtered, dep)
		}
	}
	return filtered
}

// 显示依赖
func displayDependencies(deps []*model.Dependency) {
	if len(deps) == 0 {
		fmt.Println("未找到依赖")
		return
	}

	for i, dep := range deps {
		if dep.Scope != "" {
			fmt.Printf("[%d] %s:%s:%s (%s)\n", i+1, dep.Group, dep.Name, dep.Version, dep.Scope)
		} else {
			fmt.Printf("[%d] %s:%s:%s\n", i+1, dep.Group, dep.Name, dep.Version)
		}
	}
}

// 按范围分组显示依赖
func displayDependencyByScope(sets []*model.DependencySet) {
	if len(sets) == 0 {
		fmt.Println("未找到依赖")
		return
	}

	for _, set := range sets {
		fmt.Printf("\n-- %s (%d个) --\n", set.Scope, len(set.Dependencies))
		for i, dep := range set.Dependencies {
			fmt.Printf("[%d] %s:%s:%s\n", i+1, dep.Group, dep.Name, dep.Version)
		}
	}
}

// 打印依赖统计信息
func printDependencyStats(deps []*model.Dependency) {
	// 统计总数
	fmt.Printf("总依赖数: %d\n", len(deps))

	// 统计不同范围的依赖数量
	scopeCount := make(map[string]int)
	for _, dep := range deps {
		scopeCount[dep.Scope]++
	}

	fmt.Println("依赖范围分布:")
	for scope, count := range scopeCount {
		if scope == "" {
			fmt.Printf("  未指定范围: %d个\n", count)
		} else {
			fmt.Printf("  %s: %d个\n", scope, count)
		}
	}

	// 统计不同Group的依赖数量
	groupCount := make(map[string]int)
	for _, dep := range deps {
		groupCount[dep.Group]++
	}

	fmt.Println("依赖Group分布:")
	for group, count := range groupCount {
		if group == "" {
			fmt.Printf("  未指定Group: %d个\n", count)
		} else {
			fmt.Printf("  %s: %d个\n", group, count)
		}
	}
}
