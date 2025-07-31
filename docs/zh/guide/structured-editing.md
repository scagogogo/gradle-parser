# 结构化编辑

Gradle Parser 提供强大的结构化编辑功能，允许您以编程方式修改 Gradle 构建文件，同时保持格式并最小化差异。

## 快速编辑

### 更新依赖版本

更新依赖版本的最简单方式：

```go
import "github.com/scagogogo/gradle-parser/pkg/api"

// 更新 MySQL 连接器版本
newContent, err := api.UpdateDependencyVersion(
    "build.gradle",
    "mysql",
    "mysql-connector-java", 
    "8.0.31"
)
if err != nil {
    log.Fatal(err)
}

// 写回文件
err = os.WriteFile("build.gradle", []byte(newContent), 0644)
if err != nil {
    log.Fatal(err)
}
```

## 高级编辑

### 使用 GradleEditor 进行多重修改

```go
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// 更新多个依赖
err = editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
if err != nil {
    log.Printf("更新 MySQL 失败: %v", err)
}

err = editor.UpdateDependencyVersion("org.springframework", "spring-core", "5.3.21")
if err != nil {
    log.Printf("更新 Spring 失败: %v", err)
}

// 应用所有更改
newContent, err := editor.ApplyModifications()
if err != nil {
    log.Fatal(err)
}

// 保存到文件
err = os.WriteFile("build.gradle", []byte(newContent), 0644)
if err != nil {
    log.Fatal(err)
}
```

## 下一步

- [配置选项](./configuration.md) - 自定义解析器和编辑器行为
- [API 参考](../api/editor.md) - 完整的编辑器 API 文档
- [示例](../examples/structured-editing.md) - 更多编辑示例
