# Utilities

This section documents utility functions and helper methods provided by Gradle Parser.

## Package Import

```go
import "github.com/scagogogo/gradle-parser/pkg/api"
```

## Dependency Utilities

### DependenciesByScope

Groups dependencies by their scope for easier analysis.

```go
func DependenciesByScope(dependencies []*model.Dependency) []*model.DependencySet
```

**Parameters:**
- `dependencies` ([]*model.Dependency): List of dependencies to group

**Returns:**
- `[]*model.DependencySet`: Dependencies grouped by scope

**Example:**
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

## Project Type Detection

### IsAndroidProject

Checks if the project is an Android project based on plugins.

```go
func IsAndroidProject(plugins []*model.Plugin) bool
```

**Parameters:**
- `plugins` ([]*model.Plugin): List of plugins to check

**Returns:**
- `bool`: True if Android project, false otherwise

**Detection Logic:**
- Looks for `com.android.application` plugin
- Looks for `com.android.library` plugin
- Looks for `com.android.test` plugin

### IsKotlinProject

Checks if the project uses Kotlin based on plugins.

```go
func IsKotlinProject(plugins []*model.Plugin) bool
```

**Parameters:**
- `plugins` ([]*model.Plugin): List of plugins to check

**Returns:**
- `bool`: True if Kotlin project, false otherwise

**Detection Logic:**
- Looks for `kotlin` plugin
- Looks for `org.jetbrains.kotlin.jvm` plugin
- Looks for `org.jetbrains.kotlin.android` plugin
- Looks for `kotlin-android` plugin

### IsSpringBootProject

Checks if the project is a Spring Boot project based on plugins.

```go
func IsSpringBootProject(plugins []*model.Plugin) bool
```

**Parameters:**
- `plugins` ([]*model.Plugin): List of plugins to check

**Returns:**
- `bool`: True if Spring Boot project, false otherwise

**Detection Logic:**
- Looks for `org.springframework.boot` plugin
- Looks for `io.spring.dependency-management` plugin

**Example:**
```go
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

projectTypes := []string{}

if api.IsAndroidProject(plugins) {
    projectTypes = append(projectTypes, "Android")
}

if api.IsKotlinProject(plugins) {
    projectTypes = append(projectTypes, "Kotlin")
}

if api.IsSpringBootProject(plugins) {
    projectTypes = append(projectTypes, "Spring Boot")
}

if len(projectTypes) > 0 {
    fmt.Printf("Project types: %s\n", strings.Join(projectTypes, ", "))
} else {
    fmt.Println("Standard Java project")
}
```

## Configuration Utilities

### DefaultOptions

Returns default parser options.

```go
func DefaultOptions() *Options
```

**Returns:**
- `*Options`: Default configuration options

**Default Values:**
```go
&Options{
    SkipComments:      true,
    CollectRawContent: true,
    ParsePlugins:      true,
    ParseDependencies: true,
    ParseRepositories: true,
    ParseTasks:        true,
}
```

### NewParser

Creates a new parser with custom options.

```go
func NewParser(options *Options) Parser
```

**Parameters:**
- `options` (*Options): Parser configuration options

**Returns:**
- `Parser`: Configured parser instance

**Example:**
```go
// Create custom options
options := &api.Options{
    SkipComments:      false,  // Process comments
    CollectRawContent: true,   // Store original content
    ParsePlugins:      true,   // Parse plugins
    ParseDependencies: true,   // Parse dependencies
    ParseRepositories: false,  // Skip repositories
    ParseTasks:        false,  // Skip tasks
}

// Create parser with custom options
parser := api.NewParser(options)
result, err := parser.Parse(gradleContent)
```

## File Utilities

### Version

Current version of the Gradle Parser library.

```go
const Version = "0.1.0"
```

**Example:**
```go
fmt.Printf("Using Gradle Parser v%s\n", api.Version)
```

## Helper Functions

### Validation Helpers

```go
// Check if a dependency exists in a list
func HasDependency(dependencies []*model.Dependency, group, name string) bool {
    for _, dep := range dependencies {
        if dep.Group == group && dep.Name == name {
            return true
        }
    }
    return false
}

// Find a dependency by group and name
func FindDependency(dependencies []*model.Dependency, group, name string) *model.Dependency {
    for _, dep := range dependencies {
        if dep.Group == group && dep.Name == name {
            return dep
        }
    }
    return nil
}

// Check if a plugin exists in a list
func HasPlugin(plugins []*model.Plugin, id string) bool {
    for _, plugin := range plugins {
        if plugin.ID == id {
            return true
        }
    }
    return false
}

// Find a plugin by ID
func FindPlugin(plugins []*model.Plugin, id string) *model.Plugin {
    for _, plugin := range plugins {
        if plugin.ID == id {
            return plugin
        }
    }
    return nil
}
```

## Usage Examples

### Complete Project Analysis

```go
func analyzeProject(buildFile string) error {
    result, err := api.ParseFile(buildFile)
    if err != nil {
        return err
    }

    project := result.Project
    
    // Basic info
    fmt.Printf("Project: %s v%s\n", project.Name, project.Version)
    fmt.Printf("Group: %s\n", project.Group)
    
    // Project types
    plugins := project.Plugins
    types := []string{}
    
    if api.IsAndroidProject(plugins) {
        types = append(types, "Android")
    }
    if api.IsKotlinProject(plugins) {
        types = append(types, "Kotlin")
    }
    if api.IsSpringBootProject(plugins) {
        types = append(types, "Spring Boot")
    }
    
    if len(types) > 0 {
        fmt.Printf("Types: %s\n", strings.Join(types, ", "))
    }
    
    // Dependencies by scope
    depSets := api.DependenciesByScope(project.Dependencies)
    for _, depSet := range depSets {
        fmt.Printf("\n%s (%d):\n", depSet.Scope, len(depSet.Dependencies))
        for _, dep := range depSet.Dependencies {
            fmt.Printf("  %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
        }
    }
    
    return nil
}
```

### Dependency Validation

```go
func validateDependencies(buildFile string) error {
    deps, err := api.GetDependencies(buildFile)
    if err != nil {
        return err
    }
    
    // Check for required dependencies
    required := map[string]string{
        "org.springframework.boot:spring-boot-starter-web": "Spring Boot Web",
        "mysql:mysql-connector-java": "MySQL Driver",
    }
    
    missing := []string{}
    for depKey, description := range required {
        parts := strings.Split(depKey, ":")
        if len(parts) == 2 {
            if !HasDependency(deps, parts[0], parts[1]) {
                missing = append(missing, description)
            }
        }
    }
    
    if len(missing) > 0 {
        fmt.Printf("Missing required dependencies: %s\n", strings.Join(missing, ", "))
    } else {
        fmt.Println("All required dependencies found")
    }
    
    return nil
}
```

## Best Practices

1. **Use appropriate detection functions**: Choose the right project type detection function for your needs
2. **Group dependencies by scope**: Use `DependenciesByScope` for better organization
3. **Validate project structure**: Check for required dependencies and plugins
4. **Handle edge cases**: Always check for nil values and empty collections
5. **Combine utilities**: Use multiple utility functions together for comprehensive analysis
