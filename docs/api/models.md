# Data Models

This section documents the data structures used to represent parsed Gradle projects. These models provide a structured way to access project information, dependencies, plugins, and other build configuration.

## Core Models

### Project

Represents a complete Gradle project with all its components.

```go
type Project struct {
    // Basic project information
    Group       string `json:"group"`
    Name        string `json:"name"`
    Version     string `json:"version"`
    Description string `json:"description"`

    // Java/JVM configuration
    SourceCompatibility string `json:"sourceCompatibility"`
    TargetCompatibility string `json:"targetCompatibility"`
    
    // Custom properties
    Properties map[string]string `json:"properties"`

    // Project components
    Plugins      []*Plugin      `json:"plugins"`
    Dependencies []*Dependency  `json:"dependencies"`
    Repositories []*Repository  `json:"repositories"`
    SubProjects  []*Project     `json:"subProjects"`
    Tasks        []*Task        `json:"tasks"`
    Extensions   map[string]any `json:"extensions"`

    // File information
    FilePath string `json:"filePath"`
}
```

**Fields:**
- `Group`: Maven group ID (e.g., "com.example")
- `Name`: Project name
- `Version`: Project version
- `Description`: Project description
- `SourceCompatibility`: Java source compatibility version
- `TargetCompatibility`: Java target compatibility version
- `Properties`: Custom project properties
- `Plugins`: List of applied plugins
- `Dependencies`: List of project dependencies
- `Repositories`: List of configured repositories
- `SubProjects`: Sub-projects in multi-module setup
- `Tasks`: Custom tasks defined in the build
- `Extensions`: Plugin extensions and configurations
- `FilePath`: Path to the source build file

### Dependency

Represents a project dependency.

```go
type Dependency struct {
    Group      string `json:"group"`
    Name       string `json:"name"`
    Version    string `json:"version"`
    Scope      string `json:"scope"`
    Transitive bool   `json:"transitive"`
    Raw        string `json:"raw"`
}
```

**Fields:**
- `Group`: Maven group ID (e.g., "org.springframework")
- `Name`: Artifact name (e.g., "spring-core")
- `Version`: Version string (e.g., "5.3.21")
- `Scope`: Dependency scope (e.g., "implementation", "testImplementation")
- `Transitive`: Whether transitive dependencies are included
- `Raw`: Original dependency declaration from build file

**Example:**
```go
// Dependency: org.springframework:spring-core:5.3.21
dep := &Dependency{
    Group:   "org.springframework",
    Name:    "spring-core", 
    Version: "5.3.21",
    Scope:   "implementation",
    Raw:     "implementation 'org.springframework:spring-core:5.3.21'",
}
```

### Plugin

Represents a Gradle plugin configuration.

```go
type Plugin struct {
    ID      string                 `json:"id"`
    Version string                 `json:"version,omitempty"`
    Apply   bool                   `json:"apply"`
    Config  map[string]interface{} `json:"config,omitempty"`
}
```

**Fields:**
- `ID`: Plugin identifier (e.g., "java", "org.springframework.boot")
- `Version`: Plugin version (optional)
- `Apply`: Whether the plugin is applied (true by default)
- `Config`: Plugin-specific configuration

**Example:**
```go
// Plugin: id 'org.springframework.boot' version '2.7.0'
plugin := &Plugin{
    ID:      "org.springframework.boot",
    Version: "2.7.0",
    Apply:   true,
}
```

### Repository

Represents a repository configuration.

```go
type Repository struct {
    Name string `json:"name"`
    URL  string `json:"url"`
    Type string `json:"type"`
}
```

**Fields:**
- `Name`: Repository name (e.g., "mavenCentral", "google")
- `URL`: Repository URL
- `Type`: Repository type (e.g., "maven", "ivy")

**Example:**
```go
// Repository: mavenCentral()
repo := &Repository{
    Name: "mavenCentral",
    URL:  "https://repo1.maven.org/maven2/",
    Type: "maven",
}
```

### Task

Represents a Gradle task definition.

