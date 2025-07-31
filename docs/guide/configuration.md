# Configuration

This guide covers how to configure Gradle Parser for different use cases, including performance optimization, feature selection, and custom parsing behavior.

## Parser Options

### Default Configuration

The default parser configuration is optimized for general use:

```go
import "github.com/scagogogo/gradle-parser/pkg/api"

// Using default options
result, err := api.ParseFile("build.gradle")

// Or explicitly with default options
options := api.DefaultOptions()
parser := api.NewParser(options)
result, err = parser.Parse(content)
```

### Available Options

```go
type Options struct {
    SkipComments      bool `json:"skipComments"`      // Skip comment processing
    CollectRawContent bool `json:"collectRawContent"` // Store original file content
    ParsePlugins      bool `json:"parsePlugins"`      // Parse plugins block
    ParseDependencies bool `json:"parseDependencies"` // Parse dependencies block
    ParseRepositories bool `json:"parseRepositories"` // Parse repositories block
    ParseTasks        bool `json:"parseTasks"`        // Parse task definitions
}
```

### Custom Configuration

```go
import "github.com/scagogogo/gradle-parser/pkg/parser"

// Create custom configuration
customParser := parser.NewParser().
    WithSkipComments(true).           // Skip comments for faster parsing
    WithCollectRawContent(false).     // Don't store raw content to save memory
    WithParsePlugins(true).           // Parse plugins (enabled)
    WithParseDependencies(true).      // Parse dependencies (enabled)
    WithParseRepositories(false).     // Skip repositories
    WithParseTasks(false)             // Skip tasks

result, err := customParser.ParseFile("build.gradle")
```

## Performance Optimization

### High-Performance Configuration

For scenarios requiring maximum parsing speed:

```go
// Minimal parsing configuration
fastParser := parser.NewParser().
    WithSkipComments(true).           // Skip comment processing
    WithCollectRawContent(false).     // Don't store original content
    WithParsePlugins(false).          // Skip plugin parsing
    WithParseRepositories(false).     // Skip repository parsing
    WithParseTasks(false)             // Skip task parsing

// Only dependencies will be parsed
result, err := fastParser.ParseFile("build.gradle")
dependencies := result.Project.Dependencies
```

### Memory-Optimized Configuration

For processing large files or many files:

```go
// Memory-efficient configuration
memoryOptimizedParser := parser.NewParser().
    WithCollectRawContent(false).     // Don't store raw content
    WithSkipComments(true).           // Skip comments
    WithParseTasks(false)             // Skip tasks if not needed

result, err := memoryOptimizedParser.ParseFile("build.gradle")
```

### Batch Processing Configuration

For processing multiple files efficiently:

```go
func processBuildFiles(files []string) error {
    // Create reusable parser
    batchParser := parser.NewParser().
        WithSkipComments(true).
        WithCollectRawContent(false).
        WithParseTasks(false)

    for _, file := range files {
        result, err := batchParser.ParseFile(file)
        if err != nil {
            log.Printf("Failed to parse %s: %v", file, err)
            continue
        }
        
        // Process result
        fmt.Printf("File: %s, Dependencies: %d\n", 
            file, len(result.Project.Dependencies))
    }
    
    return nil
}
```

## Feature-Specific Configuration

### Dependency Analysis Only

```go
// Configuration for dependency analysis
depAnalysisParser := parser.NewParser().
    WithParsePlugins(false).          // Skip plugins
    WithParseRepositories(false).     // Skip repositories
    WithParseTasks(false).            // Skip tasks
    WithSkipComments(true).           // Skip comments
    WithCollectRawContent(false)      // Don't store raw content

result, err := depAnalysisParser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Only dependencies are parsed
dependencies := result.Project.Dependencies
fmt.Printf("Found %d dependencies\n", len(dependencies))
```

### Plugin Analysis Only

```go
// Configuration for plugin analysis
pluginAnalysisParser := parser.NewParser().
    WithParseDependencies(false).     // Skip dependencies
    WithParseRepositories(false).     // Skip repositories
    WithParseTasks(false).            // Skip tasks
    WithSkipComments(true).           // Skip comments
    WithCollectRawContent(false)      // Don't store raw content

result, err := pluginAnalysisParser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Only plugins are parsed
plugins := result.Project.Plugins
fmt.Printf("Found %d plugins\n", len(plugins))
```

### Complete Analysis

