// Package parser 提供用于解析Gradle文件的核心功能
package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/scagogogo/gradle-parser/pkg/config"
	"github.com/scagogogo/gradle-parser/pkg/dependency"
	"github.com/scagogogo/gradle-parser/pkg/model"
)

// Parser 定义Gradle解析器接口
type Parser interface {
	// Parse 解析Gradle字符串内容
	Parse(content string) (*model.ParseResult, error)

	// ParseFile 解析Gradle文件
	ParseFile(filePath string) (*model.ParseResult, error)

	// ParseReader 从Reader中解析Gradle内容
	ParseReader(reader io.Reader) (*model.ParseResult, error)
}

// GradleParser 是默认的Gradle解析器实现
type GradleParser struct {
	// 解析配置选项
	skipComments      bool
	collectRawContent bool
	parsePlugins      bool
	parseDependencies bool
	parseRepositories bool
	parseTasks        bool

	// 当前解析状态
	currentBlock *model.ScriptBlock
	errors       []error
	warnings     []string
}

// NewParser 创建新的默认解析器实例
func NewParser() Parser {
	return &GradleParser{
		skipComments:      true,
		collectRawContent: true,
		parsePlugins:      true,
		parseDependencies: true,
		parseRepositories: true,
		parseTasks:        true,
		errors:            make([]error, 0),
		warnings:          make([]string, 0),
	}
}

// ParseFile 从文件解析Gradle配置
func (p *GradleParser) ParseFile(filePath string) (*model.ParseResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开Gradle文件: %w", err)
	}
	defer file.Close()

	result, err := p.ParseReader(file)
	if err != nil {
		return nil, err
	}

	// 设置文件路径
	if result.Project != nil {
		result.Project.FilePath = filePath
		// 如果项目名称为空，尝试从文件名推断
		if result.Project.Name == "" {
			dir := filepath.Dir(filePath)
			result.Project.Name = filepath.Base(dir)
		}
	}

	return result, nil
}

// ParseReader 从Reader中解析Gradle配置
func (p *GradleParser) ParseReader(reader io.Reader) (*model.ParseResult, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取Gradle内容失败: %w", err)
	}

	return p.Parse(string(content))
}

// Parse 从字符串解析Gradle配置
func (p *GradleParser) Parse(content string) (*model.ParseResult, error) {
	// 重置解析状态
	p.currentBlock = &model.ScriptBlock{
		Name:     "root",
		Children: make([]*model.ScriptBlock, 0),
		Values:   make(map[string]interface{}),
		Closures: make(map[string][]*model.ScriptBlock),
	}
	p.errors = make([]error, 0)
	p.warnings = make([]string, 0)

	// 记录开始时间
	startTime := time.Now()

	// 创建项目对象
	project := &model.Project{
		Properties:   make(map[string]string),
		Plugins:      make([]*model.Plugin, 0),
		Dependencies: make([]*model.Dependency, 0),
		Repositories: make([]*model.Repository, 0),
		SubProjects:  make([]*model.Project, 0),
		Tasks:        make([]*model.Task, 0),
		Extensions:   make(map[string]any),
	}

	// 使用scanner逐行解析
	scanner := bufio.NewScanner(strings.NewReader(content))
	var rawLines []string
	if p.collectRawContent {
		rawLines = make([]string, 0, strings.Count(content, "\n")+1)
	}

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// 收集原始内容
		if p.collectRawContent {
			rawLines = append(rawLines, line)
		}

		// 处理空行和注释
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || (p.skipComments && (strings.HasPrefix(trimmedLine, "//") || strings.HasPrefix(trimmedLine, "/*"))) {
			continue
		}

		// 解析行内容
		if err := p.parseLine(trimmedLine, lineNumber, project); err != nil {
			// 不把解析错误当作致命错误，只记录警告
			p.warnings = append(p.warnings, fmt.Sprintf("行 %d: %v", lineNumber, err))
		}
	}

	// 使用专门的解析器来提取依赖、插件和仓库
	if p.parseDependencies {
		depParser := dependency.NewDependencyParser()
		project.Dependencies = depParser.ExtractDependenciesFromText(content)
	}

	if p.parsePlugins {
		pluginParser := config.NewPluginParser()
		project.Plugins = pluginParser.ExtractPluginsFromText(content)
	}

	if p.parseRepositories {
		repoParser := config.NewRepositoryParser()
		project.Repositories = repoParser.ExtractRepositoriesFromText(content)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("扫描内容时出错: %w", err)
	}

	// 完成解析
	result := &model.ParseResult{
		Project:   project,
		Errors:    p.errors,
		Warnings:  p.warnings,
		ParseTime: time.Since(startTime).String(),
	}

	if p.collectRawContent {
		result.RawText = strings.Join(rawLines, "\n")
	}

	return result, nil
}

// parseLine 解析单行内容
func (p *GradleParser) parseLine(line string, lineNumber int, project *model.Project) error {
	line = strings.TrimSpace(line)

	// 跳过空行和注释
	if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
		return nil
	}

	// 解析项目基本属性
	if err := p.parseProjectProperty(line, project); err == nil {
		return nil
	}

	// 解析插件块
	if strings.HasPrefix(line, "plugins") {
		return p.parsePluginsBlock(line, project)
	}

	// 解析依赖块
	if strings.HasPrefix(line, "dependencies") {
		return p.parseDependenciesBlock(line, project)
	}

	// 解析仓库块
	if strings.HasPrefix(line, "repositories") {
		return p.parseRepositoriesBlock(line, project)
	}

	// 解析任务定义
	if strings.HasPrefix(line, "task ") || strings.Contains(line, "task(") {
		return p.parseTaskDefinition(line, project)
	}

	// 其他配置项暂时忽略，不报错
	return nil
}

