# 依赖提取示例

这个示例演示了如何使用gradle-parser库提取和分析Gradle依赖，包括：

- 提取所有依赖信息
- 过滤特定依赖
- 按范围（scope）分组显示依赖
- 统计依赖信息

## 运行示例

```
go run main.go
```

这个示例使用了硬编码的配置：
- 解析 `examples/sample_files/build.gradle` 文件
- 按范围分组显示依赖
- 过滤包含 "org.springframework" 的依赖

## 自定义配置

如果你想修改配置，请编辑 `main.go` 中的以下变量：

```go
// 硬编码配置参数，根据需要修改
filePath := "examples/sample_files/build.gradle"  // 要解析的Gradle文件路径
showScope := true                                 // 是否按范围分组显示依赖
filter := "org.springframework"                   // 过滤依赖，留空表示显示所有依赖
``` 