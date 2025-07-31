# API 参考

欢迎来到 Gradle Parser API 参考。本节为所有公共 API、类型和函数提供全面的文档。

## 概览

Gradle Parser 为解析和操作 Gradle 构建文件提供了清洁直观的 API。主要入口点是 `api` 包，它为常见操作提供高级函数。

## 包结构

- **[核心 API](./core.md)** - 主要解析函数和工具
- **[解析器](./parser.md)** - 低级解析接口和实现
- **[数据模型](./models.md)** - 表示 Gradle 项目的数据结构
- **[编辑器](./editor.md)** - 结构化编辑和修改功能
- **[工具函数](./utilities.md)** - 辅助函数和工具

## 快速参考

### 核心函数

```go
// 基础解析
result, err := api.ParseFile("build.gradle")
result, err := api.ParseString(content)
result, err := api.ParseReader(reader)

// 提取特定组件
deps, err := api.GetDependencies("build.gradle")
plugins, err := api.GetPlugins("build.gradle")
repos, err := api.GetRepositories("build.gradle")

// 项目类型检测
isAndroid := api.IsAndroidProject(plugins)
isKotlin := api.IsKotlinProject(plugins)
isSpringBoot := api.IsSpringBootProject(plugins)

// 结构化编辑
newText, err := api.UpdateDependencyVersion("build.gradle", "group", "name", "version")
newText, err := api.UpdatePluginVersion("build.gradle", "plugin", "version")
```

## 常见模式

### 错误处理

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Printf("解析失败: %v", err)
    return
}

// 检查警告
if len(result.Warnings) > 0 {
    for _, warning := range result.Warnings {
        log.Printf("警告: %s", warning)
    }
}
```

## 版本兼容性

- **Go 版本**: 需要 Go 1.19 或更高版本
- **Gradle 兼容性**: 支持 Gradle 4.0+ 语法
- **DSL 支持**: 完全支持 Groovy DSL，基本支持 Kotlin DSL
