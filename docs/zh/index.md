---
layout: home

hero:
  name: "Gradle Parser"
  text: "强大的 Gradle 构建文件解析器"
  tagline: "在 Go 中轻松解析、分析和编辑 Gradle 构建文件"
  image:
    src: /logo.svg
    alt: Gradle Parser
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/guide/getting-started
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/scagogogo/gradle-parser

features:
  - icon: 🚀
    title: 快速可靠
    details: 高性能解析 Gradle 构建文件，具有全面的错误处理和验证功能。
  
  - icon: 🔍
    title: 深度分析
    details: 提取依赖、插件、仓库和项目配置的详细信息。
  
  - icon: ✏️
    title: 结构化编辑
    details: 以编程方式修改 Gradle 文件，同时保持格式并最小化差异。
  
  - icon: 📍
    title: 源码映射
    details: 跟踪每个解析元素的精确源位置，实现准确的修改。
  
  - icon: 🌐
    title: 多格式支持
    details: 支持 Gradle 构建文件中的 Groovy DSL 和 Kotlin DSL 语法。
  
  - icon: 🛠️
    title: 高度可配置
    details: 通过灵活的选项自定义解析行为，满足您的特定需求。
---

## 快速示例

```go
package main

import (
    "fmt"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    // 解析 Gradle 文件
    result, err := api.ParseFile("build.gradle")
    if err != nil {
        panic(err)
    }

    // 访问项目信息
    fmt.Printf("项目名称: %s\n", result.Project.Name)
    fmt.Printf("项目版本: %s\n", result.Project.Version)
    
    // 列出依赖
    for _, dep := range result.Project.Dependencies {
        fmt.Printf("依赖: %s:%s:%s (%s)\n", 
            dep.Group, dep.Name, dep.Version, dep.Scope)
    }
}
```

## 安装

```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

## 主要特性

### 🔍 **全面解析**
- 提取项目元数据（组、名称、版本、描述）
- 解析带有作用域分类的依赖
- 分析插件配置并检测项目类型
- 处理包括自定义仓库在内的仓库配置

### ✏️ **结构化编辑**
- 以编程方式更新依赖版本
- 修改插件版本和配置
- 编辑项目属性同时保持格式
- 添加新依赖并正确放置

### 📍 **源位置跟踪**
- 每个元素的精确行和列信息
- 实现最小差异的准确修改
- 支持复杂的编辑场景

### 🌐 **多语言支持**
- 完全支持 Groovy DSL 语法
- Kotlin DSL 兼容性
- 处理单模块和多模块项目

## 使用场景

- **构建工具集成**: 将 Gradle 解析集成到您的构建工具和 IDE 中
- **依赖管理**: 以编程方式分析和更新项目依赖
- **项目分析**: 从 Gradle 项目中提取见解用于报告和分析
- **自动化维护**: 批量更新跨项目的依赖和配置
- **迁移工具**: 构建工具以在 Gradle 版本或配置之间迁移

## 社区

- [GitHub 仓库](https://github.com/scagogogo/gradle-parser)
- [问题跟踪](https://github.com/scagogogo/gradle-parser/issues)
- [讨论区](https://github.com/scagogogo/gradle-parser/discussions)

## 许可证

Gradle Parser 基于 [MIT 许可证](https://github.com/scagogogo/gradle-parser/blob/main/LICENSE) 发布。
