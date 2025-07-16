// Package parser 提供位置感知的Gradle解析器
package parser

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

// SourceAwareParser 位置感知的Gradle解析器
type SourceAwareParser struct {
	*GradleParser

	// 位置追踪
	currentLine   int
	currentColumn int
	currentPos    int

	// 原始文本信息
	originalText string
	lines        []string
}

// NewSourceAwareParser 创建新的位置感知解析器
func NewSourceAwareParser() *SourceAwareParser {
	return &SourceAwareParser{
		GradleParser: NewParser().(*GradleParser),
	}
}

// ParseWithSourceMapping 解析并返回带源码位置信息的结果
func (sap *SourceAwareParser) ParseWithSourceMapping(content string) (*model.SourceMappedParseResult, error) {
	// 初始化位置追踪
	sap.originalText = content
	sap.lines = strings.Split(content, "\n")
	sap.currentLine = 1
	sap.currentColumn = 1
	sap.currentPos = 0

	// 先进行常规解析
	result, err := sap.Parse(content)
	if err != nil {
		return nil, err
	}

	// 创建带源码位置信息的项目
	sourceMappedProject := &model.SourceMappedProject{
		Project:                  result.Project,
		OriginalText:             content,
		Lines:                    sap.lines,
		SourceMappedDependencies: make([]*model.SourceMappedDependency, 0),
		SourceMappedPlugins:      make([]*model.SourceMappedPlugin, 0),
		SourceMappedRepositories: make([]*model.SourceMappedRepository, 0),
		SourceMappedProperties:   make([]*model.SourceMappedProperty, 0),
	}

	// 解析带位置信息的组件
	if err := sap.parseSourceMappedComponents(content, sourceMappedProject); err != nil {
		return nil, err
	}

	return &model.SourceMappedParseResult{
		ParseResult:         result,
		SourceMappedProject: sourceMappedProject,
	}, nil
}

// parseSourceMappedComponents 解析带位置信息的组件
func (sap *SourceAwareParser) parseSourceMappedComponents(content string, project *model.SourceMappedProject) error {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNumber := 0
	currentPos := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		lineStart := currentPos
		lineEnd := currentPos + len(line)

		// 解析属性
		if err := sap.parseSourceMappedProperty(line, lineNumber, lineStart, project); err == nil {
			// 属性解析成功，继续下一行
		} else if err := sap.parseSourceMappedDependency(line, lineNumber, lineStart, project); err == nil {
			// 依赖解析成功
		} else if err := sap.parseSourceMappedPlugin(line, lineNumber, lineStart, project); err == nil {
			// 插件解析成功
		} else if err := sap.parseSourceMappedRepository(line, lineNumber, lineStart, project); err == nil {
			// 仓库解析成功
		}

		// 更新位置（+1 for newline character）
		currentPos = lineEnd + 1
	}

	return scanner.Err()
}

