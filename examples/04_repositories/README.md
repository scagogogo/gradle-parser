# 仓库提取示例

这个示例演示了如何使用gradle-parser库提取和分析Gradle仓库配置，包括：

- 提取所有仓库信息
- 检测特定仓库的使用（JitPack、自定义仓库等）
- 显示仓库凭证信息（如果有）
- 统计仓库信息

## 运行示例

```
go run main.go
```

这个示例使用了硬编码的配置：
- 解析 `examples/sample_files/build.gradle` 文件
- 启用特定仓库检查功能

## 自定义配置

如果你想修改配置，请编辑 `main.go` 中的以下变量：

```go
// 硬编码配置参数，根据需要修改
filePath := "examples/sample_files/build.gradle"  // 要解析的Gradle文件路径
checkSpecial := true                              // 是否检查特定仓库的使用
``` 