# Structured Editing

Gradle Parser provides powerful structured editing capabilities that allow you to modify Gradle build files programmatically while preserving formatting and minimizing diffs.

## Quick Edits

### Update Dependency Version

The simplest way to update a dependency version:

```go
import "github.com/scagogogo/gradle-parser/pkg/api"

// Update MySQL connector version
newContent, err := api.UpdateDependencyVersion(
    "build.gradle",
    "mysql",
    "mysql-connector-java", 
    "8.0.31"
)
if err != nil {
    log.Fatal(err)
}

// Write back to file
err = os.WriteFile("build.gradle", []byte(newContent), 0644)
if err != nil {
    log.Fatal(err)
}
```

### Update Plugin Version

Update a plugin version:

```go
// Update Spring Boot plugin
newContent, err := api.UpdatePluginVersion(
    "build.gradle",
    "org.springframework.boot",
    "2.7.2"
)
if err != nil {
    log.Fatal(err)
}

// Save the changes
err = os.WriteFile("build.gradle", []byte(newContent), 0644)
if err != nil {
    log.Fatal(err)
}
```

## Advanced Editing with GradleEditor

For more complex modifications, use the `GradleEditor`:

### Creating an Editor

```go
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}
```

### Multiple Modifications

Perform multiple edits in a single operation:

```go
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Update multiple dependencies
err = editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
if err != nil {
    log.Printf("Failed to update MySQL: %v", err)
}

err = editor.UpdateDependencyVersion("org.springframework", "spring-core", "5.3.21")
if err != nil {
    log.Printf("Failed to update Spring: %v", err)
}

// Update plugin version
err = editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
if err != nil {
    log.Printf("Failed to update Spring Boot plugin: %v", err)
}

// Update project properties
err = editor.UpdateProperty("version", "1.0.0")
if err != nil {
    log.Printf("Failed to update version: %v", err)
}

// Apply all changes
newContent, err := editor.ApplyModifications()
if err != nil {
    log.Fatal(err)
}

// Save to file
err = os.WriteFile("build.gradle", []byte(newContent), 0644)
if err != nil {
    log.Fatal(err)
}
```

### Adding Dependencies

Add new dependencies to the project:

```go
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Add a new implementation dependency
err = editor.AddDependency(
    "org.apache.commons",
    "commons-text",
    "1.9",
    "implementation"
)
if err != nil {
    log.Fatal(err)
}

// Add a test dependency
err = editor.AddDependency(
    "org.junit.jupiter",
    "junit-jupiter-api",
    "5.8.2",
    "testImplementation"
)
if err != nil {
    log.Fatal(err)
}

// Apply changes
newContent, err := editor.ApplyModifications()
if err != nil {
    log.Fatal(err)
}
```

## Batch Operations

### Bulk Dependency Updates

Update multiple dependencies efficiently:

```go
func bulkUpdateDependencies(buildFile string, updates map[string]string) error {
    editor, err := api.CreateGradleEditor(buildFile)
    if err != nil {
        return err
    }

    // Apply all updates
    for depKey, newVersion := range updates {
        parts := strings.Split(depKey, ":")
        if len(parts) != 2 {
            log.Printf("Invalid dependency key: %s", depKey)
            continue
        }
        
        group, name := parts[0], parts[1]
        err = editor.UpdateDependencyVersion(group, name, newVersion)
        if err != nil {
            log.Printf("Failed to update %s: %v", depKey, err)
        }
    }

    // Apply all modifications
    newContent, err := editor.ApplyModifications()
    if err != nil {
        return err
    }

    // Write back to file
    return os.WriteFile(buildFile, []byte(newContent), 0644)
}

// Usage
updates := map[string]string{
    "mysql:mysql-connector-java":           "8.0.31",
    "org.springframework:spring-core":      "5.3.21",
    "com.google.guava:guava":              "31.1-jre",
    "org.apache.commons:commons-lang3":    "3.12.0",
}

err := bulkUpdateDependencies("build.gradle", updates)
if err != nil {
    log.Fatal(err)
}
```

### Plugin Management

```go
func updateAllPlugins(buildFile string, pluginVersions map[string]string) error {
    editor, err := api.CreateGradleEditor(buildFile)
    if err != nil {
        return err
    }

    for pluginId, version := range pluginVersions {
        err = editor.UpdatePluginVersion(pluginId, version)
        if err != nil {
            log.Printf("Failed to update plugin %s: %v", pluginId, err)
        }
    }

    newContent, err := editor.ApplyModifications()
    if err != nil {
        return err
    }

    return os.WriteFile(buildFile, []byte(newContent), 0644)
}

// Usage
pluginVersions := map[string]string{
    "org.springframework.boot":           "2.7.2",
    "io.spring.dependency-management":    "1.0.12.RELEASE",
    "org.jetbrains.kotlin.jvm":          "1.7.10",
}

err := updateAllPlugins("build.gradle", pluginVersions)
if err != nil {
    log.Fatal(err)
}
```

## Modification Tracking

### Viewing Modifications

Track what changes will be made before applying them:

```go
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Make some changes
editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
editor.UpdateProperty("version", "1.0.0")

// Review modifications before applying
modifications := editor.GetModifications()
fmt.Printf("Planned modifications: %d\n", len(modifications))

for i, mod := range modifications {
    fmt.Printf("Modification %d:\n", i+1)
    fmt.Printf("  Type: %s\n", mod.Type)
    fmt.Printf("  Description: %s\n", mod.Description)
    fmt.Printf("  Old text: %s\n", mod.OldText)
    fmt.Printf("  New text: %s\n", mod.NewText)
    fmt.Printf("  Location: Line %d, Column %d\n", 
        mod.SourceRange.Start.Line, mod.SourceRange.Start.Column)
}
```

