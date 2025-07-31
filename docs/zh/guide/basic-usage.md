# 基本用法

本指南涵盖了使用 Gradle Parser 可以执行的基本操作，包括解析文件、提取信息和处理解析数据。

## 解析 Gradle 文件

### 简单文件解析

解析 Gradle 文件的最直接方式：

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    result, err := api.ParseFile("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    project := result.Project
    fmt.Printf("项目名称: %s\n", project.Name)
    fmt.Printf("项目组: %s\n", project.Group)
    fmt.Printf("项目版本: %s\n", project.Version)
}
```

### 解析不同来源

您可以从各种来源解析 Gradle 内容：

```go
// 从文件解析
result, err := api.ParseFile("build.gradle")

// 从字符串解析
gradleContent := `
plugins {
    id 'java'
}
group = 'com.example'
version = '1.0.0'
`
result, err = api.ParseString(gradleContent)

// 从任何 io.Reader 解析
file, err := os.Open("build.gradle")
if err != nil {
    log.Fatal(err)
}
defer file.Close()
result, err = api.ParseReader(file)
```

## 处理依赖

### 提取依赖

从 Gradle 文件获取所有依赖：

```go
deps, err := api.GetDependencies("build.gradle")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("找到 %d 个依赖:\n", len(deps))
for _, dep := range deps {
    fmt.Printf("  %s:%s:%s (%s)\n", 
        dep.Group, dep.Name, dep.Version, dep.Scope)
}
```

### 按作用域分组依赖

按作用域（implementation、testImplementation 等）组织依赖：

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

### 分析依赖

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

dependencies := result.Project.Dependencies

// 按作用域统计依赖
scopeCount := make(map[string]int)
for _, dep := range dependencies {
    scopeCount[dep.Scope]++
}

fmt.Println("按作用域统计依赖:")
for scope, count := range scopeCount {
    fmt.Printf("  %s: %d\n", scope, count)
}

// 查找特定依赖
for _, dep := range dependencies {
    if dep.Group == "org.springframework" {
        fmt.Printf("找到 Spring 依赖: %s:%s:%s\n", 
            dep.Group, dep.Name, dep.Version)
    }
}
```

## 处理插件

### 提取插件

从 Gradle 文件获取所有插件：

```go
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("找到 %d 个插件:\n", len(plugins))
for _, plugin := range plugins {
    fmt.Printf("  %s", plugin.ID)
    if plugin.Version != "" {
        fmt.Printf(" (v%s)", plugin.Version)
    }
    if !plugin.Apply {
        fmt.Printf(" [未应用]")
    }
    fmt.Println()
}
```

### 项目类型检测

根据应用的插件确定项目类型：

```go
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

// 检查项目类型
if api.IsAndroidProject(plugins) {
    fmt.Println("✓ Android 项目")
}

if api.IsKotlinProject(plugins) {
    fmt.Println("✓ Kotlin 项目")
}

if api.IsSpringBootProject(plugins) {
    fmt.Println("✓ Spring Boot 项目")
}

// 检查特定插件
hasJavaPlugin := false
for _, plugin := range plugins {
    if plugin.ID == "java" {
        hasJavaPlugin = true
        break
    }
}

if hasJavaPlugin {
    fmt.Println("✓ Java 项目")
}
```

## 处理仓库

### 提取仓库

获取仓库配置：

```go
repos, err := api.GetRepositories("build.gradle")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("找到 %d 个仓库:\n", len(repos))
for _, repo := range repos {
    fmt.Printf("  %s", repo.Name)
    if repo.URL != "" {
        fmt.Printf(" (%s)", repo.URL)
    }
    fmt.Println()
}
```

### 仓库分析

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

repositories := result.Project.Repositories

// 检查常见仓库
hasMaven := false
hasGoogle := false
hasJCenter := false

for _, repo := range repositories {
    switch repo.Name {
    case "mavenCentral":
        hasMaven = true
    case "google":
        hasGoogle = true
    case "jcenter":
        hasJCenter = true
    }
}

