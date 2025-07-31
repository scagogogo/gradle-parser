# 配置选项

本指南涵盖如何为不同用例配置 Gradle Parser，包括性能优化、功能选择和自定义解析行为。

## 解析器选项

### 默认配置

默认解析器配置针对一般用途进行了优化：

```go
import "github.com/scagogogo/gradle-parser/pkg/api"

// 使用默认选项
result, err := api.ParseFile("build.gradle")

// 或者显式使用默认选项
options := api.DefaultOptions()
parser := api.NewParser(options)
result, err = parser.Parse(content)
```

### 可用选项

```go
type Options struct {
    SkipComments      bool `json:"skipComments"`      // 跳过注释处理
    CollectRawContent bool `json:"collectRawContent"` // 存储原始文件内容
    ParsePlugins      bool `json:"parsePlugins"`      // 解析插件块
    ParseDependencies bool `json:"parseDependencies"` // 解析依赖块
    ParseRepositories bool `json:"parseRepositories"` // 解析仓库块
    ParseTasks        bool `json:"parseTasks"`        // 解析任务定义
}
```

## 性能优化

### 高性能配置

对于需要最大解析速度的场景：

```go
import "github.com/scagogogo/gradle-parser/pkg/parser"

// 最小解析配置
fastParser := parser.NewParser().
    WithSkipComments(true).           // 跳过注释处理
    WithCollectRawContent(false).     // 不存储原始内容
    WithParsePlugins(false).          // 跳过插件解析
    WithParseRepositories(false).     // 跳过仓库解析
    WithParseTasks(false)             // 跳过任务解析

// 只会解析依赖
result, err := fastParser.ParseFile("build.gradle")
dependencies := result.Project.Dependencies
```

## 最佳实践

1. **选择适当的配置**: 根据您的用例匹配配置
2. **重用解析器实例**: 创建一次解析器并重用于多个文件
3. **验证配置**: 检查冲突或无效的选项
4. **监控性能**: 为您的工作负载基准测试不同的配置
5. **记录配置**: 在代码中明确配置选择

## 下一步

- [API 参考](../api/) - 完整的 API 文档
- [示例](../examples/) - 实用的使用示例
- [高级功能](./advanced-features.md) - 高级解析功能
