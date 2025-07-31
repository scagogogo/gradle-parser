# 自定义解析器示例

为特定用例配置解析器和性能优化的示例。

## 性能优化

### 快速依赖提取

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/parser"
)

func main() {
    // 创建仅用于依赖提取的优化解析器
    fastParser := parser.NewParser().
        WithSkipComments(true).           // 跳过注释
        WithCollectRawContent(false).     // 不存储原始内容
        WithParsePlugins(false).          // 跳过插件
        WithParseRepositories(false).     // 跳过仓库
        WithParseTasks(false)             // 跳过任务

    result, err := fastParser.ParseFile("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    // 只解析依赖
    dependencies := result.Project.Dependencies
    fmt.Printf("在优化解析中找到 %d 个依赖\n", len(dependencies))
    
    for _, dep := range dependencies {
        fmt.Printf("  %s:%s:%s (%s)\n", 
            dep.Group, dep.Name, dep.Version, dep.Scope)
    }
}
```

## 自定义配置

### 内存优化解析器

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/parser"
)

func main() {
    // 内存高效配置
    memoryParser := parser.NewParser().
        WithCollectRawContent(false).     // 不存储原始内容
        WithSkipComments(true).           // 跳过注释
        WithParseTasks(false)             // 如果不需要，跳过任务

    result, err := memoryParser.ParseFile("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    project := result.Project
    fmt.Printf("内存优化解析完成:\n")
    fmt.Printf("  依赖: %d\n", len(project.Dependencies))
    fmt.Printf("  插件: %d\n", len(project.Plugins))
    fmt.Printf("  仓库: %d\n", len(project.Repositories))
}
```

## 下一步

- [基础解析示例](./basic-parsing.md)
- [依赖分析示例](./dependency-analysis.md)
