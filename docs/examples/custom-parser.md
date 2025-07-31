# Custom Parser Examples

Examples of configuring the parser for specific use cases and performance optimization.

## Performance Optimization

### Fast Dependency Extraction

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/parser"
)

func main() {
    // Create optimized parser for dependency extraction only
    fastParser := parser.NewParser().
        WithSkipComments(true).           // Skip comments
        WithCollectRawContent(false).     // Don't store raw content
        WithParsePlugins(false).          // Skip plugins
        WithParseRepositories(false).     // Skip repositories
        WithParseTasks(false)             // Skip tasks

    result, err := fastParser.ParseFile("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    // Only dependencies are parsed
    dependencies := result.Project.Dependencies
    fmt.Printf("Found %d dependencies in optimized parse\n", len(dependencies))
    
    for _, dep := range dependencies {
        fmt.Printf("  %s:%s:%s (%s)\n", 
            dep.Group, dep.Name, dep.Version, dep.Scope)
    }
}
```

## Custom Configuration

### Memory-Optimized Parser

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/parser"
)

func main() {
    // Memory-efficient configuration
    memoryParser := parser.NewParser().
        WithCollectRawContent(false).     // Don't store raw content
        WithSkipComments(true).           // Skip comments
        WithParseTasks(false)             // Skip tasks if not needed

    result, err := memoryParser.ParseFile("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    project := result.Project
    fmt.Printf("Memory-optimized parse completed:\n")
    fmt.Printf("  Dependencies: %d\n", len(project.Dependencies))
    fmt.Printf("  Plugins: %d\n", len(project.Plugins))
    fmt.Printf("  Repositories: %d\n", len(project.Repositories))
}
```

## Next Steps

- [Basic Parsing Examples](./basic-parsing.md)
- [Dependency Analysis Examples](./dependency-analysis.md)