### Conditional Modifications

Apply modifications based on conditions:

```go
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Get current project info
project := editor.GetSourceMappedProject()

// Conditional updates based on current state
for _, dep := range project.SourceMappedDependencies {
    dependency := dep.Dependency
    
    // Update Spring dependencies to latest version
    if strings.Contains(dependency.Group, "springframework") {
        err = editor.UpdateDependencyVersion(
            dependency.Group, 
            dependency.Name, 
            "5.3.21"
        )
        if err != nil {
            log.Printf("Failed to update %s: %v", dependency.Name, err)
        }
    }
    
    // Update test dependencies
    if dependency.Scope == "testImplementation" && 
       strings.Contains(dependency.Group, "junit") {
        err = editor.UpdateDependencyVersion(
            dependency.Group,
            dependency.Name,
            "5.8.2"
        )
        if err != nil {
            log.Printf("Failed to update %s: %v", dependency.Name, err)
        }
    }
}
```

## Error Handling

### Graceful Error Handling

Handle various error scenarios gracefully:

```go
func safeUpdateDependency(buildFile, group, name, version string) error {
    editor, err := api.CreateGradleEditor(buildFile)
    if err != nil {
        return fmt.Errorf("failed to create editor: %v", err)
    }

    err = editor.UpdateDependencyVersion(group, name, version)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            log.Printf("Dependency %s:%s not found, skipping", group, name)
            return nil // Not an error, just skip
        }
        return fmt.Errorf("failed to update dependency: %v", err)
    }

    newContent, err := editor.ApplyModifications()
    if err != nil {
        return fmt.Errorf("failed to apply modifications: %v", err)
    }

    err = os.WriteFile(buildFile, []byte(newContent), 0644)
    if err != nil {
        return fmt.Errorf("failed to write file: %v", err)
    }

    log.Printf("Successfully updated %s:%s to %s", group, name, version)
    return nil
}
```

### Validation Before Modification

```go
func validateAndUpdate(buildFile string) error {
    // First, parse and validate the file
    result, err := api.ParseFile(buildFile)
    if err != nil {
        return fmt.Errorf("file validation failed: %v", err)
    }

    if len(result.Warnings) > 0 {
        log.Printf("File has %d warnings, proceeding with caution", len(result.Warnings))
    }

    // Check if required dependencies exist
    requiredDeps := map[string]bool{
        "mysql:mysql-connector-java": false,
        "org.springframework:spring-core": false,
    }

    for _, dep := range result.Project.Dependencies {
        key := fmt.Sprintf("%s:%s", dep.Group, dep.Name)
        if _, exists := requiredDeps[key]; exists {
            requiredDeps[key] = true
        }
    }

    // Only proceed if all required dependencies are found
    for depKey, found := range requiredDeps {
        if !found {
            return fmt.Errorf("required dependency not found: %s", depKey)
        }
    }

    // Proceed with modifications
    editor, err := api.CreateGradleEditor(buildFile)
    if err != nil {
        return err
    }

    // Apply updates...
    return nil
}
```

## Minimal Diff Guarantee

The structured editor is designed to make minimal changes:

### Before and After Comparison

```gradle
// Original build.gradle
plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
}

group = 'com.example'
version = '0.1.0-SNAPSHOT'

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'mysql:mysql-connector-java:8.0.28'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}
```

After updating MySQL version to 8.0.31:

```gradle
// Modified build.gradle
plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
}

group = 'com.example'
version = '0.1.0-SNAPSHOT'

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'mysql:mysql-connector-java:8.0.31'  // Only this line changed
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}
```

### Preserving Formatting

The editor preserves:
- Original indentation and spacing
- Comments and their positions
- Line breaks and empty lines
- Quote styles (single vs double quotes)
- Formatting within blocks

## Best Practices

### 1. Batch Related Changes

```go
// Good: Batch related changes together
editor, _ := api.CreateGradleEditor("build.gradle")
editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
editor.UpdateDependencyVersion("org.springframework", "spring-core", "5.3.21")
editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
newContent, _ := editor.ApplyModifications()

// Avoid: Multiple separate edit operations
api.UpdateDependencyVersion("build.gradle", "mysql", "mysql-connector-java", "8.0.31")
api.UpdateDependencyVersion("build.gradle", "org.springframework", "spring-core", "5.3.21")
api.UpdatePluginVersion("build.gradle", "org.springframework.boot", "2.7.2")
```

### 2. Validate Before Editing

```go
// Always validate the file can be parsed before editing
result, err := api.ParseFile("build.gradle")
if err != nil {
    return fmt.Errorf("cannot edit invalid file: %v", err)
}

// Proceed with editing
editor, err := api.CreateGradleEditor("build.gradle")
// ...
```

### 3. Handle Missing Dependencies Gracefully

```go
err = editor.UpdateDependencyVersion("group", "name", "version")
if err != nil && strings.Contains(err.Error(), "not found") {
    // Dependency doesn't exist, add it instead
    err = editor.AddDependency("group", "name", "version", "implementation")
}
```

## Next Steps

- [Configuration](./configuration.md) - Customize parser and editor behavior
- [API Reference](../api/editor.md) - Complete editor API documentation
- [Examples](../examples/structured-editing.md) - More editing examples
