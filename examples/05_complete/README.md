# 完整功能示例

这个示例演示了gradle-parser库的完整功能集，包括：

- 自定义解析器配置
- 解析单个文件或整个项目
- 支持JSON格式输出
- 全面的解析结果展示

## 运行示例

```
go run main.go
```

这个示例使用了硬编码的配置：
- 分析 `examples/sample_files` 整个项目目录
- 使用常规文本格式输出（非JSON）
- 启用所有解析选项（插件、依赖、仓库等）

## 自定义配置

如果你想修改配置，请编辑 `main.go` 中的以下变量：

```go
// 硬编码配置参数
filePath := "examples/sample_files/build.gradle"  // 单个文件路径
projectDir := "examples/sample_files"             // 项目目录路径
useProjectMode := true                            // 是否分析整个项目
jsonOutput := false                               // 是否以JSON格式输出

// 解析器选项
skipComments := true        // 是否跳过注释
collectRawContent := true   // 是否收集原始内容
parsePlugins := true        // 是否解析插件
parseDependencies := true   // 是否解析依赖
parseRepositories := true   // 是否解析仓库
parseTasks := true          // 是否解析任务
``` 