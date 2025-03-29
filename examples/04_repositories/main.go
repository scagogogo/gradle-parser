// 04_repositories 展示gradle-parser的仓库提取功能
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/scagogogo/gradle-parser/pkg/api"
	"github.com/scagogogo/gradle-parser/pkg/config"
)

func main() {
	// 硬编码配置参数，根据需要修改
	// MODIFY HERE: 更改以下参数
	filePath := "../sample_files/build.gradle" // 要解析的Gradle文件路径
	checkSpecial := true                       // 是否检查特定仓库的使用

	// 提取仓库信息
	fmt.Printf("从文件提取仓库: %s\n", filePath)
	repos, err := api.GetRepositories(filePath)
	if err != nil {
		fmt.Printf("提取仓库失败: %v\n", err)
		os.Exit(1)
	}

	// 显示仓库信息
	fmt.Println("\n=== 仓库列表 ===")
	if len(repos) == 0 {
		fmt.Println("未找到仓库")
	} else {
		for i, repo := range repos {
			if repo.URL != "" {
				fmt.Printf("[%d] %s (%s)\n", i+1, repo.Name, repo.URL)

				// 显示凭证信息（如果有）
				if repo.Username != "" || repo.Password != "" {
					fmt.Printf("    凭证: %s / %s\n",
						obscureIfNotEmpty(repo.Username),
						maskPassword(repo.Password))
				}
			} else {
				fmt.Printf("[%d] %s\n", i+1, repo.Name)
			}
		}
	}

	// 如果需要检查特定仓库
	if checkSpecial {
		fmt.Println("\n=== 特定仓库检查 ===")
		repoParser := config.NewRepositoryParser()

		// 检查是否使用了JitPack
		if repoParser.HasJitPackRepository(repos) {
			fmt.Println("✓ 使用了JitPack仓库")
		} else {
			fmt.Println("✗ 未使用JitPack仓库")
		}

		// 检查是否使用了自定义仓库
		if repoParser.HasCustomRepository(repos) {
			fmt.Println("✓ 使用了自定义仓库")
		} else {
			fmt.Println("✗ 未使用自定义仓库")
		}

		// 显示默认仓库信息
		fmt.Println("\n常见默认仓库：")
		defaultRepos := repoParser.GetDefaultRepositories()
		for _, repo := range defaultRepos {
			if repo.URL != "" {
				fmt.Printf("- %s (%s)\n", repo.Name, repo.URL)
			} else {
				fmt.Printf("- %s\n", repo.Name)
			}
		}
	}

	// 统计仓库信息
	fmt.Println("\n=== 仓库统计 ===")
	fmt.Printf("总仓库数: %d\n", len(repos))

	// 按类型统计
	typeCount := make(map[string]int)
	for _, repo := range repos {
		typeCount[repo.Type]++
	}

	fmt.Println("仓库类型分布:")
	for repoType, count := range typeCount {
		fmt.Printf("  %s: %d个\n", repoType, count)
	}
}

// 如果不为空，显示为掩码
func obscureIfNotEmpty(value string) string {
	if value == "" {
		return "<未设置>"
	}
	return "******" // 替换为掩码
}

// 掩盖密码，只显示前两位和后两位
func maskPassword(password string) string {
	if password == "" {
		return "<未设置>"
	}

	if len(password) <= 4 {
		return "******"
	}

	masked := password[:2] + strings.Repeat("*", len(password)-4) + password[len(password)-2:]
	return masked
}