```go
// Configuration for comprehensive analysis
completeParser := parser.NewParser().
    WithSkipComments(false).          // Process comments
    WithCollectRawContent(true).      // Store original content
    WithParsePlugins(true).           // Parse plugins
    WithParseDependencies(true).      // Parse dependencies
    WithParseRepositories(true).      // Parse repositories
    WithParseTasks(true)              // Parse tasks

result, err := completeParser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

// All information is available
project := result.Project
fmt.Printf("Complete analysis:\n")
fmt.Printf("  Plugins: %d\n", len(project.Plugins))
fmt.Printf("  Dependencies: %d\n", len(project.Dependencies))
fmt.Printf("  Repositories: %d\n", len(project.Repositories))
fmt.Printf("  Tasks: %d\n", len(project.Tasks))
fmt.Printf("  Raw content: %d bytes\n", len(result.RawText))
```

## Use Case Configurations

### CI/CD Pipeline

Configuration optimized for continuous integration:

```go
// CI/CD optimized configuration
ciParser := parser.NewParser().
    WithSkipComments(true).           // Skip comments for speed
    WithCollectRawContent(false).     // Save memory
    WithParsePlugins(true).           // Need plugin info
    WithParseDependencies(true).      // Need dependency info
    WithParseRepositories(false).     // Usually not needed in CI
    WithParseTasks(false)             // Usually not needed in CI

// Use for dependency vulnerability scanning
result, err := ciParser.ParseFile("build.gradle")
if err != nil {
    return fmt.Errorf("CI parse failed: %v", err)
}

// Analyze dependencies for security issues
dependencies := result.Project.Dependencies
for _, dep := range dependencies {
    // Check against vulnerability database
    checkVulnerabilities(dep)
}
```

### IDE Integration

Configuration for IDE features:

```go
// IDE integration configuration
ideParser := parser.NewParser().
    WithSkipComments(false).          // Keep comments for documentation
    WithCollectRawContent(true).      // Need original content for editing
    WithParsePlugins(true).           // Need all plugin info
    WithParseDependencies(true).      // Need all dependency info
    WithParseRepositories(true).      // Need repository info
    WithParseTasks(true)              // Need task info for completion

// Parse with source mapping for precise editing
sourceParser := parser.NewSourceAwareParser()
result, err := sourceParser.ParseWithSourceMapping(content)
if err != nil {
    return fmt.Errorf("IDE parse failed: %v", err)
}

// Use source mapping for code completion, navigation, etc.
sourceMappedProject := result.SourceMappedProject
```

### Build Tool Integration

Configuration for build tools:

```go
// Build tool configuration
buildToolParser := parser.NewParser().
    WithSkipComments(true).           // Comments not needed
    WithCollectRawContent(false).     // Save memory
    WithParsePlugins(true).           // Need plugin configurations
    WithParseDependencies(true).      // Need dependency resolution
    WithParseRepositories(true).      // Need repository info
    WithParseTasks(false)             // Tasks handled by Gradle itself

result, err := buildToolParser.ParseFile("build.gradle")
if err != nil {
    return fmt.Errorf("build tool parse failed: %v", err)
}

// Extract information for build process
project := result.Project
repositories := project.Repositories
dependencies := project.Dependencies
```

## Advanced Configuration

### Custom Parser Factory

Create a factory for different parser configurations:

```go
type ParserFactory struct{}

func (pf *ParserFactory) CreateParser(configType string) *parser.GradleParser {
    switch configType {
    case "fast":
        return parser.NewParser().
            WithSkipComments(true).
            WithCollectRawContent(false).
            WithParseTasks(false)
    
    case "memory":
        return parser.NewParser().
            WithCollectRawContent(false).
            WithSkipComments(true)
    
    case "complete":
        return parser.NewParser().
            WithSkipComments(false).
            WithCollectRawContent(true).
            WithParseTasks(true)
    
    case "dependencies-only":
        return parser.NewParser().
            WithParsePlugins(false).
            WithParseRepositories(false).
            WithParseTasks(false).
            WithSkipComments(true).
            WithCollectRawContent(false)
    
    default:
        return parser.NewParser() // Default configuration
    }
}

// Usage
factory := &ParserFactory{}
fastParser := factory.CreateParser("fast")
result, err := fastParser.ParseFile("build.gradle")
```

### Configuration Validation

Validate configuration before use:

