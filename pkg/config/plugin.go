// Package config 提供Gradle配置解析功能
package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

var (
	// 匹配插件ID的正则表达式
	// 例如: id 'com.android.application' version '7.0.0'
	// 或者: id("org.jetbrains.kotlin.android") version "1.5.30"
	pluginRegex = regexp.MustCompile(`id\s*\(?['"](.*?)['"](\))?(\s+version\s*['"](.*?)['"])?`)

	// 匹配apply plugin的正则表达式
	// 例如: apply plugin: 'java'
	applyPluginRegex = regexp.MustCompile(`apply\s+plugin:\s*['"](.*?)['"]`)
)

// PluginParser 处理Gradle插件解析
type PluginParser struct{}

// NewPluginParser 创建新的插件解析器
func NewPluginParser() *PluginParser {
	return &PluginParser{}
}

// ParsePluginBlock 解析插件块
func (pp *PluginParser) ParsePluginBlock(block *model.ScriptBlock) ([]*model.Plugin, error) {
	if block == nil {
		return nil, fmt.Errorf("插件块为空")
	}

	plugins := make([]*model.Plugin, 0)

	// 处理plugins {} 块中的插件声明
	for _, value := range block.Values {
		valueStr := fmt.Sprintf("%v", value)
		if matches := pluginRegex.FindStringSubmatch(valueStr); len(matches) > 1 {
			plugin := &model.Plugin{
				ID:    matches[1],
				Apply: true,
			}

			// 检查是否有版本信息
			if len(matches) > 4 && matches[4] != "" {
				plugin.Version = matches[4]
			}

			plugins = append(plugins, plugin)
		}
	}

	// 处理插件块中的子闭包
	for name, closures := range block.Closures {
		// id 闭包通常用于声明插件
		if name == "id" {
			for _, closure := range closures {
				for _, value := range closure.Values {
					valueStr := fmt.Sprintf("%v", value)
					plugin := &model.Plugin{
						ID:    valueStr,
						Apply: true,
					}
					plugins = append(plugins, plugin)
				}
			}
		}
	}

	return plugins, nil
}

// ExtractPluginsFromText 从原始文本中提取插件
func (pp *PluginParser) ExtractPluginsFromText(text string) []*model.Plugin {
	plugins := make([]*model.Plugin, 0)

	// 分析文本中的插件声明
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检查plugins块中的插件声明
		if matches := pluginRegex.FindStringSubmatch(trimmedLine); len(matches) > 1 {
			plugin := &model.Plugin{
				ID:    matches[1],
				Apply: true,
			}

			// 检查是否有版本信息
			if len(matches) > 4 && matches[4] != "" {
				plugin.Version = matches[4]
			}

			plugins = append(plugins, plugin)
		}

		// 检查apply plugin语句
		if matches := applyPluginRegex.FindStringSubmatch(trimmedLine); len(matches) > 1 {
			plugin := &model.Plugin{
				ID:    matches[1],
				Apply: true,
			}
			plugins = append(plugins, plugin)
		}
	}

	return plugins
}

// GetPluginConfigurations 获取插件相关的配置块
func (pp *PluginParser) GetPluginConfigurations(rootBlock *model.ScriptBlock, plugins []*model.Plugin) map[string]*model.ScriptBlock {
	// 创建插件ID到配置块的映射
	pluginConfigs := make(map[string]*model.ScriptBlock)

	// 已知的插件配置块名称
	knownConfigs := map[string][]string{
		"com.android.application":      {"android"},
		"com.android.library":          {"android"},
		"java":                         {"java", "sourceCompatibility", "targetCompatibility"},
		"kotlin":                       {"kotlin", "kotlinOptions"},
		"org.jetbrains.kotlin.android": {"kotlin", "kotlinOptions"},
		"org.springframework.boot":     {"springBoot"},
	}

	// 为每个插件查找可能的配置块
	for _, plugin := range plugins {
		// 检查是否有已知的配置块名称
		if configNames, ok := knownConfigs[plugin.ID]; ok {
			for _, configName := range configNames {
				if blocks, ok := rootBlock.Closures[configName]; ok && len(blocks) > 0 {
					// 使用插件ID作为键，存储配置块
					pluginConfigs[plugin.ID] = blocks[0]
				}
			}
		}
	}

	return pluginConfigs
}

// IsAndroidProject 判断是否是Android项目
func (pp *PluginParser) IsAndroidProject(plugins []*model.Plugin) bool {
	for _, plugin := range plugins {
		if plugin.ID == "com.android.application" || plugin.ID == "com.android.library" {
			return true
		}
	}
	return false
}

// IsSpringBootProject 判断是否是Spring Boot项目
func (pp *PluginParser) IsSpringBootProject(plugins []*model.Plugin) bool {
	for _, plugin := range plugins {
		if plugin.ID == "org.springframework.boot" {
			return true
		}
	}
	return false
}

// IsKotlinProject 判断是否是Kotlin项目
func (pp *PluginParser) IsKotlinProject(plugins []*model.Plugin) bool {
	for _, plugin := range plugins {
		if plugin.ID == "kotlin" || plugin.ID == "org.jetbrains.kotlin.jvm" ||
			plugin.ID == "org.jetbrains.kotlin.android" {
			return true
		}
	}
	return false
}
