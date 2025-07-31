# 数据模型

本节记录用于表示解析的 Gradle 项目的数据结构。这些模型提供了访问项目信息、依赖、插件和其他构建配置的结构化方式。

## 核心模型

### Project

表示完整的 Gradle 项目及其所有组件。

```go
type Project struct {
    // 基本项目信息
    Group       string `json:"group"`
    Name        string `json:"name"`
    Version     string `json:"version"`
    Description string `json:"description"`

    // Java/JVM 配置
    SourceCompatibility string `json:"sourceCompatibility"`
    TargetCompatibility string `json:"targetCompatibility"`
    
    // 自定义属性
    Properties map[string]string `json:"properties"`

    // 项目组件
    Plugins      []*Plugin      `json:"plugins"`
    Dependencies []*Dependency  `json:"dependencies"`
    Repositories []*Repository  `json:"repositories"`
    SubProjects  []*Project     `json:"subProjects"`
    Tasks        []*Task        `json:"tasks"`
    Extensions   map[string]any `json:"extensions"`

    // 文件信息
    FilePath string `json:"filePath"`
}
```

**字段:**
- `Group`: Maven 组 ID（例如，"com.example"）
- `Name`: 项目名称
- `Version`: 项目版本
- `Description`: 项目描述
- `Dependencies`: 项目依赖列表
- `Plugins`: 应用的插件列表
- `Repositories`: 配置的仓库列表

### Dependency

表示项目依赖。

```go
type Dependency struct {
    Group      string `json:"group"`
    Name       string `json:"name"`
    Version    string `json:"version"`
    Scope      string `json:"scope"`
    Transitive bool   `json:"transitive"`
    Raw        string `json:"raw"`
}
```

**字段:**
- `Group`: Maven 组 ID（例如，"org.springframework"）
- `Name`: 构件名称（例如，"spring-core"）
- `Version`: 版本字符串（例如，"5.3.21"）
- `Scope`: 依赖作用域（例如，"implementation"、"testImplementation"）
- `Raw`: 构建文件中的原始依赖声明

### Plugin

表示 Gradle 插件配置。

```go
type Plugin struct {
    ID      string                 `json:"id"`
    Version string                 `json:"version,omitempty"`
    Apply   bool                   `json:"apply"`
    Config  map[string]interface{} `json:"config,omitempty"`
}
```

**字段:**
- `ID`: 插件标识符（例如，"java"、"org.springframework.boot"）
- `Version`: 插件版本（可选）
- `Apply`: 插件是否被应用（默认为 true）
- `Config`: 插件特定配置

## 结果模型

### ParseResult

包含解析 Gradle 文件的完整结果。

```go
type ParseResult struct {
    Project   *Project `json:"project"`
    RawText   string   `json:"rawText,omitempty"`
    Errors    []error  `json:"errors,omitempty"`
    Warnings  []string `json:"warnings,omitempty"`
    ParseTime string   `json:"parseTime,omitempty"`
}
```

**字段:**
- `Project`: 解析的项目信息
- `RawText`: 原始文件内容（如果启用收集）
- `Errors`: 致命解析错误
- `Warnings`: 非致命解析警告
- `ParseTime`: 解析文件所用的时间

## 使用示例

### 访问项目信息

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project
fmt.Printf("项目: %s\n", project.Name)
fmt.Printf("组: %s\n", project.Group)
fmt.Printf("版本: %s\n", project.Version)
```

### 处理依赖

```go
for _, dep := range project.Dependencies {
    fmt.Printf("依赖: %s:%s", dep.Group, dep.Name)
    if dep.Version != "" {
        fmt.Printf(":%s", dep.Version)
    }
    fmt.Printf(" (%s)\n", dep.Scope)
}
```

## JSON 序列化

所有模型都支持 JSON 序列化，便于与 Web API 和数据存储集成：

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

// 序列化为 JSON
jsonData, err := json.MarshalIndent(result.Project, "", "  ")
if err != nil {
    log.Fatal(err)
}

fmt.Println(string(jsonData))
```

## 类型安全

所有模型都使用 Go 的类型系统来确保数据完整性：

- 文本数据使用字符串字段
- 集合使用切片
- 键值配置使用映射
- 可选关系使用指针
- 可扩展配置使用接口
