# Parser API

The Parser API provides low-level parsing interfaces and implementations for advanced use cases. While most users should use the [Core API](./core.md), the Parser API offers fine-grained control over parsing behavior.

## Package Import

```go
import "github.com/scagogogo/gradle-parser/pkg/parser"
```

## Parser Interface

### Parser

The main parser interface that all implementations must satisfy.

```go
type Parser interface {
    Parse(content string) (*model.ParseResult, error)
    ParseFile(filePath string) (*model.ParseResult, error)
}
```

## GradleParser

The default implementation of the Parser interface.

### Constructor

```go
func NewParser() *GradleParser
```

**Returns:**
- `*GradleParser`: New parser instance with default configuration

**Example:**
```go
parser := parser.NewParser()
result, err := parser.Parse(gradleContent)
if err != nil {
    log.Fatal(err)
}
```

### Configuration Methods

The GradleParser supports method chaining for configuration:

#### WithSkipComments

Controls whether comments are processed during parsing.

```go
func (p *GradleParser) WithSkipComments(skip bool) *GradleParser
```

**Parameters:**
- `skip` (bool): If true, comments are ignored during parsing

**Example:**
```go
parser := parser.NewParser().WithSkipComments(false)
```

#### WithCollectRawContent

Controls whether the original file content is preserved.

```go
func (p *GradleParser) WithCollectRawContent(collect bool) *GradleParser
```

**Parameters:**
- `collect` (bool): If true, original content is stored in ParseResult.RawText

#### WithParsePlugins

Controls whether plugins are parsed.

```go
func (p *GradleParser) WithParsePlugins(parse bool) *GradleParser
```

**Parameters:**
- `parse` (bool): If true, plugins block is parsed

#### WithParseDependencies

Controls whether dependencies are parsed.

```go
func (p *GradleParser) WithParseDependencies(parse bool) *GradleParser
```

**Parameters:**
- `parse` (bool): If true, dependencies block is parsed

#### WithParseRepositories

Controls whether repositories are parsed.

```go
func (p *GradleParser) WithParseRepositories(parse bool) *GradleParser
```

**Parameters:**
- `parse` (bool): If true, repositories block is parsed

#### WithParseTasks

Controls whether tasks are parsed.

```go
func (p *GradleParser) WithParseTasks(parse bool) *GradleParser
```

**Parameters:**
- `parse` (bool): If true, task definitions are parsed

### Configuration Example

```go
parser := parser.NewParser().
    WithSkipComments(false).
    WithCollectRawContent(true).
    WithParsePlugins(true).
    WithParseDependencies(true).
    WithParseRepositories(true).
    WithParseTasks(false)

result, err := parser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}
```

## Source-Aware Parser

For applications requiring precise source location tracking:

### SourceAwareParser

Extended parser that tracks source locations for all parsed elements.

```go
type SourceAwareParser struct {
    *GradleParser
}
```

### Constructor

```go
func NewSourceAwareParser() *SourceAwareParser
```

**Returns:**
- `*SourceAwareParser`: New source-aware parser instance

### ParseWithSourceMapping

Parses content and returns source mapping information.

```go
func (sap *SourceAwareParser) ParseWithSourceMapping(content string) (*model.SourceMappedParseResult, error)
```

**Parameters:**
- `content` (string): Gradle file content to parse

**Returns:**
- `*model.SourceMappedParseResult`: Parse result with source location information
- `error`: Error if parsing fails

**Example:**
```go
sourceParser := parser.NewSourceAwareParser()
result, err := sourceParser.ParseWithSourceMapping(gradleContent)
if err != nil {
    log.Fatal(err)
}

// Access source-mapped project
sourceMappedProject := result.SourceMappedProject

// Each dependency has source location information
for _, dep := range sourceMappedProject.SourceMappedDependencies {
    fmt.Printf("Dependency %s:%s at line %d\n", 
        dep.Dependency.Group, 
        dep.Dependency.Name, 
        dep.SourceRange.Start.Line)
}
```

## Parser Options

### Options

Configuration structure for parser behavior.

```go
type Options struct {
    SkipComments      bool `json:"skipComments"`
    CollectRawContent bool `json:"collectRawContent"`
    ParsePlugins      bool `json:"parsePlugins"`
    ParseDependencies bool `json:"parseDependencies"`
    ParseRepositories bool `json:"parseRepositories"`
    ParseTasks        bool `json:"parseTasks"`
}
```

### DefaultOptions

Returns default parser options.

```go
func DefaultOptions() *Options
```

**Returns:**
- `*Options`: Default configuration options

**Example:**
```go
options := parser.DefaultOptions()
options.SkipComments = false
options.ParseTasks = false

// Use with API
parser := api.NewParser(options)
```

## Advanced Usage

### Custom Parser Configuration

```go
// Create parser with specific configuration
parser := parser.NewParser().
    WithSkipComments(true).           // Skip comments for faster parsing
    WithCollectRawContent(false).     // Don't store raw content to save memory
    WithParseTasks(false)             // Skip task parsing if not needed

result, err := parser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}
```

### Performance Optimization

For large files or batch processing:

```go
// Minimal parsing for dependency extraction only
parser := parser.NewParser().
    WithParsePlugins(false).
    WithParseRepositories(false).
    WithParseTasks(false).
    WithSkipComments(true).
    WithCollectRawContent(false)

result, err := parser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Only dependencies will be parsed
dependencies := result.Project.Dependencies
```

### Source Location Tracking

```go
sourceParser := parser.NewSourceAwareParser()
result, err := sourceParser.ParseWithSourceMapping(content)
if err != nil {
    log.Fatal(err)
}

// Access precise source locations
for _, dep := range result.SourceMappedProject.SourceMappedDependencies {
    sourceRange := dep.SourceRange
    fmt.Printf("Dependency at line %d, column %d to line %d, column %d\n",
        sourceRange.Start.Line, sourceRange.Start.Column,
        sourceRange.End.Line, sourceRange.End.Column)
}
```

## Error Handling

Parser-specific error handling:

```go
result, err := parser.ParseFile("build.gradle")
if err != nil {
    // Handle different error types
    if os.IsNotExist(err) {
        log.Printf("File not found: %v", err)
    } else if strings.Contains(err.Error(), "permission") {
        log.Printf("Permission denied: %v", err)
    } else {
        log.Printf("Parse error: %v", err)
    }
    return
}

// Check for parsing warnings
if len(result.Warnings) > 0 {
    log.Printf("Parsing completed with warnings:")
    for _, warning := range result.Warnings {
        log.Printf("  - %s", warning)
    }
}
```

## Integration with Core API

The Parser API integrates seamlessly with the Core API:

```go
// Using parser directly
parser := parser.NewParser().WithSkipComments(false)
result, err := parser.ParseFile("build.gradle")

// Using via Core API with custom options
options := &api.Options{
    SkipComments:      false,
    CollectRawContent: true,
    ParsePlugins:      true,
    ParseDependencies: true,
    ParseRepositories: true,
    ParseTasks:        true,
}
apiParser := api.NewParser(options)
result, err = apiParser.Parse(content)
```

## Best Practices

1. **Use appropriate parsing level**: Choose between basic parsing and source-aware parsing based on your needs
2. **Configure for performance**: Disable unnecessary parsing features for better performance
3. **Handle errors gracefully**: Always check for errors and warnings
4. **Reuse parser instances**: Create parser once and reuse for multiple files
5. **Memory management**: Disable raw content collection for large-scale processing
