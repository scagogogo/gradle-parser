# 结构化编辑示例

以编程方式修改 Gradle 构建文件同时保持格式的示例。

## 快速更新

### 更新依赖版本

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
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

    fmt.Println("✅ 成功将 MySQL 连接器更新到 v8.0.31")
}
```

## 高级编辑

### 使用编辑器进行批量更新

```go
package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    editor, err := api.CreateGradleEditor("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    // 更新多个依赖
    updates := map[string]string{
        "mysql:mysql-connector-java":        "8.0.31",
        "org.springframework:spring-core":   "5.3.21",
        "com.google.guava:guava":           "31.1-jre",
    }

    for depKey, version := range updates {
        parts := strings.Split(depKey, ":")
        if len(parts) == 2 {
            err = editor.UpdateDependencyVersion(parts[0], parts[1], version)
            if err != nil {
                log.Printf("更新 %s 失败: %v", depKey, err)
            } else {
                fmt.Printf("✅ 将 %s 更新到 %s\n", depKey, version)
            }
        }
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

    fmt.Println("🎉 所有更新成功应用！")
}
```

## 下一步

- [自定义解析器示例](./custom-parser.md)
- [基础解析示例](./basic-parsing.md)
