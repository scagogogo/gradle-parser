# Gradle Parser | Gradle 解析器

![build](https://github.com/scagogogo/gradle-parser/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/gradle-parser)](https://goreportcard.com/report/github.com/scagogogo/gradle-parser)
[![GoDoc](https://godoc.org/github.com/scagogogo/gradle-parser?status.svg)](https://pkg.go.dev/github.com/scagogogo/gradle-parser)
[![License](https://img.shields.io/github/license/scagogogo/gradle-parser)](/LICENSE)

一个用于解析Gradle构建文件的Go库，可提取依赖信息、插件配置、仓库设置等。 

A Go library for parsing Gradle build files, extracting dependencies, plugins, repositories and other configuration data.

## 功能特点 | Features

- 🚀 解析Gradle构建文件 (支持build.gradle和build.gradle.kts)
- 🔍 深入提取项目信息（组、名称、版本、描述）
- 📦 解析和分组依赖信息，支持作用域分类
- 🔌 提取插件配置，检测项目类型（Android/Kotlin/Spring Boot）
- 📝 解析仓库配置，包括自定义仓库和认证信息
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

## 项目结构 | Project Structure

整个项目采用模块化设计，代码组织如下：

The project uses a modular design with the following code organization:

```
├── pkg/                  # 核心包目录 | Core packages
│   ├── api/              # 主API接口 | Main API
│   ├── config/           # 配置解析相关 | Configuration parsing
│   ├── dependency/       # 依赖解析相关 | Dependency parsing
│   ├── model/            # 数据模型定义 | Data models
│   ├── parser/           # 解析器核心 | Parser core
│   └── util/             # 工具函数 | Utility functions
└── examples/             # 示例代码 | Example code
    ├── 01_basic/         # 基本使用示例 | Basic usage
    ├── 02_dependencies/  # 依赖提取示例 | Dependency extraction
    ├── 03_plugins/       # 插件提取示例 | Plugin extraction
    ├── 04_repositories/  # 仓库提取示例 | Repository extraction
    ├── 05_complete/      # 完整功能示例 | Complete features
    └── sample_files/     # 示例Gradle文件 | Sample Gradle files
```

## 示例程序 | Examples

查看 [examples](examples/) 目录获取更详细的示例代码。

Check the [examples](examples/) directory for more detailed example code.

每个示例程序都展示了库的不同功能，从基本解析到完整的项目分析。

Each example demonstrates different features of the library, from basic parsing to complete project analysis.

## 持续集成 | Continuous Integration

本项目使用GitHub Actions进行持续集成，确保代码质量和功能正常。每次提交代码时，CI系统会自动：

This project uses GitHub Actions for continuous integration to ensure code quality. On each commit, the CI system automatically:

- 运行所有单元测试 | Runs all unit tests
- 运行所有示例程序 | Runs all example programs
- 确保代码能够正常构建 | Ensures the code builds correctly

## 协议 | License

MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

MIT License - See [LICENSE](LICENSE) file for details 