```go
type Task struct {
    Name        string                 `json:"name"`
    Type        string                 `json:"type,omitempty"`
    Description string                 `json:"description,omitempty"`
    Group       string                 `json:"group,omitempty"`
    DependsOn   []string               `json:"dependsOn,omitempty"`
    Config      map[string]interface{} `json:"config,omitempty"`
}
```

**Fields:**
- `Name`: Task name
- `Type`: Task type (e.g., "Copy", "Jar")
- `Description`: Task description
- `Group`: Task group for organization
- `DependsOn`: List of task dependencies
- `Config`: Task-specific configuration

## Result Models

### ParseResult

Contains the complete result of parsing a Gradle file.

```go
type ParseResult struct {
    Project   *Project `json:"project"`
    RawText   string   `json:"rawText,omitempty"`
    Errors    []error  `json:"errors,omitempty"`
    Warnings  []string `json:"warnings,omitempty"`
    ParseTime string   `json:"parseTime,omitempty"`
}
```

**Fields:**
- `Project`: Parsed project information
- `RawText`: Original file content (if collection enabled)
- `Errors`: Fatal parsing errors
- `Warnings`: Non-fatal parsing warnings
- `ParseTime`: Time taken to parse the file

### DependencySet

Groups dependencies by scope for easier analysis.

```go
type DependencySet struct {
    Scope        string        `json:"scope"`
    Dependencies []*Dependency `json:"dependencies"`
}
```

**Fields:**
- `Scope`: Dependency scope (e.g., "implementation", "testImplementation")
- `Dependencies`: List of dependencies in this scope

## Source Mapping Models

For advanced use cases requiring precise source location tracking:

### SourceMappedProject

Extended project model with source location information.

```go
type SourceMappedProject struct {
    *Project
    
    // Source mapping information
    SourceMappedDependencies []*SourceMappedDependency
    SourceMappedPlugins      []*SourceMappedPlugin
    SourceMappedProperties   []*SourceMappedProperty
    SourceMappedRepositories []*SourceMappedRepository
    
    // Original file information
    FilePath     string
    OriginalText string
    Lines        []string
}
```

### SourcePosition

Represents a position in the source file.

```go
type SourcePosition struct {
    Line   int `json:"line"`
    Column int `json:"column"`
    Start  int `json:"start"`
    End    int `json:"end"`
    Length int `json:"length"`
}
```

### SourceRange

Represents a range in the source file.

```go
type SourceRange struct {
    Start SourcePosition `json:"start"`
    End   SourcePosition `json:"end"`
}
```

## Usage Examples

### Accessing Project Information

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

project := result.Project
fmt.Printf("Project: %s\n", project.Name)
fmt.Printf("Group: %s\n", project.Group)
fmt.Printf("Version: %s\n", project.Version)
```

### Working with Dependencies

```go
for _, dep := range project.Dependencies {
    fmt.Printf("Dependency: %s:%s", dep.Group, dep.Name)
    if dep.Version != "" {
        fmt.Printf(":%s", dep.Version)
    }
    fmt.Printf(" (%s)\n", dep.Scope)
}
```

### Analyzing Plugins

```go
for _, plugin := range project.Plugins {
    fmt.Printf("Plugin: %s", plugin.ID)
    if plugin.Version != "" {
        fmt.Printf(" v%s", plugin.Version)
    }
    if !plugin.Apply {
        fmt.Printf(" (not applied)")
    }
    fmt.Println()
}
```

### Checking Repositories

```go
for _, repo := range project.Repositories {
    fmt.Printf("Repository: %s", repo.Name)
    if repo.URL != "" {
        fmt.Printf(" (%s)", repo.URL)
    }
    fmt.Println()
}
```

## JSON Serialization

All models support JSON serialization for easy integration with web APIs and data storage:

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Serialize to JSON
jsonData, err := json.MarshalIndent(result.Project, "", "  ")
if err != nil {
    log.Fatal(err)
}

fmt.Println(string(jsonData))
```

## Type Safety

All models use Go's type system to ensure data integrity:

- String fields for textual data
- Slices for collections
- Maps for key-value configurations
- Pointers for optional relationships
- Interfaces for extensible configurations
