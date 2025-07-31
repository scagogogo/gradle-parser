# Gradle Parser

[![CI](https://github.com/scagogogo/gradle-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/gradle-parser/actions/workflows/ci.yml)
[![Quality](https://github.com/scagogogo/gradle-parser/actions/workflows/quality.yml/badge.svg)](https://github.com/scagogogo/gradle-parser/actions/workflows/quality.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/gradle-parser)](https://goreportcard.com/report/github.com/scagogogo/gradle-parser)
[![GoDoc](https://godoc.org/github.com/scagogogo/gradle-parser?status.svg)](https://pkg.go.dev/github.com/scagogogo/gradle-parser)
[![codecov](https://codecov.io/gh/scagogogo/gradle-parser/branch/main/graph/badge.svg)](https://codecov.io/gh/scagogogo/gradle-parser)
[![License](https://img.shields.io/github/license/scagogogo/gradle-parser)](/LICENSE)
[![Release](https://img.shields.io/github/v/release/scagogogo/gradle-parser)](https://github.com/scagogogo/gradle-parser/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/scagogogo/gradle-parser)](https://golang.org/)

一个功能强大的Go库，用于解析和操作Gradle构建文件。轻松提取依赖项、插件、仓库和其他配置数据。具有结构化编辑功能，可进行程序化修改，同时保持原始格式。

## 📚 文档

> **🌟 完整文档：[https://scagogogo.github.io/gradle-parser/](https://scagogogo.github.io/gradle-parser/)**
>
> 📖 [English README](README.md) | 🚀 [快速开始](https://scagogogo.github.io/gradle-parser/zh/guide/getting-started.html) | 📋 [API 参考](https://scagogogo.github.io/gradle-parser/zh/api/) | 💡 [示例代码](https://scagogogo.github.io/gradle-parser/zh/examples/)

## ✨ 功能特点

### 🔍 **全面解析**
- 解析Gradle构建文件（支持 `build.gradle` 和 `build.gradle.kts`）
- 提取项目元数据（组、名称、版本、描述）
- 解析和分类依赖项，支持作用域分类
- 分析插件配置，检测项目类型（Android/Kotlin/Spring Boot）
- 处理仓库配置，包括自定义仓库和身份验证

### ✏️ **结构化编辑**
- **精确修改**：更新依赖版本、插件版本、项目属性
- **最小差异**：保持原始格式，仅修改必要部分
- **源码位置跟踪**：记录元素在源文件中的精确位置
- **批量操作**：在单个操作中应用多个修改

### 🛠️ **高级功能**
- 支持Groovy DSL和Kotlin DSL语法
- 可自定义解析器配置以满足不同需求
- 多模块Gradle项目支持
- 全面的错误处理和验证

## 🚀 快速开始

### 安装

```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    // 解析Gradle文件
    result, err := api.ParseFile("path/to/build.gradle")
    if err != nil {
        panic(err)
    }

    // 访问项目信息
    project := result.Project
    fmt.Printf("项目：%s\n", project.Name)
    fmt.Printf("组：%s\n", project.Group)
    fmt.Printf("版本：%s\n", project.Version)
    
    // 列出依赖项
    for _, dep := range project.Dependencies {
        fmt.Printf("依赖：%s:%s:%s (%s)\n", 
            dep.Group, dep.Name, dep.Version, dep.Scope)
    }
    
    // 列出插件
    for _, plugin := range project.Plugins {
        fmt.Printf("插件：%s", plugin.ID)
        if plugin.Version != "" {
            fmt.Printf(" v%s", plugin.Version)
        }
        fmt.Println()
    }
    
    // 列出仓库
    for _, repo := range project.Repositories {
        fmt.Printf("仓库：%s (%s)\n", repo.Name, repo.URL)
    }
}
```

## 📖 核心功能

### 🔍 **依赖分析**

```go
// 直接提取依赖项
dependencies, err := api.GetDependencies("build.gradle")
if err != nil {
    log.Fatal(err)
}

// 按作用域分组依赖项
dependencySets := api.DependenciesByScope(dependencies)
for _, set := range dependencySets {
    fmt.Printf("作用域：%s\n", set.Scope)
    for _, dep := range set.Dependencies {
        fmt.Printf("  %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
    }
}
```

### 🔌 **插件检测**

```go
// 提取插件信息
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

// 检测项目类型
if api.IsAndroidProject(plugins) {
    fmt.Println("检测到Android项目")
}
if api.IsKotlinProject(plugins) {
    fmt.Println("检测到Kotlin项目")
}
if api.IsSpringBootProject(plugins) {
    fmt.Println("检测到Spring Boot项目")
}
```

### 📝 **仓库配置**

```go
// 提取仓库配置
repos, err := api.GetRepositories("build.gradle")
if err != nil {
    log.Fatal(err)
}

for _, repo := range repos {
    fmt.Printf("仓库：%s (%s)\n", repo.Name, repo.URL)
}
```

### ✏️ **结构化编辑**

```go
// 简单版本更新
newText, err := api.UpdateDependencyVersion("build.gradle", "mysql", "mysql-connector-java", "8.0.31")
if err != nil {
    log.Fatal(err)
}

newText, err = api.UpdatePluginVersion("build.gradle", "org.springframework.boot", "2.7.2")
if err != nil {
    log.Fatal(err)
}
```

### 🛠️ **高级编辑**

```go
// 创建编辑器进行批量修改
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// 执行多个修改
editor.UpdateProperty("version", "1.0.0")
editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
editor.UpdateDependencyVersion("com.google.guava", "guava", "31.1-jre")
editor.AddDependency("org.apache.commons", "commons-text", "1.9", "implementation")

// 应用所有修改
serializer := editor.NewGradleSerializer(editor.GetSourceMappedProject().OriginalText)
finalText, err := serializer.ApplyModifications(editor.GetModifications())
if err != nil {
    log.Fatal(err)
}

// 生成差异以供审查
diffLines := serializer.GenerateDiff(editor.GetModifications())
for _, diffLine := range diffLines {
    fmt.Println(diffLine.String())
}
```

### ⚙️ **自定义解析器配置**

```go
// 创建自定义解析器选项
options := api.DefaultOptions()
options.SkipComments = true
options.ParsePlugins = true 
options.ParseDependencies = true
options.ParseRepositories = true

// 创建自定义解析器
parser := api.NewParser(options)
result, err := parser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}
```

## 🏗️ 项目结构

项目采用模块化设计，关注点清晰分离：

```
├── pkg/                  # 核心包
│   ├── api/              # 主API接口
│   ├── config/           # 配置解析
│   ├── dependency/       # 依赖解析
│   ├── editor/           # 结构化编辑器
│   ├── model/            # 数据模型
│   ├── parser/           # 解析器核心
│   └── util/             # 工具函数
└── examples/             # 示例代码
    ├── 01_basic/         # 基本用法
    ├── 02_dependencies/  # 依赖提取
    ├── 03_plugins/       # 插件提取
    ├── 04_repositories/  # 仓库提取
    ├── 05_complete/      # 完整功能
    ├── 06_editor/        # 结构化编辑
    └── sample_files/     # 示例Gradle文件
```

## 📚 资源

### 📖 **示例**
探索 [examples](examples/) 目录，查看展示不同功能的综合代码示例：

- **基本用法**：简单解析和数据提取
- **高级功能**：复杂解析场景和自定义
- **结构化编辑**：程序化文件修改
- **项目分析**：完整的项目检查和报告

### 🧪 **测试**
运行全面的测试套件：

```bash
# 运行所有测试
go test ./...

# 运行带覆盖率的测试套件
cd test && ./scripts/run-tests.sh

# 运行示例
cd examples && ./run-all-examples.sh
```

**测试覆盖率**：目标 >90%

### 🔄 **持续集成**
本项目使用GitHub Actions进行全面的质量保证：

- **🔄 CI**：多版本Go测试、代码验证、示例验证
- **📊 Quality**：代码覆盖率、安全扫描、复杂度分析
- **📚 Docs**：文档构建和部署
- **🚀 Release**：自动发布和资产构建

**质量标准：**
- ✅ 单元测试和集成测试
- 🔍 代码质量检查（golangci-lint）
- 🛡️ 安全漏洞扫描
- 📈 性能基准测试
- 📝 文档链接验证

## 🤝 贡献

欢迎贡献！请查看[贡献指南](CONTRIBUTING.md)了解详情。

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/scagogogo/gradle-parser.git
cd gradle-parser

# 安装依赖
go mod download

# 运行测试
go test ./...

# 尝试示例
cd examples/01_basic && go run main.go
```

### 报告问题

发现了bug或有功能请求？请在[GitHub Issues](https://github.com/scagogogo/gradle-parser/issues)中报告。

## 🗺️ 路线图

- [ ] 增强Kotlin DSL支持
- [ ] 性能优化
- [ ] 更多编辑功能
- [ ] IDE插件支持
- [ ] 额外的Gradle DSL语法支持

## 📄 许可证

MIT许可证 - 详见[LICENSE](LICENSE)文件。

---

<div align="center">

**⭐ 如果这个项目对您有帮助，请给个星标！**

Made with ❤️ by [scagogogo](https://github.com/scagogogo)

</div>
