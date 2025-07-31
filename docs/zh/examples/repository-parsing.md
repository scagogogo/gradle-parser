# 仓库解析示例

解析和分析 Gradle 项目中仓库配置的示例。

## 基础仓库提取

### 列出所有仓库

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    repos, err := api.GetRepositories("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("找到 %d 个仓库:\n\n", len(repos))
    
    for i, repo := range repos {
        fmt.Printf("%d. %s", i+1, repo.Name)
        if repo.URL != "" {
            fmt.Printf(" (%s)", repo.URL)
        }
        fmt.Printf(" [%s]\n", repo.Type)
    }
}
```

## 仓库分析

### 分析仓库类型

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    repos, err := api.GetRepositories("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("📊 仓库分析")
    fmt.Println("======================\n")

    // 按类型计数
    typeCount := make(map[string]int)
    for _, repo := range repos {
        typeCount[repo.Type]++
    }

    fmt.Println("仓库类型:")
    for repoType, count := range typeCount {
        fmt.Printf("  %s: %d\n", repoType, count)
    }

    // 检查常见仓库
    commonRepos := map[string]bool{
        "mavenCentral": false,
        "google":       false,
        "jcenter":      false,
    }

    for _, repo := range repos {
        if _, exists := commonRepos[repo.Name]; exists {
            commonRepos[repo.Name] = true
        }
    }

    fmt.Println("\n常见仓库:")
    for name, found := range commonRepos {
        status := "❌"
        if found {
            status = "✅"
        }
        fmt.Printf("  %s %s\n", status, name)
    }
}
```

## 下一步

- [结构化编辑示例](./structured-editing.md)
- [自定义解析器示例](./custom-parser.md)
