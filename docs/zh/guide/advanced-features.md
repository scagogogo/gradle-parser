# 高级功能

本指南涵盖了 Gradle Parser 的高级功能，包括源码映射、自定义解析配置和专门的使用场景。

## 源码感知解析

源码感知解析跟踪原始文件中每个解析元素的确切位置，实现精确修改和详细分析。

### 基本源码映射

```go
import "github.com/scagogogo/gradle-parser/pkg/api"

// 使用源码映射解析
result, err := api.ParseFileWithSourceMapping("build.gradle")
if err != nil {
    log.Fatal(err)
}

sourceMappedProject := result.SourceMappedProject

// 访问带有源位置的依赖
for _, dep := range sourceMappedProject.SourceMappedDependencies {
    dependency := dep.Dependency
    sourceRange := dep.SourceRange
    
    fmt.Printf("依赖: %s:%s:%s\n", 
        dependency.Group, dependency.Name, dependency.Version)
    fmt.Printf("  位置: 第%d行，第%d列 到 第%d行，第%d列\n",
        sourceRange.Start.Line, sourceRange.Start.Column,
        sourceRange.End.Line, sourceRange.End.Column)
    fmt.Printf("  原始文本: %s\n", dep.RawText)
}
```

## 自定义解析器配置

### 性能优化

为特定用例配置解析器：

```go
import "github.com/scagogogo/gradle-parser/pkg/parser"

// 仅用于依赖提取的高性能配置
fastParser := parser.NewParser().
    WithSkipComments(true).           // 跳过注释处理
    WithCollectRawContent(false).     // 不存储原始内容
    WithParsePlugins(false).          // 跳过插件解析
    WithParseRepositories(false).     // 跳过仓库解析
    WithParseTasks(false)             // 跳过任务解析

result, err := fastParser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

// 只会解析依赖
dependencies := result.Project.Dependencies
fmt.Printf("找到 %d 个依赖\n", len(dependencies))
```

## 下一步

探索更多高级主题：

- [结构化编辑](./structured-editing.md) - 以编程方式修改 Gradle 文件
- [配置选项](./configuration.md) - 详细的解析器配置选项
- [API 参考](../api/) - 完整的 API 文档
- [示例](../examples/) - 实用的使用示例
