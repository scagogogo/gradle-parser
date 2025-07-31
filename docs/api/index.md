# API Reference

Welcome to the Gradle Parser API reference. This section provides comprehensive documentation for all public APIs, types, and functions.

## Overview

Gradle Parser provides a clean and intuitive API for parsing and manipulating Gradle build files. The main entry point is the `api` package, which offers high-level functions for common operations.

## Package Structure

- **[Core API](./core.md)** - Main parsing functions and utilities
- **[Parser](./parser.md)** - Low-level parsing interfaces and implementations
- **[Models](./models.md)** - Data structures representing Gradle projects
- **[Editor](./editor.md)** - Structured editing and modification capabilities
- **[Utilities](./utilities.md)** - Helper functions and utilities

## Quick Reference

### Core Functions

```go
// Basic parsing
result, err := api.ParseFile("build.gradle")
result, err := api.ParseString(content)
result, err := api.ParseReader(reader)

// Extract specific components
deps, err := api.GetDependencies("build.gradle")
plugins, err := api.GetPlugins("build.gradle")
repos, err := api.GetRepositories("build.gradle")

// Project type detection
isAndroid := api.IsAndroidProject(plugins)
isKotlin := api.IsKotlinProject(plugins)
isSpringBoot := api.IsSpringBootProject(plugins)

// Structured editing
newText, err := api.UpdateDependencyVersion("build.gradle", "group", "name", "version")
newText, err := api.UpdatePluginVersion("build.gradle", "plugin", "version")
```

### Advanced Features

```go
// Source-aware parsing
result, err := api.ParseFileWithSourceMapping("build.gradle")

// Custom parser configuration
options := api.DefaultOptions()
options.SkipComments = false
parser := api.NewParser(options)

// Structured editing with editor
editor, err := api.CreateGradleEditor("build.gradle")
err = editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
err = editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
```

## Common Patterns

### Error Handling

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Printf("Failed to parse: %v", err)
    return
}

// Check for warnings
if len(result.Warnings) > 0 {
    for _, warning := range result.Warnings {
        log.Printf("Warning: %s", warning)
    }
}
```

### Working with Dependencies

```go
// Get all dependencies
deps, err := api.GetDependencies("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Group by scope
depSets := api.DependenciesByScope(deps)
for _, depSet := range depSets {
    fmt.Printf("Scope: %s\n", depSet.Scope)
    for _, dep := range depSet.Dependencies {
        fmt.Printf("  %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
    }
}
```

### Project Analysis

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project

// Basic project info
fmt.Printf("Project: %s\n", project.Name)
fmt.Printf("Group: %s\n", project.Group)
fmt.Printf("Version: %s\n", project.Version)

// Analyze plugins
plugins := project.Plugins
if api.IsAndroidProject(plugins) {
    fmt.Println("This is an Android project")
}
if api.IsKotlinProject(plugins) {
    fmt.Println("This project uses Kotlin")
}
if api.IsSpringBootProject(plugins) {
    fmt.Println("This is a Spring Boot project")
}
```

## Type System

The API uses a rich type system to represent Gradle projects:

- **Project** - Root project information
- **Dependency** - Individual dependencies with scope and version info
- **Plugin** - Plugin configurations and versions
- **Repository** - Repository definitions and URLs
- **Task** - Task definitions and configurations

See the [Models](./models.md) section for detailed type documentation.

## Error Types

Common error scenarios:

- **File not found** - When the specified Gradle file doesn't exist
- **Parse errors** - When the Gradle file contains invalid syntax
- **Permission errors** - When the file cannot be read due to permissions

## Best Practices

1. **Always handle errors** - Check return values and handle errors appropriately
2. **Use appropriate parsing level** - Choose between basic parsing and source-aware parsing based on your needs
3. **Configure parser options** - Customize parsing behavior for better performance
4. **Validate results** - Check for warnings and validate parsed data
5. **Use structured editing** - Prefer structured editing over string manipulation for modifications

## Version Compatibility

- **Go Version**: Requires Go 1.19 or later
- **Gradle Compatibility**: Supports Gradle 4.0+ syntax
- **DSL Support**: Full Groovy DSL support, basic Kotlin DSL support
