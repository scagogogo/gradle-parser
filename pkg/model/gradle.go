// Package model 提供解析Gradle配置文件所需的数据结构
package model

// Project 表示Gradle项目结构
type Project struct {
	// 项目基本信息
	Group       string `json:"group"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	// 项目配置
	SourceCompatibility string            `json:"sourceCompatibility"`
	TargetCompatibility string            `json:"targetCompatibility"`
	Properties          map[string]string `json:"properties"`

	// 核心组件
	Plugins      []*Plugin      `json:"plugins"`
	Dependencies []*Dependency  `json:"dependencies"`
	Repositories []*Repository  `json:"repositories"`
	SubProjects  []*Project     `json:"subProjects"`
	Tasks        []*Task        `json:"tasks"`
	Extensions   map[string]any `json:"extensions"`

	// 原始文件路径
	FilePath string `json:"filePath"`
}

// Dependency 表示Gradle依赖
type Dependency struct {
	Group      string `json:"group"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Scope      string `json:"scope"` // implementation, api, testImplementation, etc.
	Transitive bool   `json:"transitive"`
	Raw        string `json:"raw"` // 原始依赖声明
}

// Plugin 表示Gradle插件
type Plugin struct {
	ID      string                 `json:"id"`
	Version string                 `json:"version,omitempty"`
	Apply   bool                   `json:"apply"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

// Repository 表示Gradle仓库配置
type Repository struct {
	Name     string                 `json:"name"`
	URL      string                 `json:"url,omitempty"`
	Type     string                 `json:"type"` // maven, ivy, flatDir, etc.
	Config   map[string]interface{} `json:"config,omitempty"`
	Username string                 `json:"username,omitempty"`
	Password string                 `json:"password,omitempty"`
}

// Task 表示Gradle任务
type Task struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type,omitempty"`
	Description string                 `json:"description,omitempty"`
	Group       string                 `json:"group,omitempty"`
	DependsOn   []string               `json:"dependsOn,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// ScriptBlock 表示Gradle脚本块
type ScriptBlock struct {
	Name     string                    `json:"name"`
	Parent   *ScriptBlock              `json:"-"`
	Children []*ScriptBlock            `json:"children,omitempty"`
	Values   map[string]interface{}    `json:"values,omitempty"`
	Closures map[string][]*ScriptBlock `json:"closures,omitempty"`
}

// DependencySet 表示一组依赖，用于按范围分组
type DependencySet struct {
	Scope        string        `json:"scope"`
	Dependencies []*Dependency `json:"dependencies"`
}

// ParseResult 表示解析结果
type ParseResult struct {
	Project   *Project `json:"project"`
	RawText   string   `json:"rawText,omitempty"`
	Errors    []error  `json:"errors,omitempty"`
	Warnings  []string `json:"warnings,omitempty"`
	ParseTime string   `json:"parseTime,omitempty"`
}
