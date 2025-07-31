# 工具函数

本节记录 Gradle Parser 提供的工具函数和辅助方法。

## 包导入

```go
import "github.com/scagogogo/gradle-parser/pkg/api"
```

## 依赖工具

### DependenciesByScope

按作用域对依赖进行分组以便于分析。

```go
func DependenciesByScope(dependencies []*model.Dependency) []*model.DependencySet
```

**参数:**
- `dependencies` ([]*model.Dependency): 要分组的依赖列表

**返回值:**
- `[]*model.DependencySet`: 按作用域分组的依赖

**示例:**
```go
deps, err := api.GetDependencies("build.gradle")
if err != nil {
    log.Fatal(err)
}

depSets := api.DependenciesByScope(deps)
for _, depSet := range depSets {
    fmt.Printf("\n%s 依赖:\n", depSet.Scope)
    for _, dep := range depSet.Dependencies {
        fmt.Printf("  %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
    }
}
```

## 项目类型检测

### IsAndroidProject

基于插件检查项目是否为 Android 项目。

```go
func IsAndroidProject(plugins []*model.Plugin) bool
```

**检测逻辑:**
- 查找 `com.android.application` 插件
- 查找 `com.android.library` 插件
- 查找 `com.android.test` 插件

### IsKotlinProject

基于插件检查项目是否使用 Kotlin。

```go
func IsKotlinProject(plugins []*model.Plugin) bool
```

**检测逻辑:**
- 查找 `kotlin` 插件
- 查找 `org.jetbrains.kotlin.jvm` 插件
- 查找 `org.jetbrains.kotlin.android` 插件

### IsSpringBootProject

基于插件检查项目是否为 Spring Boot 项目。

```go
func IsSpringBootProject(plugins []*model.Plugin) bool
```

**检测逻辑:**
- 查找 `org.springframework.boot` 插件
- 查找 `io.spring.dependency-management` 插件

**示例:**
```go
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

projectTypes := []string{}

if api.IsAndroidProject(plugins) {
    projectTypes = append(projectTypes, "Android")
}

if api.IsKotlinProject(plugins) {
    projectTypes = append(projectTypes, "Kotlin")
}

if api.IsSpringBootProject(plugins) {
    projectTypes = append(projectTypes, "Spring Boot")
}

if len(projectTypes) > 0 {
    fmt.Printf("项目类型: %s\n", strings.Join(projectTypes, ", "))
} else {
    fmt.Println("标准 Java 项目")
}
```

## 配置工具

### DefaultOptions

返回默认解析器选项。

```go
func DefaultOptions() *Options
```

**返回值:**
- `*Options`: 默认配置选项

**默认值:**
```go
&Options{
    SkipComments:      true,
    CollectRawContent: true,
    ParsePlugins:      true,
    ParseDependencies: true,
    ParseRepositories: true,
    ParseTasks:        true,
}
```

## 版本信息

### Version

Gradle Parser 库的当前版本。

```go
const Version = "0.1.0"
```

**示例:**
```go
fmt.Printf("使用 Gradle Parser v%s\n", api.Version)
```

## 最佳实践

1. **使用适当的检测函数**: 为您的需求选择正确的项目类型检测函数
2. **按作用域分组依赖**: 使用 `DependenciesByScope` 获得更好的组织
3. **验证项目结构**: 检查必需的依赖和插件
4. **处理边缘情况**: 始终检查 nil 值和空集合
5. **组合工具**: 将多个工具函数一起使用进行全面分析