// WithSkipComments 设置是否跳过注释
func (p *GradleParser) WithSkipComments(skip bool) *GradleParser {
	p.skipComments = skip
	return p
}

// WithCollectRawContent 设置是否收集原始内容
func (p *GradleParser) WithCollectRawContent(collect bool) *GradleParser {
	p.collectRawContent = collect
	return p
}

// WithParsePlugins 设置是否解析插件
func (p *GradleParser) WithParsePlugins(parse bool) *GradleParser {
	p.parsePlugins = parse
	return p
}

// WithParseDependencies 设置是否解析依赖
func (p *GradleParser) WithParseDependencies(parse bool) *GradleParser {
	p.parseDependencies = parse
	return p
}

// WithParseRepositories 设置是否解析仓库
func (p *GradleParser) WithParseRepositories(parse bool) *GradleParser {
	p.parseRepositories = parse
	return p
}

// WithParseTasks 设置是否解析任务
func (p *GradleParser) WithParseTasks(parse bool) *GradleParser {
	p.parseTasks = parse
	return p
}

// parseProjectProperty 解析项目基本属性
func (p *GradleParser) parseProjectProperty(line string, project *model.Project) error {
	// 匹配 key = value 格式
	if strings.Contains(line, "=") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid assignment format")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 移除引号
		value = strings.Trim(value, `"'`)

		switch key {
		case "group":
			project.Group = value
		case "version":
			project.Version = value
		case "description":
			project.Description = value
		case "sourceCompatibility":
			project.SourceCompatibility = value
		case "targetCompatibility":
			project.TargetCompatibility = value
		default:
			// 其他属性存储在Properties中
			if project.Properties == nil {
				project.Properties = make(map[string]string)
			}
			project.Properties[key] = value
		}
		return nil
	}

	return fmt.Errorf("not a property assignment")
}

// parsePluginsBlock 解析插件块
func (p *GradleParser) parsePluginsBlock(line string, project *model.Project) error {
	if !p.parsePlugins {
		return nil
	}

	// 简单的插件解析 - 这里可以扩展为更复杂的块解析
	// 目前只处理单行插件声明
	if strings.Contains(line, "id") {
		// 匹配 id 'plugin-name' version 'version'
		// 或 id("plugin-name") version "version"
		plugin := &model.Plugin{Apply: true}

		// 提取插件ID
		if idMatch := extractQuotedValue(line, "id"); idMatch != "" {
			plugin.ID = idMatch
		}

		// 提取版本
		if versionMatch := extractQuotedValue(line, "version"); versionMatch != "" {
			plugin.Version = versionMatch
		}

		if plugin.ID != "" {
			project.Plugins = append(project.Plugins, plugin)
		}
	}

	return nil
}

// parseDependenciesBlock 解析依赖块
func (p *GradleParser) parseDependenciesBlock(line string, project *model.Project) error {
	if !p.parseDependencies {
		return nil
	}

	// 依赖解析已经在dependency包中实现，这里暂时跳过
	// 实际应用中可以在这里调用dependency.ExtractDependenciesFromText
	return nil
}

// parseRepositoriesBlock 解析仓库块
func (p *GradleParser) parseRepositoriesBlock(line string, project *model.Project) error {
	if !p.parseRepositories {
		return nil
	}

	// 简单的仓库解析
	if strings.Contains(line, "mavenCentral") {
		repo := &model.Repository{
			Name: "mavenCentral",
			Type: "maven",
			URL:  "https://repo1.maven.org/maven2/",
		}
		project.Repositories = append(project.Repositories, repo)
	} else if strings.Contains(line, "google") {
		repo := &model.Repository{
			Name: "google",
			Type: "maven",
			URL:  "https://dl.google.com/dl/android/maven2/",
		}
		project.Repositories = append(project.Repositories, repo)
	} else if strings.Contains(line, "maven") && strings.Contains(line, "url") {
		// 解析自定义maven仓库
		repo := &model.Repository{
			Name: "custom",
			Type: "maven",
		}
		if url := extractQuotedValue(line, "url"); url != "" {
			repo.URL = url
		}
		project.Repositories = append(project.Repositories, repo)
	}

	return nil
}

// parseTaskDefinition 解析任务定义
func (p *GradleParser) parseTaskDefinition(line string, project *model.Project) error {
	if !p.parseTasks {
		return nil
	}

	// 简单的任务解析
	task := &model.Task{}

	// 提取任务名称
	if strings.HasPrefix(line, "task ") {
		parts := strings.Fields(line)
		if len(parts) > 1 {
			task.Name = parts[1]
		}
	}

	if task.Name != "" {
		project.Tasks = append(project.Tasks, task)
	}

	return nil
}

// extractQuotedValue 从行中提取引号包围的值
func extractQuotedValue(line, keyword string) string {
	// 查找关键字位置
	keywordIndex := strings.Index(line, keyword)
	if keywordIndex == -1 {
		return ""
	}

	// 从关键字后开始查找引号
	searchStart := keywordIndex + len(keyword)
	remaining := line[searchStart:]

	// 查找第一个引号
	var quote rune
	var start int = -1
	for i, r := range remaining {
		if r == '"' || r == '\'' {
			quote = r
			start = i + 1
			break
		}
	}

	if start == -1 {
		return ""
	}

	// 查找匹配的结束引号
	for i := start; i < len(remaining); i++ {
		if rune(remaining[i]) == quote {
			return remaining[start:i]
		}
	}

	return ""
}
