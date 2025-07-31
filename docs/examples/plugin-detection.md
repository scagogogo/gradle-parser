# Plugin Detection Examples

This page provides examples of detecting and analyzing plugins in Gradle projects.

## Basic Plugin Detection

### List All Plugins

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

    fmt.Printf("Found %d plugins:\n\n", len(plugins))
    
    for i, plugin := range plugins {
        fmt.Printf("%d. %s", i+1, plugin.ID)
        if plugin.Version != "" {
            fmt.Printf(" (v%s)", plugin.Version)
        }
        if !plugin.Apply {
            fmt.Printf(" [not applied]")
        }
        fmt.Println()
    }
}
```

## Project Type Detection

### Detect Project Types

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

    fmt.Println("🔍 Analyzing project type...")
    
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
    
    // Check for other common project types
    for _, plugin := range plugins {
        switch plugin.ID {
        case "java":
            projectTypes = append(projectTypes, "☕ Java")
        case "application":
            projectTypes = append(projectTypes, "🚀 Application")
        case "java-library":
            projectTypes = append(projectTypes, "📚 Java Library")
        }
    }
    
    if len(projectTypes) > 0 {
        fmt.Printf("\n✅ Detected project types:\n")
        for _, pType := range projectTypes {
            fmt.Printf("   %s\n", pType)
        }
    } else {
        fmt.Println("\n❓ Unknown project type")
    }
}
```

## Advanced Plugin Analysis

### Plugin Categorization

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func categorizePlugins(plugins []*api.Plugin) map[string][]*api.Plugin {
    categories := map[string][]string{
        "Language":    {"java", "kotlin", "groovy", "scala"},
        "Framework":   {"org.springframework.boot", "io.quarkus", "io.micronaut"},
        "Android":     {"com.android.application", "com.android.library"},
        "Build":       {"maven-publish", "signing", "distribution"},
        "Quality":     {"checkstyle", "pmd", "spotbugs", "jacoco"},
        "Development": {"idea", "eclipse", "application"},
    }

    result := make(map[string][]*api.Plugin)
    uncategorized := []*api.Plugin{}

    for _, plugin := range plugins {
        categorized := false
        for category, pluginIds := range categories {
            for _, id := range pluginIds {
                if plugin.ID == id || strings.Contains(plugin.ID, id) {
                    result[category] = append(result[category], plugin)
                    categorized = true
                    break
                }
            }
            if categorized {
                break
            }
        }
        
        if !categorized {
            uncategorized = append(uncategorized, plugin)
        }
    }
    
    if len(uncategorized) > 0 {
        result["Other"] = uncategorized
    }
    
    return result
}

