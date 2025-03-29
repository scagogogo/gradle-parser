# 插件提取示例

这个示例演示了如何使用gradle-parser库提取和分析Gradle插件，包括：

- 提取所有插件信息
- 检测项目类型（Android、Kotlin、Spring Boot）
- 统计插件信息

## 运行示例

```
go run main.go
```

这个示例使用了硬编码的配置：
- 解析 `examples/sample_files/build.gradle` 文件
- 启用项目类型检测功能

## 自定义配置

如果你想修改配置，请编辑 `main.go` 中的以下变量：

```go
// 硬编码配置参数，根据需要修改
filePath := "examples/sample_files/build.gradle"  // 要解析的Gradle文件路径
detectType := true                                // 是否检测项目类型
``` 