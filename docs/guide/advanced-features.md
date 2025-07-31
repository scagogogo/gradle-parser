# Advanced Features

This guide covers advanced features of Gradle Parser, including source mapping, custom parsing configurations, and specialized use cases.

## Source-Aware Parsing

Source-aware parsing tracks the exact location of every parsed element in the original file, enabling precise modifications and detailed analysis.

### Basic Source Mapping

```go
import "github.com/scagogogo/gradle-parser/pkg/api"

// Parse with source mapping
result, err := api.ParseFileWithSourceMapping("build.gradle")
if err != nil {
    log.Fatal(err)
}

sourceMappedProject := result.SourceMappedProject

// Access dependencies with source locations
for _, dep := range sourceMappedProject.SourceMappedDependencies {
    dependency := dep.Dependency
    sourceRange := dep.SourceRange
    
    fmt.Printf("Dependency: %s:%s:%s\n", 
        dependency.Group, dependency.Name, dependency.Version)
    fmt.Printf("  Location: Line %d, Column %d to Line %d, Column %d\n",
        sourceRange.Start.Line, sourceRange.Start.Column,
        sourceRange.End.Line, sourceRange.End.Column)
    fmt.Printf("  Raw text: %s\n", dep.RawText)
}
```

### Source Location Details

```go
result, err := api.ParseFileWithSourceMapping("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Access detailed source information
sourceMappedProject := result.SourceMappedProject

fmt.Printf("File: %s\n", sourceMappedProject.FilePath)
fmt.Printf("Total lines: %d\n", len(sourceMappedProject.Lines))

// Examine plugins with source locations
for _, plugin := range sourceMappedProject.SourceMappedPlugins {
    fmt.Printf("Plugin: %s\n", plugin.Plugin.ID)
    
    pos := plugin.SourceRange.Start
    fmt.Printf("  Position: Line %d, Column %d\n", pos.Line, pos.Column)
    fmt.Printf("  Character range: %d-%d\n", pos.Start, pos.End)
    fmt.Printf("  Length: %d characters\n", pos.Length)
}
```

## Custom Parser Configuration

### Performance Optimization

Configure the parser for specific use cases:

```go
import "github.com/scagogogo/gradle-parser/pkg/parser"

// High-performance configuration for dependency extraction only
fastParser := parser.NewParser().
    WithSkipComments(true).           // Skip comment processing
    WithCollectRawContent(false).     // Don't store original content
    WithParsePlugins(false).          // Skip plugin parsing
    WithParseRepositories(false).     // Skip repository parsing
    WithParseTasks(false)             // Skip task parsing

result, err := fastParser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Only dependencies will be parsed
dependencies := result.Project.Dependencies
fmt.Printf("Found %d dependencies\n", len(dependencies))
```

### Detailed Analysis Configuration

```go
// Comprehensive parsing with all features enabled
detailedParser := parser.NewParser().
    WithSkipComments(false).          // Process comments
    WithCollectRawContent(true).      // Store original content
    WithParsePlugins(true).           // Parse plugins
    WithParseDependencies(true).      // Parse dependencies
    WithParseRepositories(true).      // Parse repositories
    WithParseTasks(true)              // Parse tasks

result, err := detailedParser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Access all parsed information
project := result.Project
fmt.Printf("Plugins: %d\n", len(project.Plugins))
fmt.Printf("Dependencies: %d\n", len(project.Dependencies))
fmt.Printf("Repositories: %d\n", len(project.Repositories))
fmt.Printf("Tasks: %d\n", len(project.Tasks))
fmt.Printf("Raw content length: %d\n", len(result.RawText))
```

## Working with Multi-Module Projects

### Parsing Multiple Build Files

