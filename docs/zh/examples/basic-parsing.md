# 基础解析示例

本页提供基础 Gradle 文件解析操作的实用示例。

## 简单文件解析

### 解析并显示项目信息

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    // 解析 build.gradle 文件
    result, err := api.ParseFile("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    project := result.Project

    // 显示基本项目信息
    fmt.Println("=== 项目信息 ===")
    fmt.Printf("名称: %s\n", project.Name)
    fmt.Printf("组: %s\n", project.Group)
    fmt.Printf("版本: %s\n", project.Version)
    fmt.Printf("描述: %s\n", project.Description)

    // 显示解析统计
    fmt.Printf("\n=== 解析统计 ===\n")
    fmt.Printf("依赖: %d\n", len(project.Dependencies))
    fmt.Printf("插件: %d\n", len(project.Plugins))
    fmt.Printf("仓库: %d\n", len(project.Repositories))
    fmt.Printf("解析时间: %s\n", result.ParseTime)

    if len(result.Warnings) > 0 {
        fmt.Printf("警告: %d\n", len(result.Warnings))
    }
}
```

## 错误处理示例

### 健壮的错误处理

```go
package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func parseWithErrorHandling(filePath string) {
    result, err := api.ParseFile(filePath)
    if err != nil {
        // 处理不同类型的错误
        if os.IsNotExist(err) {
            log.Printf("文件未找到: %s", filePath)
            return
        }
        
        if strings.Contains(err.Error(), "permission") {
            log.Printf("权限被拒绝: %s", filePath)
            return
        }
        
        log.Printf("解析错误: %v", err)
        return
    }

    // 检查解析警告
    if len(result.Warnings) > 0 {
        fmt.Printf("⚠️  解析完成，有 %d 个警告:\n", len(result.Warnings))
        for i, warning := range result.Warnings {
            fmt.Printf("  %d. %s\n", i+1, warning)
        }
        fmt.Println()
    }

    // 验证解析结果
    project := result.Project
    if project == nil {
        log.Println("❌ 未找到项目信息")
        return
    }

    fmt.Printf("✅ 成功解析: %s\n", filePath)
    fmt.Printf("   项目: %s v%s\n", project.Name, project.Version)
}

func main() {
    files := []string{
        "build.gradle",
        "app/build.gradle",
        "non-existent.gradle",
    }

    for _, file := range files {
        parseWithErrorHandling(file)
        fmt.Println()
    }
}
```

## 运行示例

1. 创建新的 Go 模块：
```bash
mkdir gradle-parser-examples
cd gradle-parser-examples
go mod init examples
```

2. 安装 Gradle Parser：
```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

3. 创建示例 `build.gradle` 文件

4. 将任何示例代码复制到 `main.go`

5. 运行示例：
```bash
go run main.go
```
