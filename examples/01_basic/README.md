# 基本使用示例

这个示例演示了如何使用gradle-parser库解析Gradle文件并提取基本信息，包括：

- 项目名称、组、版本和描述
- 依赖列表
- 插件列表
- 仓库配置

## 运行示例

```
go run main.go
```

这个示例使用了硬编码的文件路径 `examples/sample_files/build.gradle`。

如果你想解析不同的文件，请修改 `main.go` 中的 `filePath` 变量：

```go
// 修改此路径以指向你要解析的Gradle文件
filePath := "examples/sample_files/build.gradle"
``` 