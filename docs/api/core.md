# Core API

The core API provides the main entry points for parsing Gradle build files. These functions are designed for ease of use and cover the most common parsing scenarios.

## Package Import

```go
import "github.com/scagogogo/gradle-parser/pkg/api"
```

## Basic Parsing Functions

### ParseFile

Parses a Gradle file from the filesystem.

```go
func ParseFile(filePath string) (*model.ParseResult, error)
```

**Parameters:**
- `filePath` (string): Path to the Gradle file to parse

**Returns:**
- `*model.ParseResult`: Parsed project information and metadata
- `error`: Error if parsing fails

**Example:**
```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project
fmt.Printf("Project: %s v%s\n", project.Name, project.Version)
```

### ParseString

Parses Gradle content from a string.

```go
func ParseString(content string) (*model.ParseResult, error)
```

**Parameters:**
- `content` (string): Gradle file content as string

**Returns:**
- `*model.ParseResult`: Parsed project information and metadata
- `error`: Error if parsing fails

**Example:**
```go
gradleContent := `
plugins {
    id 'java'
}
group = 'com.example'
version = '1.0.0'
`

result, err := api.ParseString(gradleContent)
if err != nil {
    log.Fatal(err)
}
```

### ParseReader

Parses Gradle content from an io.Reader.

```go
func ParseReader(reader io.Reader) (*model.ParseResult, error)
```

**Parameters:**
- `reader` (io.Reader): Reader containing Gradle file content

**Returns:**
- `*model.ParseResult`: Parsed project information and metadata
- `error`: Error if parsing fails

**Example:**
```go
file, err := os.Open("build.gradle")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

result, err := api.ParseReader(file)
if err != nil {
    log.Fatal(err)
}
```

## Component Extraction Functions

### GetDependencies

Extracts dependencies from a Gradle file.

```go
func GetDependencies(filePath string) ([]*model.Dependency, error)
```

**Parameters:**
- `filePath` (string): Path to the Gradle file

**Returns:**
- `[]*model.Dependency`: List of dependencies found in the file
- `error`: Error if extraction fails

**Example:**
```go
deps, err := api.GetDependencies("build.gradle")
if err != nil {
    log.Fatal(err)
}

for _, dep := range deps {
    fmt.Printf("%s:%s:%s (%s)\n", dep.Group, dep.Name, dep.Version, dep.Scope)
}
```

### GetPlugins

Extracts plugins from a Gradle file.

```go
func GetPlugins(filePath string) ([]*model.Plugin, error)
```

**Parameters:**
- `filePath` (string): Path to the Gradle file

**Returns:**
- `[]*model.Plugin`: List of plugins found in the file
- `error`: Error if extraction fails

**Example:**
```go
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

for _, plugin := range plugins {
    fmt.Printf("Plugin: %s", plugin.ID)
    if plugin.Version != "" {
        fmt.Printf(" v%s", plugin.Version)
    }
    fmt.Println()
}
```

### GetRepositories

Extracts repositories from a Gradle file.

```go
func GetRepositories(filePath string) ([]*model.Repository, error)
```

**Parameters:**
- `filePath` (string): Path to the Gradle file

**Returns:**
- `[]*model.Repository`: List of repositories found in the file
- `error`: Error if extraction fails

**Example:**
```go
repos, err := api.GetRepositories("build.gradle")
if err != nil {
    log.Fatal(err)
}

for _, repo := range repos {
    fmt.Printf("Repository: %s (%s)\n", repo.Name, repo.URL)
}
```

## Utility Functions

### DependenciesByScope

Groups dependencies by their scope (implementation, testImplementation, etc.).

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

### IsKotlinProject

Checks if the project uses Kotlin based on plugins.

```go
func IsKotlinProject(plugins []*model.Plugin) bool
```

**Parameters:**
- `plugins` ([]*model.Plugin): List of plugins to check

**Returns:**
- `bool`: True if Kotlin project, false otherwise

### IsSpringBootProject

Checks if the project is a Spring Boot project based on plugins.

```go
func IsSpringBootProject(plugins []*model.Plugin) bool
```

**Parameters:**
- `plugins` ([]*model.Plugin): List of plugins to check

**Returns:**
- `bool`: True if Spring Boot project, false otherwise

**Example:**
```go
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

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

## Version Information

### Version

Current version of the Gradle Parser library.

```go
const Version = "0.1.0"
```

**Example:**
```go
fmt.Printf("Using Gradle Parser v%s\n", api.Version)
```

## Error Handling

All parsing functions return detailed error information. Common error types include:

- **File not found**: When the specified file doesn't exist
- **Permission denied**: When the file cannot be read
- **Parse errors**: When the Gradle syntax is invalid

**Best Practices:**
1. Always check for errors before using results
2. Log or handle errors appropriately for your use case
3. Check the `Warnings` field in `ParseResult` for non-fatal issues

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Printf("Failed to parse file: %v", err)
    return
}

if len(result.Warnings) > 0 {
    log.Printf("Parsing completed with %d warnings", len(result.Warnings))
    for _, warning := range result.Warnings {
        log.Printf("Warning: %s", warning)
    }
}
```
