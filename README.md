# Gradle Parser | Gradle 解析器

[![CI](https://github.com/scagogogo/gradle-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/gradle-parser/actions/workflows/ci.yml)
[![Quality](https://github.com/scagogogo/gradle-parser/actions/workflows/quality.yml/badge.svg)](https://github.com/scagogogo/gradle-parser/actions/workflows/quality.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/gradle-parser)](https://goreportcard.com/report/github.com/scagogogo/gradle-parser)
[![GoDoc](https://godoc.org/github.com/scagogogo/gradle-parser?status.svg)](https://pkg.go.dev/github.com/scagogogo/gradle-parser)
[![codecov](https://codecov.io/gh/scagogogo/gradle-parser/branch/main/graph/badge.svg)](https://codecov.io/gh/scagogogo/gradle-parser)
[![License](https://img.shields.io/github/license/scagogogo/gradle-parser)](/LICENSE)
[![Release](https://img.shields.io/github/v/release/scagogogo/gradle-parser)](https://github.com/scagogogo/gradle-parser/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/scagogogo/gradle-parser)](https://golang.org/)

一个用于解析Gradle构建文件的Go库，可提取依赖信息、插件配置、仓库设置等。 

A Go library for parsing Gradle build files, extracting dependencies, plugins, repositories and other configuration data.

## 功能特点 | Features

- 🚀 解析Gradle构建文件 (支持build.gradle和build.gradle.kts)
- 🔍 深入提取项目信息（组、名称、版本、描述）
- 📦 解析和分组依赖信息，支持作用域分类
- 🔌 提取插件配置，检测项目类型（Android/Kotlin/Spring Boot）
- 📝 解析仓库配置，包括自定义仓库和认证信息
- ✏️ **结构化编辑**: 精确修改依赖版本、插件版本、项目属性
- 🎯 **最小diff**: 保持原始格式，只修改必要部分
- 📍 **源码位置追踪**: 记录每个元素在源文件中的精确位置
- 🌐 支持Groovy DSL和Kotlin DSL语法
- 🛠️ 可自定义解析器配置，满足不同需求
- 🔄 支持解析整个多模块Gradle项目

## 安装 | Installation

```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

## 快速开始 | Quick Start

```go
package main

import (
    "fmt"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    // 解析Gradle文件
    // Parse Gradle file
    result, err := api.ParseFile("path/to/build.gradle")
    if err != nil {
        panic(err)
    }

    // 获取项目信息
    // Get project info
    fmt.Printf("项目名称: %s\n", result.Project.Name)
    fmt.Printf("项目版本: %s\n", result.Project.Version)
    
    // 获取依赖
    // Get dependencies
    for _, dep := range result.Project.Dependencies {
        fmt.Printf("依赖: %s:%s:%s (%s)\n", 
            dep.Group, dep.Name, dep.Version, dep.Scope)
    }
    
    // 获取插件
    // Get plugins
    for _, plugin := range result.Project.Plugins {
        fmt.Printf("插件: %s (版本: %s)\n", plugin.ID, plugin.Version)
    }
    
    // 获取仓库
    // Get repositories
    for _, repo := range result.Project.Repositories {
        fmt.Printf("仓库: %s (%s)\n", repo.Name, repo.URL)
    }
}
```

## 主要功能 | Main Features

### 依赖提取 | Dependency Extraction

```go
// 直接提取依赖
// Extract dependencies directly
dependencies, err := api.GetDependencies("build.gradle")

// 按范围分组
// Group by scope
dependencySets := api.DependenciesByScope(dependencies)
```

### 插件分析 | Plugin Analysis

```go
// 提取插件信息
// Extract plugins
plugins, err := api.GetPlugins("build.gradle")

// 检测项目类型
// Detect project type
isAndroid := api.IsAndroidProject(plugins)
isKotlin := api.IsKotlinProject(plugins)
isSpringBoot := api.IsSpringBootProject(plugins)
```

### 仓库解析 | Repository Parsing

```go
// 提取仓库配置
// Extract repositories
repos, err := api.GetRepositories("build.gradle")
```

### 自定义解析器 | Custom Parser

```go
// 自定义解析器选项
// Customize parser options
options := api.DefaultOptions()
options.SkipComments = true
options.ParsePlugins = true 
options.ParseDependencies = true

// 创建定制解析器
// Create custom parser
parser := api.NewParser(options)
result, err := parser.ParseFile("build.gradle")
```

### 结构化编辑 | Structured Editing

```go
// 更新依赖版本（便捷方法）
// Update dependency version (convenient method)
newText, err := api.UpdateDependencyVersion("build.gradle", "mysql", "mysql-connector-java", "8.0.31")

// 更新插件版本
// Update plugin version
newText, err := api.UpdatePluginVersion("build.gradle", "org.springframework.boot", "2.7.2")

// 批量修改（高级用法）
// Batch modifications (advanced usage)
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// 执行多个修改
// Perform multiple modifications
editor.UpdateProperty("version", "1.0.0")
editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
editor.UpdateDependencyVersion("com.google.guava", "guava", "31.1-jre")
editor.AddDependency("org.apache.commons", "commons-text", "1.9", "implementation")

// 应用所有修改
// Apply all modifications
serializer := editor.NewGradleSerializer(editor.GetSourceMappedProject().OriginalText)
finalText, err := serializer.ApplyModifications(editor.GetModifications())

// 生成修改diff
// Generate modification diff
diffLines := serializer.GenerateDiff(editor.GetModifications())
for _, diffLine := range diffLines {
    fmt.Println(diffLine.String())
}
```

## 项目结构 | Project Structure

整个项目采用模块化设计，代码组织如下：

The project uses a modular design with the following code organization:

```
├── pkg/                  # 核心包目录 | Core packages
│   ├── api/              # 主API接口 | Main API
│   ├── config/           # 配置解析相关 | Configuration parsing
│   ├── dependency/       # 依赖解析相关 | Dependency parsing
│   ├── editor/           # 结构化编辑器 | Structured editor
│   ├── model/            # 数据模型定义 | Data models
│   ├── parser/           # 解析器核心 | Parser core
│   └── util/             # 工具函数 | Utility functions
└── examples/             # 示例代码 | Example code
    ├── 01_basic/         # 基本使用示例 | Basic usage
    ├── 02_dependencies/  # 依赖提取示例 | Dependency extraction
    ├── 03_plugins/       # 插件提取示例 | Plugin extraction
    ├── 04_repositories/  # 仓库提取示例 | Repository extraction
    ├── 05_complete/      # 完整功能示例 | Complete features
    ├── 06_editor/        # 结构化编辑示例 | Structured editing
    └── sample_files/     # 示例Gradle文件 | Sample Gradle files
```

## 示例程序 | Examples

查看 [examples](examples/) 目录获取更详细的示例代码。

Check the [examples](examples/) directory for more detailed example code.

每个示例程序都展示了库的不同功能，从基本解析到完整的项目分析。

Each example demonstrates different features of the library, from basic parsing to complete project analysis.

## 文档 | Documentation

📚 **完整文档**: [https://scagogogo.github.io/gradle-parser/](https://scagogogo.github.io/gradle-parser/)

**Complete Documentation**: [https://scagogogo.github.io/gradle-parser/](https://scagogogo.github.io/gradle-parser/)

- [快速开始指南 | Getting Started Guide](https://scagogogo.github.io/gradle-parser/guide/getting-started.html)
- [API 参考 | API Reference](https://scagogogo.github.io/gradle-parser/api/)
- [示例代码 | Examples](https://scagogogo.github.io/gradle-parser/examples/)
- [高级功能 | Advanced Features](https://scagogogo.github.io/gradle-parser/guide/advanced-features.html)

## 测试 | Testing

本项目包含全面的测试套件：

This project includes a comprehensive test suite:

```bash
# 运行所有测试 | Run all tests
go test ./...

# 运行测试套件 | Run test suite
cd test && ./scripts/run-tests.sh

# 运行示例 | Run examples
cd examples && ./run-all-examples.sh
```

**测试覆盖率 | Test Coverage**: 目标 >90% | Target >90%

## 持续集成 | Continuous Integration

本项目使用GitHub Actions进行持续集成，确保代码质量和功能正常：

This project uses GitHub Actions for continuous integration to ensure code quality:

### CI 工作流 | CI Workflows

- **🔄 CI**: 多Go版本测试、代码检查、示例验证
- **📊 Quality**: 代码覆盖率、安全扫描、复杂度分析
- **📚 Docs**: 文档构建和部署
- **🚀 Release**: 自动发布和资产构建

### 质量保证 | Quality Assurance

- ✅ 单元测试和集成测试 | Unit and integration tests
- 🔍 代码质量检查 (golangci-lint) | Code quality checks
- 🛡️ 安全漏洞扫描 | Security vulnerability scanning
- 📈 性能基准测试 | Performance benchmarking
- 📝 文档链接验证 | Documentation link validation

## 贡献 | Contributing

欢迎贡献代码！请查看 [贡献指南](CONTRIBUTING.md) 了解详情。

Contributions are welcome! Please see the [Contributing Guide](CONTRIBUTING.md) for details.

### 开发环境 | Development Environment

```bash
# 克隆仓库 | Clone repository
git clone https://github.com/scagogogo/gradle-parser.git
cd gradle-parser

# 安装依赖 | Install dependencies
go mod download

# 运行测试 | Run tests
go test ./...

# 运行示例 | Run examples
cd examples/01_basic && go run main.go
```

### 报告问题 | Reporting Issues

如果您发现了bug或有功能请求，请在 [GitHub Issues](https://github.com/scagogogo/gradle-parser/issues) 中报告。

If you find a bug or have a feature request, please report it in [GitHub Issues](https://github.com/scagogogo/gradle-parser/issues).

## 路线图 | Roadmap

- [ ] 支持更多Gradle DSL语法 | Support more Gradle DSL syntax
- [ ] 增强Kotlin DSL支持 | Enhanced Kotlin DSL support
- [ ] 性能优化 | Performance optimizations
- [ ] 更多编辑功能 | More editing capabilities
- [ ] IDE插件支持 | IDE plugin support

## 协议 | License

MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

MIT License - See [LICENSE](LICENSE) file for details

---

<div align="center">

**⭐ 如果这个项目对您有帮助，请给个星标！**

**⭐ If this project helps you, please give it a star!**

Made with ❤️ by [scagogogo](https://github.com/scagogogo)

</div>