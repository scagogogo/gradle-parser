# 插件检测示例

本页提供检测和分析 Gradle 项目中插件的示例。

## 基础插件检测

### 列出所有插件

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    plugins, err := api.GetPlugins("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("找到 %d 个插件:\n\n", len(plugins))
    
    for i, plugin := range plugins {
        fmt.Printf("%d. %s", i+1, plugin.ID)
        if plugin.Version != "" {
            fmt.Printf(" (v%s)", plugin.Version)
        }
        if !plugin.Apply {
            fmt.Printf(" [未应用]")
        }
        fmt.Println()
    }
}
```

## 项目类型检测

### 检测项目类型

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    plugins, err := api.GetPlugins("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("🔍 分析项目类型...")
    
    projectTypes := []string{}
    
    if api.IsAndroidProject(plugins) {
        projectTypes = append(projectTypes, "📱 Android")
    }
    
    if api.IsKotlinProject(plugins) {
        projectTypes = append(projectTypes, "🎯 Kotlin")
    }
    
    if api.IsSpringBootProject(plugins) {
        projectTypes = append(projectTypes, "🍃 Spring Boot")
    }
    
    // 检查其他常见项目类型
    for _, plugin := range plugins {
        switch plugin.ID {
        case "java":
            projectTypes = append(projectTypes, "☕ Java")
        case "application":
            projectTypes = append(projectTypes, "🚀 应用程序")
        case "java-library":
            projectTypes = append(projectTypes, "📚 Java 库")
        }
    }
    
    if len(projectTypes) > 0 {
        fmt.Printf("\n✅ 检测到的项目类型:\n")
        for _, pType := range projectTypes {
            fmt.Printf("   %s\n", pType)
        }
    } else {
        fmt.Println("\n❓ 未知项目类型")
    }
}
```

## 下一步

- [仓库解析示例](./repository-parsing.md)
- [结构化编辑示例](./structured-editing.md)
- [自定义解析器示例](./custom-parser.md)