```go
func parseMultiModuleProject(rootDir string) error {
    // Parse root build.gradle
    rootBuild := filepath.Join(rootDir, "build.gradle")
    rootResult, err := api.ParseFile(rootBuild)
    if err != nil {
        return fmt.Errorf("failed to parse root build: %v", err)
    }

    fmt.Printf("Root project: %s\n", rootResult.Project.Name)

    // Parse settings.gradle to find subprojects
    settingsFile := filepath.Join(rootDir, "settings.gradle")
    if _, err := os.Stat(settingsFile); err == nil {
        settingsResult, err := api.ParseFile(settingsFile)
        if err != nil {
            log.Printf("Warning: failed to parse settings.gradle: %v", err)
        } else {
            fmt.Printf("Settings parsed successfully\n")
        }
    }

    // Parse subproject build files
    subprojectDirs := []string{"app", "lib", "core", "common"}
    for _, subdir := range subprojectDirs {
        subBuildFile := filepath.Join(rootDir, subdir, "build.gradle")
        if _, err := os.Stat(subBuildFile); err == nil {
            subResult, err := api.ParseFile(subBuildFile)
            if err != nil {
                log.Printf("Warning: failed to parse %s: %v", subBuildFile, err)
                continue
            }
            
            fmt.Printf("Subproject %s: %d dependencies\n", 
                subdir, len(subResult.Project.Dependencies))
        }
    }

    return nil
}
```

### Project Dependency Analysis

```go
func analyzeProjectDependencies(projectDir string) error {
    buildFiles := []string{
        filepath.Join(projectDir, "build.gradle"),
        filepath.Join(projectDir, "app", "build.gradle"),
        filepath.Join(projectDir, "lib", "build.gradle"),
    }

    allDependencies := make(map[string]*api.Dependency)
    dependencyUsage := make(map[string][]string)

    for _, buildFile := range buildFiles {
        if _, err := os.Stat(buildFile); err != nil {
            continue // Skip non-existent files
        }

        deps, err := api.GetDependencies(buildFile)
        if err != nil {
            log.Printf("Warning: failed to parse %s: %v", buildFile, err)
            continue
        }

        moduleName := filepath.Base(filepath.Dir(buildFile))
        if moduleName == filepath.Base(projectDir) {
            moduleName = "root"
        }

        for _, dep := range deps {
            key := fmt.Sprintf("%s:%s", dep.Group, dep.Name)
            allDependencies[key] = dep
            dependencyUsage[key] = append(dependencyUsage[key], moduleName)
        }
    }

    // Report shared dependencies
    fmt.Println("Shared dependencies:")
    for depKey, modules := range dependencyUsage {
        if len(modules) > 1 {
            dep := allDependencies[depKey]
            fmt.Printf("  %s:%s (used in: %s)\n", 
                dep.Group, dep.Name, strings.Join(modules, ", "))
        }
    }

    return nil
}
```

## Advanced Dependency Analysis

### Version Conflict Detection

```go
func detectVersionConflicts(buildFile string) error {
    deps, err := api.GetDependencies(buildFile)
    if err != nil {
        return err
    }

    // Group dependencies by group:name
    depVersions := make(map[string][]string)
    for _, dep := range deps {
        key := fmt.Sprintf("%s:%s", dep.Group, dep.Name)
        depVersions[key] = append(depVersions[key], dep.Version)
    }

    // Find conflicts
    conflicts := []string{}
    for depKey, versions := range depVersions {
        uniqueVersions := make(map[string]bool)
        for _, version := range versions {
            uniqueVersions[version] = true
        }
        
        if len(uniqueVersions) > 1 {
            versionList := make([]string, 0, len(uniqueVersions))
            for version := range uniqueVersions {
                versionList = append(versionList, version)
            }
            conflicts = append(conflicts, 
                fmt.Sprintf("%s: %s", depKey, strings.Join(versionList, ", ")))
        }
    }

    if len(conflicts) > 0 {
        fmt.Println("Version conflicts detected:")
        for _, conflict := range conflicts {
            fmt.Printf("  %s\n", conflict)
        }
    } else {
        fmt.Println("No version conflicts found")
    }

    return nil
}
```

### Dependency Tree Analysis

```go
func analyzeDependencyTree(buildFile string) error {
    result, err := api.ParseFile(buildFile)
    if err != nil {
        return err
    }

    dependencies := result.Project.Dependencies
    
    // Group by scope
    depSets := api.DependenciesByScope(dependencies)
    
    for _, depSet := range depSets {
        fmt.Printf("\n%s dependencies (%d):\n", depSet.Scope, len(depSet.Dependencies))
        
        // Group by organization
        orgGroups := make(map[string][]*api.Dependency)
        for _, dep := range depSet.Dependencies {
            org := strings.Split(dep.Group, ".")[0]
            orgGroups[org] = append(orgGroups[org], dep)
        }
        
        for org, deps := range orgGroups {
            fmt.Printf("  %s:\n", org)
            for _, dep := range deps {
                fmt.Printf("    %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
            }
        }
    }

    return nil
}
```

