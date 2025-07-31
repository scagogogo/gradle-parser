# Editor API

The Editor API provides structured editing capabilities for Gradle build files. It allows you to modify dependencies, plugins, and properties while preserving the original formatting and minimizing diffs.

## Package Import

```go
import "github.com/scagogogo/gradle-parser/pkg/api"
```

## Quick Edit Functions

### UpdateDependencyVersion

Updates a dependency version in a Gradle file.

```go
func UpdateDependencyVersion(filePath, group, name, newVersion string) (string, error)
```

**Parameters:**
- `filePath` (string): Path to the Gradle file
- `group` (string): Dependency group ID
- `name` (string): Dependency artifact name  
- `newVersion` (string): New version to set

**Returns:**
- `string`: Updated file content
- `error`: Error if update fails

**Example:**
```go
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

### UpdatePluginVersion

Updates a plugin version in a Gradle file.

```go
func UpdatePluginVersion(filePath, pluginId, newVersion string) (string, error)
```

**Parameters:**
- `filePath` (string): Path to the Gradle file
- `pluginId` (string): Plugin identifier
- `newVersion` (string): New version to set

**Returns:**
- `string`: Updated file content
- `error`: Error if update fails

**Example:**
```go
// Update Spring Boot plugin version
newContent, err := api.UpdatePluginVersion(
    "build.gradle",
    "org.springframework.boot",
    "2.7.2"
)
if err != nil {
    log.Fatal(err)
}
```

## Advanced Editor

### CreateGradleEditor

Creates a structured editor for more complex modifications.

```go
func CreateGradleEditor(filePath string) (*editor.GradleEditor, error)
```

**Parameters:**
- `filePath` (string): Path to the Gradle file

**Returns:**
- `*editor.GradleEditor`: Editor instance for the file
- `error`: Error if creation fails

**Example:**
```go
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Perform multiple edits
err = editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
if err != nil {
    log.Fatal(err)
}

err = editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
if err != nil {
    log.Fatal(err)
}

err = editor.UpdateProperty("version", "1.0.0")
if err != nil {
    log.Fatal(err)
}
```

## GradleEditor Methods

### UpdateDependencyVersion

Updates a dependency version using the editor.

```go
func (ge *GradleEditor) UpdateDependencyVersion(group, name, newVersion string) error
```

**Parameters:**
- `group` (string): Dependency group ID
- `name` (string): Dependency artifact name
- `newVersion` (string): New version to set

**Returns:**
- `error`: Error if update fails

### UpdatePluginVersion

Updates a plugin version using the editor.

```go
func (ge *GradleEditor) UpdatePluginVersion(pluginId, newVersion string) error
```

**Parameters:**
- `pluginId` (string): Plugin identifier
- `newVersion` (string): New version to set

**Returns:**
- `error`: Error if update fails

### UpdateProperty

Updates a project property using the editor.

```go
func (ge *GradleEditor) UpdateProperty(key, newValue string) error
```

**Parameters:**
- `key` (string): Property name (e.g., "version", "group")
- `newValue` (string): New property value

**Returns:**
- `error`: Error if update fails

### AddDependency

Adds a new dependency to the project.

```go
func (ge *GradleEditor) AddDependency(group, name, version, scope string) error
```

**Parameters:**
- `group` (string): Dependency group ID
- `name` (string): Dependency artifact name
- `version` (string): Dependency version
- `scope` (string): Dependency scope (e.g., "implementation", "testImplementation")

**Returns:**
- `error`: Error if addition fails

**Example:**
```go
err = editor.AddDependency(
    "org.apache.commons",
    "commons-text", 
    "1.9",
    "implementation"
)
if err != nil {
    log.Fatal(err)
}
```

### GetModifications

Returns all modifications made by the editor.

```go
func (ge *GradleEditor) GetModifications() []Modification
```

**Returns:**
- `[]Modification`: List of all modifications

### ClearModifications

Clears all pending modifications.

```go
func (ge *GradleEditor) ClearModifications()
```

## Serialization

### ApplyModifications

Applies all modifications and returns the updated content.

```go
func (ge *GradleEditor) ApplyModifications() (string, error)
```

**Returns:**
- `string`: Updated file content with all modifications applied
- `error`: Error if serialization fails

**Example:**
```go
// Make multiple changes
editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
editor.AddDependency("org.apache.commons", "commons-text", "1.9", "implementation")

// Apply all changes
newContent, err := editor.ApplyModifications()
if err != nil {
    log.Fatal(err)
}

// Write to file
err = os.WriteFile("build.gradle", []byte(newContent), 0644)
if err != nil {
    log.Fatal(err)
}
```

## Modification Types

### Modification

Represents a single modification operation.

```go
type Modification struct {
    Type        ModificationType  `json:"type"`
    SourceRange model.SourceRange `json:"sourceRange"`
    OldText     string            `json:"oldText"`
    NewText     string            `json:"newText"`
    Description string            `json:"description"`
}
```

### ModificationType

Types of modifications supported.

```go
type ModificationType string

const (
    ModificationTypeReplace ModificationType = "replace"
    ModificationTypeInsert  ModificationType = "insert"
    ModificationTypeDelete  ModificationType = "delete"
)
```

## Advanced Features

### Source-Aware Parsing

For precise editing, use source-aware parsing:

```go
result, err := api.ParseFileWithSourceMapping("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Create editor from source-mapped result
editor := editor.NewGradleEditor(result.SourceMappedProject)
```

### Batch Operations

Perform multiple operations efficiently:

```go
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Batch updates
updates := []struct {
    group, name, version string
}{
    {"mysql", "mysql-connector-java", "8.0.31"},
    {"org.springframework", "spring-core", "5.3.21"},
    {"com.google.guava", "guava", "31.1-jre"},
}

for _, update := range updates {
    err = editor.UpdateDependencyVersion(update.group, update.name, update.version)
    if err != nil {
        log.Printf("Failed to update %s:%s: %v", update.group, update.name, err)
    }
}

// Apply all changes at once
newContent, err := editor.ApplyModifications()
if err != nil {
    log.Fatal(err)
}
```

## Error Handling

Common error scenarios:

- **Dependency not found**: When trying to update a non-existent dependency
- **Plugin not found**: When trying to update a non-existent plugin
- **Property not found**: When trying to update a non-existent property
- **Invalid syntax**: When the modification would create invalid Gradle syntax

**Best Practices:**
```go
err = editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        log.Printf("Dependency not found, skipping update")
    } else {
        log.Printf("Update failed: %v", err)
        return err
    }
}
```

## Minimal Diff Guarantee

The editor is designed to make minimal changes:

- **Preserves formatting**: Original indentation and spacing are maintained
- **Preserves comments**: Comments are not modified or removed
- **Minimal changes**: Only the specific values being updated are changed
- **Line-by-line**: Changes are made line-by-line to minimize diff size

**Example:**
```gradle
// Original
implementation 'mysql:mysql-connector-java:8.0.28'

// After UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
implementation 'mysql:mysql-connector-java:8.0.31'
```

Only the version number is changed, everything else remains identical.
