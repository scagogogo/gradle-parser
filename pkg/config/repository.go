// Package config 提供Gradle配置解析功能
package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

var (
	// 匹配Maven仓库URL的正则表达式.
	// 例如: maven { url 'https://jitpack.io' }
	// 或者: maven { url = uri("https://maven.aliyun.com/repository/public") }
	mavenUrlRegex = regexp.MustCompile(`url\s*=?\s*(?:uri\()?['"](https?://[^'"]+)['"]`)

	// 匹配Maven仓库名称的正则表达式.
	// 例如: mavenCentral()
	mavenNameRegex = regexp.MustCompile(`(mavenCentral|mavenLocal|jcenter|google)\(\)`)
)

// RepositoryParser 处理Gradle仓库解析.
type RepositoryParser struct{}

// NewRepositoryParser 创建新的仓库解析器.
func NewRepositoryParser() *RepositoryParser {
	return &RepositoryParser{}
}

// ParseRepositoryBlock 解析仓库块.
func (rp *RepositoryParser) ParseRepositoryBlock(block *model.ScriptBlock) ([]*model.Repository, error) {
	if block == nil {
		return nil, fmt.Errorf("仓库块为空")
	}

	repos := make([]*model.Repository, 0)

	// 处理repositories {} 块中的直接值
	for _, value := range block.Values {
		valueStr := fmt.Sprintf("%v", value)

		// 检查是否是预定义仓库名称
		if match := mavenNameRegex.FindStringSubmatch(valueStr); len(match) > 1 {
			repos = append(repos, &model.Repository{
				Name: match[1],
				Type: "maven",
			})
		}
	}

	// 处理子闭包
	for name, closures := range block.Closures {
		switch name {
		case "mavenCentral", "mavenLocal", "jcenter", "google":
			// 预定义的Maven仓库
			repo := &model.Repository{
				Name: name,
				Type: "maven",
			}
			repos = append(repos, repo)

		case "maven":
			// 自定义Maven仓库
			for _, closure := range closures {
				repo := &model.Repository{
					Type: "maven",
				}

				// 寻找URL
				for _, value := range closure.Values {
					valueStr := fmt.Sprintf("%v", value)
					if match := mavenUrlRegex.FindStringSubmatch(valueStr); len(match) > 1 {
						repo.URL = match[1]

						// 从URL推断名称
						parts := strings.Split(match[1], "/")
						if len(parts) > 2 {
							repo.Name = parts[2] // 使用域名作为名称
						} else {
							repo.Name = "custom-maven"
						}
					}
				}

				// 查找凭证信息
				for subName, subClosures := range closure.Closures {
					if subName == "credentials" && len(subClosures) > 0 {
						for key, value := range subClosures[0].Values {
							valueStr := fmt.Sprintf("%v", value)
							switch key {
							case "username":
								repo.Username = strings.Trim(valueStr, "'\"")
							case "password":
								repo.Password = strings.Trim(valueStr, "'\"")
							}
						}
					}
				}

				repos = append(repos, repo)
			}

		case "ivy":
			// Ivy仓库
			for _, closure := range closures {
				repo := &model.Repository{
					Name: "ivy",
					Type: "ivy",
				}

				// 寻找URL
				for _, value := range closure.Values {
					valueStr := fmt.Sprintf("%v", value)
					if match := mavenUrlRegex.FindStringSubmatch(valueStr); len(match) > 1 {
						repo.URL = match[1]
					}
				}

				repos = append(repos, repo)
			}

		case "flatDir":
			// 平面目录仓库
			for _, closure := range closures {
				repo := &model.Repository{
					Name:   "flatDir",
					Type:   "flatDir",
					Config: make(map[string]interface{}),
				}

				// 收集配置
				for key, value := range closure.Values {
					repo.Config[key] = value
				}

				repos = append(repos, repo)
			}
		}
	}

	return repos, nil
}

// ExtractRepositoriesFromText 从原始文本中提取仓库
func (rp *RepositoryParser) ExtractRepositoriesFromText(text string) []*model.Repository {
	repos := make([]*model.Repository, 0)

	// 分析文本中的仓库声明
	lines := strings.Split(text, "\n")
	inRepoBlock := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检查是否进入repositories块
		if strings.Contains(trimmedLine, "repositories") && strings.Contains(trimmedLine, "{") {
			inRepoBlock = true
			continue
		}

		// 检查是否离开repositories块
		if inRepoBlock && trimmedLine == "}" {
			inRepoBlock = false
			continue
		}

		// 在repositories块内部
		if inRepoBlock {
			// 检查预定义仓库
			if match := mavenNameRegex.FindStringSubmatch(trimmedLine); len(match) > 1 {
				repos = append(repos, &model.Repository{
					Name: match[1],
					Type: "maven",
				})
				continue
			}

			// 检查Maven URL
			if match := mavenUrlRegex.FindStringSubmatch(trimmedLine); len(match) > 1 {
				url := match[1]

				// 从URL推断名称
				name := "custom-maven"
				parts := strings.Split(url, "/")
				if len(parts) > 2 {
					name = parts[2]
				}

				repos = append(repos, &model.Repository{
					Name: name,
					URL:  url,
					Type: "maven",
				})
			}
		}
	}

	return repos
}

// GetDefaultRepositories 获取常见的默认仓库
func (rp *RepositoryParser) GetDefaultRepositories() []*model.Repository {
	return []*model.Repository{
		{
			Name: "mavenCentral",
			Type: "maven",
			URL:  "https://repo.maven.apache.org/maven2/",
		},
		{
			Name: "mavenLocal",
			Type: "maven",
		},
		{
			Name: "google",
			Type: "maven",
			URL:  "https://dl.google.com/android/maven2/",
		},
		{
			Name: "jcenter",
			Type: "maven",
			URL:  "https://jcenter.bintray.com/",
		},
	}
}

// HasJitPackRepository 检查是否使用了JitPack仓库
func (rp *RepositoryParser) HasJitPackRepository(repos []*model.Repository) bool {
	for _, repo := range repos {
		if repo.URL != "" && strings.Contains(repo.URL, "jitpack.io") {
			return true
		}
	}
	return false
}

// HasCustomRepository 检查是否使用了自定义仓库
func (rp *RepositoryParser) HasCustomRepository(repos []*model.Repository) bool {
	for _, repo := range repos {
		if repo.Name != "mavenCentral" && repo.Name != "mavenLocal" &&
			repo.Name != "google" && repo.Name != "jcenter" {
			return true
		}
	}
	return false
}
