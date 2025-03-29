// Package api 提供解析Gradle配置文件的API
package api

import (
	"io"
	"os"

	"github.com/scagogogo/gradle-parser/pkg/config"
	"github.com/scagogogo/gradle-parser/pkg/dependency"
	"github.com/scagogogo/gradle-parser/pkg/model"
	"github.com/scagogogo/gradle-parser/pkg/parser"
)

// 版本信息
const (
	Version = "0.1.0"
)

// ParseFile 解析指定路径的Gradle文件
func ParseFile(filePath string) (*model.ParseResult, error) {
	parser := parser.NewParser()
	return parser.ParseFile(filePath)
}

// ParseString 解析Gradle字符串内容
func ParseString(content string) (*model.ParseResult, error) {
	parser := parser.NewParser()
	return parser.Parse(content)
}

// ParseReader 从Reader解析Gradle内容
func ParseReader(reader io.Reader) (*model.ParseResult, error) {
	parser := parser.NewParser()
	return parser.ParseReader(reader)
}

// GetDependencies 从文件提取依赖信息
func GetDependencies(filePath string) ([]*model.Dependency, error) {
	// 尝试打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取整个文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// 创建依赖解析器
	depParser := dependency.NewDependencyParser()

	// 直接从文本提取依赖
	return depParser.ExtractDependenciesFromText(string(content)), nil
}

// GetPlugins 从文件提取插件信息
func GetPlugins(filePath string) ([]*model.Plugin, error) {
	// 尝试打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取整个文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// 创建插件解析器
	pluginParser := config.NewPluginParser()

	// 直接从文本提取插件
	return pluginParser.ExtractPluginsFromText(string(content)), nil
}

// GetRepositories 从文件提取仓库信息
func GetRepositories(filePath string) ([]*model.Repository, error) {
	// 尝试打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取整个文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// 创建仓库解析器
	repoParser := config.NewRepositoryParser()

	// 直接从文本提取仓库
	return repoParser.ExtractRepositoriesFromText(string(content)), nil
}

// DependenciesByScope 按范围对依赖进行分组
func DependenciesByScope(dependencies []*model.Dependency) []*model.DependencySet {
	depParser := dependency.NewDependencyParser()
	return depParser.GroupDependenciesByScope(dependencies)
}

// IsAndroidProject 检查是否是Android项目
func IsAndroidProject(plugins []*model.Plugin) bool {
	pluginParser := config.NewPluginParser()
	return pluginParser.IsAndroidProject(plugins)
}

// IsKotlinProject 检查是否是Kotlin项目
func IsKotlinProject(plugins []*model.Plugin) bool {
	pluginParser := config.NewPluginParser()
	return pluginParser.IsKotlinProject(plugins)
}

// IsSpringBootProject 检查是否是Spring Boot项目
func IsSpringBootProject(plugins []*model.Plugin) bool {
	pluginParser := config.NewPluginParser()
	return pluginParser.IsSpringBootProject(plugins)
}

// Options 解析选项
type Options struct {
	SkipComments      bool
	CollectRawContent bool
	ParsePlugins      bool
	ParseDependencies bool
	ParseRepositories bool
	ParseTasks        bool
}

// DefaultOptions 创建默认选项
func DefaultOptions() *Options {
	return &Options{
		SkipComments:      true,
		CollectRawContent: true,
		ParsePlugins:      true,
		ParseDependencies: true,
		ParseRepositories: true,
		ParseTasks:        true,
	}
}

// NewParser 创建自定义配置的解析器
func NewParser(options *Options) parser.Parser {
	p := parser.NewParser().(*parser.GradleParser)

	if options != nil {
		p.WithSkipComments(options.SkipComments)
		p.WithCollectRawContent(options.CollectRawContent)
		p.WithParsePlugins(options.ParsePlugins)
		p.WithParseDependencies(options.ParseDependencies)
		p.WithParseRepositories(options.ParseRepositories)
		p.WithParseTasks(options.ParseTasks)
	}

	return p
}