```go
func validateConfiguration(options *api.Options) error {
    // Check for conflicting options
    if !options.ParseDependencies && !options.ParsePlugins && 
       !options.ParseRepositories && !options.ParseTasks {
        return fmt.Errorf("at least one parsing feature must be enabled")
    }
    
    // Warn about performance implications
    if options.CollectRawContent && !options.SkipComments {
        log.Println("Warning: Collecting raw content with comment processing may impact performance")
    }
    
    return nil
}

// Usage
options := &api.Options{
    SkipComments:      true,
    CollectRawContent: false,
    ParsePlugins:      false,
    ParseDependencies: true,
    ParseRepositories: false,
    ParseTasks:        false,
}

if err := validateConfiguration(options); err != nil {
    log.Fatal(err)
}

parser := api.NewParser(options)
```

## Environment-Based Configuration

### Configuration from Environment Variables

```go
func createParserFromEnv() *parser.GradleParser {
    p := parser.NewParser()
    
    if os.Getenv("GRADLE_PARSER_SKIP_COMMENTS") == "true" {
        p = p.WithSkipComments(true)
    }
    
    if os.Getenv("GRADLE_PARSER_COLLECT_RAW") == "false" {
        p = p.WithCollectRawContent(false)
    }
    
    if os.Getenv("GRADLE_PARSER_SKIP_TASKS") == "true" {
        p = p.WithParseTasks(false)
    }
    
    return p
}

// Usage
parser := createParserFromEnv()
result, err := parser.ParseFile("build.gradle")
```

### Configuration from File

```go
type ParserConfig struct {
    SkipComments      bool `json:"skipComments"`
    CollectRawContent bool `json:"collectRawContent"`
    ParsePlugins      bool `json:"parsePlugins"`
    ParseDependencies bool `json:"parseDependencies"`
    ParseRepositories bool `json:"parseRepositories"`
    ParseTasks        bool `json:"parseTasks"`
}

func loadConfigFromFile(configFile string) (*ParserConfig, error) {
    data, err := os.ReadFile(configFile)
    if err != nil {
        return nil, err
    }
    
    var config ParserConfig
    err = json.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }
    
    return &config, nil
}

func createParserFromConfig(config *ParserConfig) *parser.GradleParser {
    return parser.NewParser().
        WithSkipComments(config.SkipComments).
        WithCollectRawContent(config.CollectRawContent).
        WithParsePlugins(config.ParsePlugins).
        WithParseDependencies(config.ParseDependencies).
        WithParseRepositories(config.ParseRepositories).
        WithParseTasks(config.ParseTasks)
}

// Usage
config, err := loadConfigFromFile("parser-config.json")
if err != nil {
    log.Fatal(err)
}

parser := createParserFromConfig(config)
result, err := parser.ParseFile("build.gradle")
```

## Performance Benchmarking

### Measuring Configuration Impact

```go
func benchmarkConfigurations(buildFile string) {
    configs := map[string]*parser.GradleParser{
        "default":     parser.NewParser(),
        "fast":        parser.NewParser().WithSkipComments(true).WithCollectRawContent(false),
        "memory":      parser.NewParser().WithCollectRawContent(false),
        "minimal":     parser.NewParser().WithSkipComments(true).WithCollectRawContent(false).WithParseTasks(false),
        "deps-only":   parser.NewParser().WithParsePlugins(false).WithParseRepositories(false).WithParseTasks(false),
    }
    
    for name, p := range configs {
        start := time.Now()
        result, err := p.ParseFile(buildFile)
        duration := time.Since(start)
        
        if err != nil {
            log.Printf("Config %s failed: %v", name, err)
            continue
        }
        
        fmt.Printf("Config %s: %v (%d deps, %d plugins)\n", 
            name, duration, len(result.Project.Dependencies), len(result.Project.Plugins))
    }
}
```

## Best Practices

1. **Choose appropriate configuration**: Match configuration to your use case
2. **Reuse parser instances**: Create parser once and reuse for multiple files
3. **Validate configuration**: Check for conflicting or invalid options
4. **Monitor performance**: Benchmark different configurations for your workload
5. **Document configuration**: Make configuration choices explicit in your code

## Next Steps

- [API Reference](../api/) - Complete API documentation
- [Examples](../examples/) - Practical usage examples
- [Advanced Features](./advanced-features.md) - Advanced parsing capabilities