fmt.Println("仓库分析:")
fmt.Printf("  Maven Central: %v\n", hasMaven)
fmt.Printf("  Google: %v\n", hasGoogle)
fmt.Printf("  JCenter: %v\n", hasJCenter)
```

## 处理项目属性

### 访问基本属性

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project

// 基本项目信息
fmt.Printf("项目信息:\n")
fmt.Printf("  名称: %s\n", project.Name)
fmt.Printf("  组: %s\n", project.Group)
fmt.Printf("  版本: %s\n", project.Version)
fmt.Printf("  描述: %s\n", project.Description)

// Java 兼容性
if project.SourceCompatibility != "" {
    fmt.Printf("  源码兼容性: %s\n", project.SourceCompatibility)
}
if project.TargetCompatibility != "" {
    fmt.Printf("  目标兼容性: %s\n", project.TargetCompatibility)
}
```

### 自定义属性

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project

// 访问自定义属性
if len(project.Properties) > 0 {
    fmt.Println("自定义属性:")
    for key, value := range project.Properties {
        fmt.Printf("  %s = %s\n", key, value)
    }
}
```

## 错误处理和验证

### 全面的错误处理

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    // 处理不同类型的错误
    if os.IsNotExist(err) {
        log.Printf("Gradle 文件未找到: %v", err)
        return
    }
    
    if strings.Contains(err.Error(), "permission") {
        log.Printf("权限被拒绝: %v", err)
        return
    }
    
    log.Printf("解析错误: %v", err)
    return
}

// 检查警告
if len(result.Warnings) > 0 {
    fmt.Printf("解析完成，有 %d 个警告:\n", len(result.Warnings))
    for i, warning := range result.Warnings {
        fmt.Printf("  %d. %s\n", i+1, warning)
    }
}

// 验证结果
project := result.Project
if project == nil {
    log.Println("警告: 未找到项目信息")
    return
}

if len(project.Dependencies) == 0 {
    log.Println("警告: 未找到依赖")
}

if len(project.Plugins) == 0 {
    log.Println("警告: 未找到插件")
}
```

### 数据验证

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project

// 验证项目结构
errors := []string{}

if project.Group == "" {
    errors = append(errors, "缺少项目组")
}

if project.Version == "" {
    errors = append(errors, "缺少项目版本")
}

// 检查必需的插件
hasJavaPlugin := false
for _, plugin := range project.Plugins {
    if plugin.ID == "java" || plugin.ID == "java-library" {
        hasJavaPlugin = true
        break
    }
}

if !hasJavaPlugin {
    errors = append(errors, "未找到 Java 插件")
}

// 报告验证结果
if len(errors) > 0 {
    fmt.Println("验证错误:")
    for _, err := range errors {
        fmt.Printf("  - %s\n", err)
    }
} else {
    fmt.Println("✓ 项目结构有效")
}
```

## 性能考虑

### 解析大文件

对于大型 Gradle 文件，考虑使用自定义解析器选项：

```go
// 创建优化设置的解析器
options := api.DefaultOptions()
options.SkipComments = true        // 跳过注释以加快解析速度
options.CollectRawContent = false  // 不存储原始内容以节省内存
options.ParseTasks = false         // 如果不需要，跳过任务解析

parser := api.NewParser(options)
result, err := parser.Parse(gradleContent)
if err != nil {
    log.Fatal(err)
}
```

### 批量处理

处理多个文件时：

```go
// 重用解析器实例
parser := api.NewParser(api.DefaultOptions())

files := []string{"app/build.gradle", "lib/build.gradle", "core/build.gradle"}
for _, file := range files {
    result, err := parser.ParseFile(file)
    if err != nil {
        log.Printf("解析 %s 失败: %v", file, err)
        continue
    }
    
    fmt.Printf("已处理 %s: %d 个依赖\n", 
        file, len(result.Project.Dependencies))
}
```

## 下一步

现在您了解了基础知识，探索更多高级功能：

- [高级功能](./advanced-features.md) - 源码映射、自定义解析等
- [结构化编辑](./structured-editing.md) - 以编程方式修改 Gradle 文件
- [配置选项](./configuration.md) - 自定义解析器行为
- [API 参考](../api/) - 完整的 API 文档