## Plugin Analysis

### Advanced Plugin Detection

```go
func analyzePlugins(buildFile string) error {
    plugins, err := api.GetPlugins(buildFile)
    if err != nil {
        return err
    }

    // Categorize plugins
    categories := map[string][]string{
        "Language":    {"java", "kotlin", "groovy", "scala"},
        "Framework":   {"org.springframework.boot", "io.quarkus", "io.micronaut"},
        "Android":     {"com.android.application", "com.android.library"},
        "Build":       {"maven-publish", "signing", "distribution"},
        "Quality":     {"checkstyle", "pmd", "spotbugs", "jacoco"},
        "Development": {"idea", "eclipse", "application"},
    }

    pluginsByCategory := make(map[string][]*api.Plugin)
    uncategorized := []*api.Plugin{}

    for _, plugin := range plugins {
        categorized := false
        for category, pluginIds := range categories {
            for _, id := range pluginIds {
                if plugin.ID == id || strings.Contains(plugin.ID, id) {
                    pluginsByCategory[category] = append(pluginsByCategory[category], plugin)
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

    // Report categorized plugins
    for category, categoryPlugins := range pluginsByCategory {
        if len(categoryPlugins) > 0 {
            fmt.Printf("%s plugins:\n", category)
            for _, plugin := range categoryPlugins {
                fmt.Printf("  %s", plugin.ID)
                if plugin.Version != "" {
                    fmt.Printf(" v%s", plugin.Version)
                }
                fmt.Println()
            }
        }
    }

    if len(uncategorized) > 0 {
        fmt.Println("Other plugins:")
        for _, plugin := range uncategorized {
            fmt.Printf("  %s", plugin.ID)
            if plugin.Version != "" {
                fmt.Printf(" v%s", plugin.Version)
            }
            fmt.Println()
        }
    }

    return nil
}
```

## Custom Data Extraction

### Extract Custom Properties

```go
func extractCustomProperties(buildFile string) error {
    result, err := api.ParseFile(buildFile)
    if err != nil {
        return err
    }

    project := result.Project

    // Standard properties
    fmt.Println("Standard properties:")
    if project.Group != "" {
        fmt.Printf("  group: %s\n", project.Group)
    }
    if project.Version != "" {
        fmt.Printf("  version: %s\n", project.Version)
    }
    if project.Description != "" {
        fmt.Printf("  description: %s\n", project.Description)
    }

    // Java properties
    if project.SourceCompatibility != "" {
        fmt.Printf("  sourceCompatibility: %s\n", project.SourceCompatibility)
    }
    if project.TargetCompatibility != "" {
        fmt.Printf("  targetCompatibility: %s\n", project.TargetCompatibility)
    }

    // Custom properties
    if len(project.Properties) > 0 {
        fmt.Println("Custom properties:")
        for key, value := range project.Properties {
            fmt.Printf("  %s: %s\n", key, value)
        }
    }

    return nil
}
```

## Performance Monitoring

### Parse Time Analysis

```go
func benchmarkParsing(buildFile string) error {
    // Measure different parsing configurations
    configs := map[string]*parser.GradleParser{
        "minimal": parser.NewParser().
            WithSkipComments(true).
            WithCollectRawContent(false).
            WithParseTasks(false),
        "standard": parser.NewParser(),
        "comprehensive": parser.NewParser().
            WithSkipComments(false).
            WithCollectRawContent(true).
            WithParseTasks(true),
    }

    for name, p := range configs {
        start := time.Now()
        result, err := p.ParseFile(buildFile)
        duration := time.Since(start)
        
        if err != nil {
            log.Printf("Config %s failed: %v", name, err)
            continue
        }

        fmt.Printf("Config %s:\n", name)
        fmt.Printf("  Parse time: %v\n", duration)
        fmt.Printf("  Dependencies: %d\n", len(result.Project.Dependencies))
        fmt.Printf("  Plugins: %d\n", len(result.Project.Plugins))
        fmt.Printf("  Raw content: %d bytes\n", len(result.RawText))
        fmt.Println()
    }

    return nil
}
```

## Next Steps

Explore more advanced topics:

- [Structured Editing](./structured-editing.md) - Modify Gradle files programmatically
- [Configuration](./configuration.md) - Detailed parser configuration options
- [API Reference](../api/) - Complete API documentation
- [Examples](../examples/) - Practical usage examples
