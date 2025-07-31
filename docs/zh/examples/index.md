# 示例

本节提供在各种场景中使用 Gradle Parser 的实用示例。每个示例都包含完整的、可运行的代码和解释。

## 可用示例

### 基础用法
- **[基础解析](./basic-parsing.md)** - 简单的文件解析和信息提取
- **[依赖分析](./dependency-analysis.md)** - 处理项目依赖
- **[插件检测](./plugin-detection.md)** - 分析插件和项目类型
- **[仓库解析](./repository-parsing.md)** - 提取仓库配置

### 高级功能
- **[结构化编辑](./structured-editing.md)** - 以编程方式修改 Gradle 文件
- **[自定义解析器](./custom-parser.md)** - 为特定需求配置解析器选项

## 运行示例

所有示例都设计为自包含且可运行。要运行示例：

1. 创建新的 Go 模块：
```bash
mkdir gradle-parser-example
cd gradle-parser-example
go mod init example
```

2. 安装 Gradle Parser：
```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

3. 将示例代码复制到 `main.go`

4. 创建示例 `build.gradle` 文件（或使用提供的示例）

5. 运行示例：
```bash
go run main.go
```

## 示例分类

### 🔍 **解析和分析**
学习如何解析 Gradle 文件并提取信息：
- 项目元数据（名称、版本、组）
- 带作用域分析的依赖
- 插件配置
- 仓库设置

### ✏️ **编辑和修改**
了解如何以编程方式修改 Gradle 文件：
- 更新依赖版本
- 修改插件配置
- 添加新依赖
- 保持格式并最小化差异

### 🛠️ **高级用法**
探索高级功能和自定义：
- 自定义解析器配置
- 源位置跟踪
- 错误处理策略
- 性能优化

## 常见用例

### 构建工具集成
```go
// 解析和分析项目
result, err := api.ParseFile("build.gradle")
if err != nil {
    return err
}

// 检查过时的依赖
for _, dep := range result.Project.Dependencies {
    if isOutdated(dep) {
        fmt.Printf("过时: %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
    }
}
```

### 依赖管理
```go
// 更新所有 Spring Boot 依赖
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    return err
}

springBootVersion := "2.7.2"
err = editor.UpdatePluginVersion("org.springframework.boot", springBootVersion)
if err != nil {
    return err
}
```

## 获取帮助

如果您对示例有疑问或需要特定用例的帮助：

- 查看 [API 参考](../api/) 获取详细文档
- 访问我们的 [GitHub 讨论](https://github.com/scagogogo/gradle-parser/discussions)
- 为错误或功能请求开启 [Issue](https://github.com/scagogogo/gradle-parser/issues)
