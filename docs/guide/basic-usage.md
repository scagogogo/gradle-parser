# Basic Usage

This guide covers the fundamental operations you can perform with Gradle Parser, including parsing files, extracting information, and working with the parsed data.

## Parsing Gradle Files

### Simple File Parsing

The most straightforward way to parse a Gradle file:

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    result, err := api.ParseFile("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    project := result.Project
    fmt.Printf("Project: %s\n", project.Name)
    fmt.Printf("Group: %s\n", project.Group)
    fmt.Printf("Version: %s\n", project.Version)
}
```

### Parsing Different Sources

You can parse Gradle content from various sources:

```go
// From file
result, err := api.ParseFile("build.gradle")

// From string
gradleContent := `
plugins {
    id 'java'
}
group = 'com.example'
version = '1.0.0'
`
result, err = api.ParseString(gradleContent)

// From any io.Reader
file, err := os.Open("build.gradle")
if err != nil {
    log.Fatal(err)
}
defer file.Close()
result, err = api.ParseReader(file)
```

## Working with Dependencies

### Extracting Dependencies

Get all dependencies from a Gradle file:

```go
deps, err := api.GetDependencies("build.gradle")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d dependencies:\n", len(deps))
for _, dep := range deps {
    fmt.Printf("  %s:%s:%s (%s)\n", 
        dep.Group, dep.Name, dep.Version, dep.Scope)
}
```

### Grouping Dependencies by Scope

Organize dependencies by their scope (implementation, testImplementation, etc.):

```go
deps, err := api.GetDependencies("build.gradle")
if err != nil {
    log.Fatal(err)
}

