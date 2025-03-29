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
			p.errors = append(p.errors, fmt.Errorf("行 %d: %w", lineNumber, err))
		}
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
	// TODO: 实现行解析逻辑
	// 这里将根据不同的语法规则来解析Gradle配置

	// 这是一个占位函数，实际实现会更复杂
	// 需要处理各种语法结构，如块、赋值、方法调用等

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