// parseSourceMappedProperty 解析带位置信息的属性
func (sap *SourceAwareParser) parseSourceMappedProperty(line string, lineNumber, lineStart int, project *model.SourceMappedProject) error {
	trimmedLine := strings.TrimSpace(line)

	// 匹配 key = value 格式
	if strings.Contains(trimmedLine, "=") {
		parts := strings.SplitN(trimmedLine, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid assignment format")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)

		// 计算在行内的位置
		keyStart := strings.Index(line, key)
		if keyStart == -1 {
			return fmt.Errorf("key not found in line")
		}

		valueStart := strings.Index(line, parts[1])
		if valueStart == -1 {
			return fmt.Errorf("value not found in line")
		}

		// 创建源码位置信息
		sourceRange := model.SourceRange{
			Start: model.SourcePosition{
				Line:     lineNumber,
				Column:   keyStart + 1,
				StartPos: lineStart + keyStart,
				EndPos:   lineStart + len(line),
				Length:   len(line) - keyStart,
			},
			End: model.SourcePosition{
				Line:     lineNumber,
				Column:   len(line),
				StartPos: lineStart + len(line),
				EndPos:   lineStart + len(line),
				Length:   0,
			},
		}

		sourceMappedProperty := &model.SourceMappedProperty{
			Key:         key,
			Value:       value,
			SourceRange: sourceRange,
			RawText:     line,
		}

		project.SourceMappedProperties = append(project.SourceMappedProperties, sourceMappedProperty)
		return nil
	}

	return fmt.Errorf("not a property assignment")
}

// parseSourceMappedDependency 解析带位置信息的依赖
func (sap *SourceAwareParser) parseSourceMappedDependency(line string, lineNumber, lineStart int, project *model.SourceMappedProject) error {
	trimmedLine := strings.TrimSpace(line)

	// 使用依赖解析器的正则表达式
	patterns := []string{
		`['"]([^'"]+):([^'"]+):([^'"]+)['"]`,           // "group:name:version"
		`['"]([^'"]+):([^'"]+)['"]`,                    // "group:name" (没有版本号)
		`['"]([^'"]+)\.([^'"]+):([^'"]+):([^'"]+)['"]`, // "group.name:name:version"
		`project\(['"]:(.*)['"]\)`,                     // project(":name")
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(trimmedLine, -1)

		for _, match := range matches {
			if len(match) > 0 {
				rawDep := match[0]

				// 查找依赖在行中的位置
				depStart := strings.Index(line, rawDep)
				if depStart == -1 {
					continue
				}

				// 解析依赖 - 使用简单的解析逻辑
				dep := &model.Dependency{
					Raw: rawDep,
				}

				// 简单解析group:name:version格式
				if strings.Contains(rawDep, ":") {
					parts := strings.Split(strings.Trim(rawDep, `"'`), ":")
					if len(parts) >= 2 {
						dep.Group = parts[0]
						dep.Name = parts[1]
						if len(parts) >= 3 {
							dep.Version = parts[2]
						}
					}
				}

				// 创建源码位置信息
				sourceRange := model.SourceRange{
					Start: model.SourcePosition{
						Line:     lineNumber,
						Column:   depStart + 1,
						StartPos: lineStart + depStart,
						EndPos:   lineStart + depStart + len(rawDep),
						Length:   len(rawDep),
					},
					End: model.SourcePosition{
						Line:     lineNumber,
						Column:   depStart + len(rawDep),
						StartPos: lineStart + depStart + len(rawDep),
						EndPos:   lineStart + depStart + len(rawDep),
						Length:   0,
					},
				}

				sourceMappedDep := &model.SourceMappedDependency{
					Dependency:  dep,
					SourceRange: sourceRange,
					RawText:     rawDep,
				}

				project.SourceMappedDependencies = append(project.SourceMappedDependencies, sourceMappedDep)
				return nil
			}
		}
	}

	return fmt.Errorf("not a dependency")
}

// parseSourceMappedPlugin 解析带位置信息的插件
func (sap *SourceAwareParser) parseSourceMappedPlugin(line string, lineNumber, lineStart int, project *model.SourceMappedProject) error {
	trimmedLine := strings.TrimSpace(line)

	// 使用插件解析器的正则表达式
	pluginRegex := regexp.MustCompile(`id\s*\(?['"](.*?)['"](\))?(\s+version\s*['"](.*?)['"])?`)

	if matches := pluginRegex.FindStringSubmatch(trimmedLine); len(matches) > 1 {
		// 查找插件声明在行中的位置
		pluginStart := strings.Index(line, matches[0])
		if pluginStart == -1 {
			return fmt.Errorf("plugin declaration not found in line")
		}

		plugin := &model.Plugin{
			ID:    matches[1],
			Apply: true,
		}

		// 检查是否有版本信息
		if len(matches) > 4 && matches[4] != "" {
			plugin.Version = matches[4]
		}

		// 创建源码位置信息
		sourceRange := model.SourceRange{
			Start: model.SourcePosition{
				Line:     lineNumber,
				Column:   pluginStart + 1,
				StartPos: lineStart + pluginStart,
				EndPos:   lineStart + pluginStart + len(matches[0]),
				Length:   len(matches[0]),
			},
			End: model.SourcePosition{
				Line:     lineNumber,
				Column:   pluginStart + len(matches[0]),
				StartPos: lineStart + pluginStart + len(matches[0]),
				EndPos:   lineStart + pluginStart + len(matches[0]),
				Length:   0,
			},
		}

		sourceMappedPlugin := &model.SourceMappedPlugin{
			Plugin:      plugin,
			SourceRange: sourceRange,
			RawText:     matches[0],
		}

		project.SourceMappedPlugins = append(project.SourceMappedPlugins, sourceMappedPlugin)
		return nil
	}

	return fmt.Errorf("not a plugin")
}

// parseSourceMappedRepository 解析带位置信息的仓库
func (sap *SourceAwareParser) parseSourceMappedRepository(line string, lineNumber, lineStart int, project *model.SourceMappedProject) error {
	trimmedLine := strings.TrimSpace(line)

	// 检查常见的仓库声明
	repoPatterns := map[string]string{
		"mavenCentral()": "mavenCentral",
		"google()":       "google",
		"jcenter()":      "jcenter",
		"mavenLocal()":   "mavenLocal",
	}

	for pattern, name := range repoPatterns {
		if strings.Contains(trimmedLine, pattern) {
			repoStart := strings.Index(line, pattern)
			if repoStart == -1 {
				continue
			}

			repo := &model.Repository{
				Name: name,
				Type: "maven",
			}

			// 创建源码位置信息
			sourceRange := model.SourceRange{
				Start: model.SourcePosition{
					Line:     lineNumber,
					Column:   repoStart + 1,
					StartPos: lineStart + repoStart,
					EndPos:   lineStart + repoStart + len(pattern),
					Length:   len(pattern),
				},
				End: model.SourcePosition{
					Line:     lineNumber,
					Column:   repoStart + len(pattern),
					StartPos: lineStart + repoStart + len(pattern),
					EndPos:   lineStart + repoStart + len(pattern),
					Length:   0,
				},
			}

			sourceMappedRepo := &model.SourceMappedRepository{
				Repository:  repo,
				SourceRange: sourceRange,
				RawText:     pattern,
			}

			project.SourceMappedRepositories = append(project.SourceMappedRepositories, sourceMappedRepo)
			return nil
		}
	}

	return fmt.Errorf("not a repository")
}
