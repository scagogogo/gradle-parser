# 依赖分析示例

本页提供分析 Gradle 项目中依赖的示例。

## 基础依赖提取

### 列出所有依赖

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    deps, err := api.GetDependencies("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("找到 %d 个依赖:\n\n", len(deps))
    
    for i, dep := range deps {
        fmt.Printf("%d. %s:%s:%s (%s)\n", 
            i+1, dep.Group, dep.Name, dep.Version, dep.Scope)
    }
}
```

## 依赖分组

### 按作用域分组

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    deps, err := api.GetDependencies("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    depSets := api.DependenciesByScope(deps)
    
    for _, depSet := range depSets {
        fmt.Printf("\n=== %s 依赖 (%d) ===\n", 
            depSet.Scope, len(depSet.Dependencies))
        
        for _, dep := range depSet.Dependencies {
            fmt.Printf("  %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
        }
    }
}
```

## 依赖分析

### 查找过时的依赖

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

// 模拟函数检查版本是否过时
func isOutdated(dep *api.Dependency) (bool, string) {
    // 这是一个简化的示例
    // 在实际使用中，您会检查依赖数据库
    outdatedVersions := map[string]string{
        "mysql:mysql-connector-java": "8.0.31",
        "org.springframework:spring-core": "5.3.23",
        "com.google.guava:guava": "31.1-jre",
    }
    
    key := fmt.Sprintf("%s:%s", dep.Group, dep.Name)
    if latestVersion, exists := outdatedVersions[key]; exists {
        if dep.Version != latestVersion {
            return true, latestVersion
        }
    }
    
    return false, ""
}

func main() {
    deps, err := api.GetDependencies("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("检查过时的依赖...\n")
    
    outdatedCount := 0
    for _, dep := range deps {
        if outdated, latestVersion := isOutdated(dep); outdated {
            fmt.Printf("⚠️  %s:%s\n", dep.Group, dep.Name)
            fmt.Printf("   当前: %s → 最新: %s\n", dep.Version, latestVersion)
            outdatedCount++
        }
    }
    
    if outdatedCount == 0 {
        fmt.Println("✅ 所有依赖都是最新的！")
    } else {
        fmt.Printf("\n📊 找到 %d 个过时的依赖\n", outdatedCount)
    }
}
```

## 下一步

- [插件检测示例](./plugin-detection.md)
- [仓库解析示例](./repository-parsing.md)
- [结构化编辑示例](./structured-editing.md)