depSets := api.DependenciesByScope(deps)
for _, depSet := range depSets {
    fmt.Printf("\n%s dependencies:\n", depSet.Scope)
    for _, dep := range depSet.Dependencies {
        fmt.Printf("  %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
    }
}
```

### Analyzing Dependencies

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

dependencies := result.Project.Dependencies

// Count dependencies by scope
scopeCount := make(map[string]int)
for _, dep := range dependencies {
    scopeCount[dep.Scope]++
}

fmt.Println("Dependencies by scope:")
for scope, count := range scopeCount {
    fmt.Printf("  %s: %d\n", scope, count)
}

// Find specific dependencies
for _, dep := range dependencies {
    if dep.Group == "org.springframework" {
        fmt.Printf("Found Spring dependency: %s:%s:%s\n", 
            dep.Group, dep.Name, dep.Version)
    }
}
```

## Working with Plugins

### Extracting Plugins

Get all plugins from a Gradle file:

```go
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d plugins:\n", len(plugins))
for _, plugin := range plugins {
    fmt.Printf("  %s", plugin.ID)
    if plugin.Version != "" {
        fmt.Printf(" (v%s)", plugin.Version)
    }
    if !plugin.Apply {
        fmt.Printf(" [not applied]")
    }
    fmt.Println()
}
```

### Project Type Detection

Determine the type of project based on applied plugins:

```go
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Check project types
if api.IsAndroidProject(plugins) {
    fmt.Println("✓ Android project")
}

if api.IsKotlinProject(plugins) {
    fmt.Println("✓ Kotlin project")
}

if api.IsSpringBootProject(plugins) {
    fmt.Println("✓ Spring Boot project")
}

// Check for specific plugins
hasJavaPlugin := false
for _, plugin := range plugins {
    if plugin.ID == "java" {
        hasJavaPlugin = true
        break
    }
}

if hasJavaPlugin {
    fmt.Println("✓ Java project")
}
```

## Working with Repositories

### Extracting Repositories

Get repository configurations:

```go
repos, err := api.GetRepositories("build.gradle")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d repositories:\n", len(repos))
for _, repo := range repos {
    fmt.Printf("  %s", repo.Name)
    if repo.URL != "" {
        fmt.Printf(" (%s)", repo.URL)
    }
    fmt.Println()
}
```

### Repository Analysis

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

repositories := result.Project.Repositories

// Check for common repositories
hasMaven := false
hasGoogle := false
hasJCenter := false

for _, repo := range repositories {
    switch repo.Name {
    case "mavenCentral":
        hasMaven = true
    case "google":
        hasGoogle = true
    case "jcenter":
        hasJCenter = true
    }
}

fmt.Println("Repository analysis:")
fmt.Printf("  Maven Central: %v\n", hasMaven)
fmt.Printf("  Google: %v\n", hasGoogle)
fmt.Printf("  JCenter: %v\n", hasJCenter)
```

## Working with Project Properties

### Accessing Basic Properties

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project

// Basic project information
fmt.Printf("Project Information:\n")
fmt.Printf("  Name: %s\n", project.Name)
fmt.Printf("  Group: %s\n", project.Group)
fmt.Printf("  Version: %s\n", project.Version)
fmt.Printf("  Description: %s\n", project.Description)

// Java compatibility
if project.SourceCompatibility != "" {
    fmt.Printf("  Source Compatibility: %s\n", project.SourceCompatibility)
}
if project.TargetCompatibility != "" {
    fmt.Printf("  Target Compatibility: %s\n", project.TargetCompatibility)
}
```

### Custom Properties

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project

// Access custom properties
if len(project.Properties) > 0 {
    fmt.Println("Custom properties:")
    for key, value := range project.Properties {
        fmt.Printf("  %s = %s\n", key, value)
    }
}
```

## Error Handling and Validation

### Comprehensive Error Handling

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    // Handle different types of errors
    if os.IsNotExist(err) {
        log.Printf("Gradle file not found: %v", err)
        return
    }
    
    if strings.Contains(err.Error(), "permission") {
        log.Printf("Permission denied: %v", err)
        return
    }
    
    log.Printf("Parse error: %v", err)
    return
}

// Check for warnings
if len(result.Warnings) > 0 {
    fmt.Printf("Parsing completed with %d warnings:\n", len(result.Warnings))
    for i, warning := range result.Warnings {
        fmt.Printf("  %d. %s\n", i+1, warning)
    }
}

// Validate results
project := result.Project
if project == nil {
    log.Println("Warning: No project information found")
    return
}

if len(project.Dependencies) == 0 {
    log.Println("Warning: No dependencies found")
}

if len(project.Plugins) == 0 {
    log.Println("Warning: No plugins found")
}
```

### Data Validation

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project

// Validate project structure
errors := []string{}

if project.Group == "" {
    errors = append(errors, "Missing project group")
}

if project.Version == "" {
    errors = append(errors, "Missing project version")
}

// Check for required plugins
hasJavaPlugin := false
for _, plugin := range project.Plugins {
    if plugin.ID == "java" || plugin.ID == "java-library" {
        hasJavaPlugin = true
        break
    }
}

if !hasJavaPlugin {
    errors = append(errors, "No Java plugin found")
}

// Report validation results
if len(errors) > 0 {
    fmt.Println("Validation errors:")
    for _, err := range errors {
        fmt.Printf("  - %s\n", err)
    }
} else {
    fmt.Println("✓ Project structure is valid")
}
```

## Performance Considerations

### Parsing Large Files

For large Gradle files, consider using custom parser options:

```go
// Create parser with optimized settings
options := api.DefaultOptions()
options.SkipComments = true        // Skip comments for faster parsing
options.CollectRawContent = false  // Don't store raw content to save memory
options.ParseTasks = false         // Skip task parsing if not needed

parser := api.NewParser(options)
result, err := parser.Parse(gradleContent)
if err != nil {
    log.Fatal(err)
}
```

### Batch Processing

When processing multiple files:

```go
// Reuse parser instance
parser := api.NewParser(api.DefaultOptions())

files := []string{"app/build.gradle", "lib/build.gradle", "core/build.gradle"}
for _, file := range files {
    result, err := parser.ParseFile(file)
    if err != nil {
        log.Printf("Failed to parse %s: %v", file, err)
        continue
    }
    
    fmt.Printf("Processed %s: %d dependencies\n", 
        file, len(result.Project.Dependencies))
}
```

## Next Steps

Now that you understand the basics, explore more advanced features:

- [Advanced Features](./advanced-features.md) - Source mapping, custom parsing, and more
- [Structured Editing](./structured-editing.md) - Modify Gradle files programmatically
- [Configuration](./configuration.md) - Customize parser behavior
- [API Reference](../api/) - Complete API documentation
