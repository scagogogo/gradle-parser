# 快速开始

欢迎使用 Gradle Parser！本指南将帮助您开始在 Go 中解析 Gradle 构建文件。

## 安装

使用 Go 模块安装 Gradle Parser：

```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

## 快速开始

这是一个简单的示例来帮助您开始：

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    // 解析 Gradle 文件
    result, err := api.ParseFile("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    // 打印项目信息
    project := result.Project
    fmt.Printf("项目名称: %s\n", project.Name)
    fmt.Printf("项目组: %s\n", project.Group)
    fmt.Printf("项目版本: %s\n", project.Version)
    fmt.Printf("项目描述: %s\n", project.Description)
}
```

## 基础解析

### 从文件解析

解析 Gradle 文件的最常见方式：

```go
result, err := api.ParseFile("path/to/build.gradle")
if err != nil {
    log.Fatal(err)
}

// 访问解析的项目
project := result.Project
```

### 从字符串解析

如果您有 Gradle 内容作为字符串：

```go
gradleContent := `
plugins {
    id 'java'
}

group = 'com.example'
version = '1.0.0'
`

result, err := api.ParseString(gradleContent)
if err != nil {
    log.Fatal(err)
}
```

### 从 Reader 解析

用于流式传输或其他 I/O 源：

```go
import "strings"

reader := strings.NewReader(gradleContent)
result, err := api.ParseReader(reader)
if err != nil {
    log.Fatal(err)
}
```

## 理解解析结果

`ParseResult` 包含：

- **Project**: 主要项目信息和组件
- **RawText**: 原始文件内容（如果启用）
- **Errors**: 遇到的任何解析错误
- **Warnings**: 非致命解析警告
- **ParseTime**: 解析文件所用的时间

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("解析时间: %s\n", result.ParseTime)
fmt.Printf("警告数量: %d\n", len(result.Warnings))

// 访问项目组件
project := result.Project
fmt.Printf("依赖数量: %d\n", len(project.Dependencies))
fmt.Printf("插件数量: %d\n", len(project.Plugins))
fmt.Printf("仓库数量: %d\n", len(project.Repositories))
```

## 错误处理

Gradle Parser 提供详细的错误信息：

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Printf("解析文件失败: %v", err)
    return
}

// 检查解析警告
if len(result.Warnings) > 0 {
    fmt.Println("解析警告:")
    for _, warning := range result.Warnings {
        fmt.Printf("  - %s\n", warning)
    }
}
```

## 下一步

现在您已经掌握了基础知识，探索更多功能：

- [基本用法](./basic-usage.md) - 学习提取特定信息
- [高级功能](./advanced-features.md) - 发现强大的解析功能
- [结构化编辑](./structured-editing.md) - 以编程方式修改 Gradle 文件
- [配置选项](./configuration.md) - 自定义解析器行为

## 系统要求

- Go 1.19 或更高版本
- 无需外部依赖

## 支持的格式

- Groovy DSL (`build.gradle`)
- Kotlin DSL (`build.gradle.kts`) - 基本支持
- 单模块和多模块项目
