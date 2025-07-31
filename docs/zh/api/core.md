# 核心 API

核心 API 提供了解析 Gradle 构建文件的主要入口点。这些函数设计简单易用，涵盖了最常见的解析场景。

## 包导入

```go
import "github.com/scagogogo/gradle-parser/pkg/api"
```

## 基础解析函数

### ParseFile

从文件系统解析 Gradle 文件。

```go
func ParseFile(filePath string) (*model.ParseResult, error)
```

**参数:**
- `filePath` (string): 要解析的 Gradle 文件路径

**返回值:**
- `*model.ParseResult`: 解析的项目信息和元数据
- `error`: 解析失败时的错误

**示例:**
```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project
fmt.Printf("项目: %s v%s\n", project.Name, project.Version)
```

### ParseString

从字符串解析 Gradle 内容。

```go
func ParseString(content string) (*model.ParseResult, error)
```

**参数:**
- `content` (string): 作为字符串的 Gradle 文件内容

**返回值:**
- `*model.ParseResult`: 解析的项目信息和元数据
- `error`: 解析失败时的错误

## 组件提取函数

### GetDependencies

从 Gradle 文件提取依赖。

```go
func GetDependencies(filePath string) ([]*model.Dependency, error)
```

**参数:**
- `filePath` (string): Gradle 文件路径

**返回值:**
- `[]*model.Dependency`: 文件中找到的依赖列表
- `error`: 提取失败时的错误

**示例:**
```go
deps, err := api.GetDependencies("build.gradle")
if err != nil {
    log.Fatal(err)
}

for _, dep := range deps {
    fmt.Printf("%s:%s:%s (%s)\n", dep.Group, dep.Name, dep.Version, dep.Scope)
}
```

## 项目类型检测

### IsAndroidProject

基于插件检查项目是否为 Android 项目。

```go
func IsAndroidProject(plugins []*model.Plugin) bool
```

### IsKotlinProject

基于插件检查项目是否使用 Kotlin。

```go
func IsKotlinProject(plugins []*model.Plugin) bool
```

### IsSpringBootProject

基于插件检查项目是否为 Spring Boot 项目。

```go
func IsSpringBootProject(plugins []*model.Plugin) bool
```

**示例:**
```go
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

if api.IsAndroidProject(plugins) {
    fmt.Println("这是一个 Android 项目")
}
if api.IsKotlinProject(plugins) {
    fmt.Println("这个项目使用 Kotlin")
}
if api.IsSpringBootProject(plugins) {
    fmt.Println("这是一个 Spring Boot 项目")
}
```

## 错误处理

所有解析函数都返回详细的错误信息。常见错误类型包括：

- **文件未找到**: 当指定的文件不存在时
- **权限被拒绝**: 当文件由于权限无法读取时
- **解析错误**: 当 Gradle 语法无效时

**最佳实践:**
1. 在使用结果之前始终检查错误
2. 适当地记录或处理错误
3. 检查 `ParseResult` 中的 `Warnings` 字段以了解非致命问题

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Printf("解析文件失败: %v", err)
    return
}

if len(result.Warnings) > 0 {
    log.Printf("解析完成，有 %d 个警告", len(result.Warnings))
    for _, warning := range result.Warnings {
        log.Printf("警告: %s", warning)
    }
}
```