func main() {
    plugins, err := api.GetPlugins("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    categorized := categorizePlugins(plugins)
    
    fmt.Println("📊 Plugin Analysis by Category")
    fmt.Println("==============================\n")
    
    for category, categoryPlugins := range categorized {
        if len(categoryPlugins) > 0 {
            fmt.Printf("🏷️  %s Plugins (%d):\n", category, len(categoryPlugins))
            for _, plugin := range categoryPlugins {
                fmt.Printf("   • %s", plugin.ID)
                if plugin.Version != "" {
                    fmt.Printf(" v%s", plugin.Version)
                }
                fmt.Println()
            }
            fmt.Println()
        }
    }
}
```

## Plugin Version Analysis

### Check for Outdated Plugins

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

// Mock latest versions database
var latestVersions = map[string]string{
    "org.springframework.boot":           "2.7.2",
    "io.spring.dependency-management":    "1.0.12.RELEASE",
    "org.jetbrains.kotlin.jvm":          "1.7.10",
    "com.android.application":            "7.2.1",
    "com.android.library":                "7.2.1",
}

func checkPluginVersions(plugins []*api.Plugin) {
    fmt.Println("🔍 Checking plugin versions...")
    fmt.Println("==============================\n")
    
    outdatedCount := 0
    
    for _, plugin := range plugins {
        if plugin.Version == "" {
            continue // Skip plugins without explicit versions
        }
        
        if latestVersion, exists := latestVersions[plugin.ID]; exists {
            if plugin.Version != latestVersion {
                fmt.Printf("⚠️  %s\n", plugin.ID)
                fmt.Printf("   Current: %s → Latest: %s\n\n", plugin.Version, latestVersion)
                outdatedCount++
            }
        }
    }
    
    if outdatedCount == 0 {
        fmt.Println("✅ All plugins are up to date!")
    } else {
        fmt.Printf("📊 Found %d potentially outdated plugins\n", outdatedCount)
    }
}

func main() {
    plugins, err := api.GetPlugins("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    checkPluginVersions(plugins)
}
```

## Plugin Configuration Analysis

### Extract Plugin Configurations

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func analyzePluginConfigurations(filePath string) {
    result, err := api.ParseFile(filePath)
    if err != nil {
        log.Fatal(err)
    }

    plugins := result.Project.Plugins
    
    fmt.Println("🔧 Plugin Configuration Analysis")
    fmt.Println("================================\n")
    
    for _, plugin := range plugins {
        fmt.Printf("📦 %s", plugin.ID)
        if plugin.Version != "" {
            fmt.Printf(" v%s", plugin.Version)
        }
        fmt.Println()
        
        if len(plugin.Config) > 0 {
            fmt.Println("   Configuration:")
            for key, value := range plugin.Config {
                fmt.Printf("     %s: %v\n", key, value)
            }
        } else {
            fmt.Println("   No explicit configuration")
        }
        
        fmt.Printf("   Applied: %v\n\n", plugin.Apply)
    }
}

func main() {
    analyzePluginConfigurations("build.gradle")
}
```

## Multi-Project Plugin Analysis

### Compare Plugins Across Modules

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func comparePluginsAcrossModules(modules []string) {
    modulePlugins := make(map[string][]*api.Plugin)
    
    // Collect plugins from each module
    for _, module := range modules {
        buildFile := filepath.Join(module, "build.gradle")
        plugins, err := api.GetPlugins(buildFile)
        if err != nil {
            log.Printf("Warning: failed to parse %s: %v", buildFile, err)
            continue
        }
        modulePlugins[module] = plugins
    }
    
    // Find common plugins
    pluginUsage := make(map[string][]string)
    for module, plugins := range modulePlugins {
        for _, plugin := range plugins {
            pluginUsage[plugin.ID] = append(pluginUsage[plugin.ID], module)
        }
    }
    
    fmt.Println("🔍 Multi-Module Plugin Analysis")
    fmt.Println("===============================\n")
    
    // Report shared plugins
    fmt.Println("📊 Shared Plugins:")
    for pluginID, modules := range pluginUsage {
        if len(modules) > 1 {
            fmt.Printf("   %s (used in: %v)\n", pluginID, modules)
        }
    }
    
    // Report module-specific plugins
    fmt.Println("\n📦 Module-Specific Plugins:")
    for pluginID, modules := range pluginUsage {
        if len(modules) == 1 {
            fmt.Printf("   %s (only in: %s)\n", pluginID, modules[0])
        }
    }
}

func main() {
    modules := []string{".", "app", "lib", "core"}
    comparePluginsAcrossModules(modules)
}
```

## Plugin Recommendation System

### Suggest Missing Plugins

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func suggestPlugins(filePath string) {
    result, err := api.ParseFile(filePath)
    if err != nil {
        log.Fatal(err)
    }

    plugins := result.Project.Plugins
    dependencies := result.Project.Dependencies
    
    // Create plugin ID set for quick lookup
    pluginIDs := make(map[string]bool)
    for _, plugin := range plugins {
        pluginIDs[plugin.ID] = true
    }
    
    suggestions := []string{}
    
    // Analyze dependencies to suggest plugins
    for _, dep := range dependencies {
        switch {
        case dep.Group == "org.springframework.boot" && !pluginIDs["org.springframework.boot"]:
            suggestions = append(suggestions, "org.springframework.boot - for Spring Boot projects")
            
        case dep.Group == "org.junit.jupiter" && !pluginIDs["java"]:
            suggestions = append(suggestions, "java - for JUnit 5 support")
            
        case (dep.Group == "mysql" || dep.Group == "org.postgresql") && !pluginIDs["org.flywaydb.flyway"]:
            suggestions = append(suggestions, "org.flywaydb.flyway - for database migrations")
            
        case dep.Name == "kotlin-stdlib" && !pluginIDs["org.jetbrains.kotlin.jvm"]:
            suggestions = append(suggestions, "org.jetbrains.kotlin.jvm - for Kotlin support")
        }
    }
    
    // Check for common missing plugins
    if pluginIDs["java"] && !pluginIDs["jacoco"] {
        suggestions = append(suggestions, "jacoco - for code coverage")
    }
    
    if pluginIDs["java"] && !pluginIDs["checkstyle"] {
        suggestions = append(suggestions, "checkstyle - for code style checking")
    }
    
    fmt.Println("💡 Plugin Suggestions")
    fmt.Println("=====================\n")
    
    if len(suggestions) > 0 {
        fmt.Println("Consider adding these plugins:")
        for _, suggestion := range suggestions {
            fmt.Printf("   • %s\n", suggestion)
        }
    } else {
        fmt.Println("✅ No additional plugins recommended")
    }
}

func main() {
    suggestPlugins("build.gradle")
}
```

## Next Steps

- [Repository Parsing Examples](./repository-parsing.md)
- [Structured Editing Examples](./structured-editing.md)
- [Custom Parser Examples](./custom-parser.md)
