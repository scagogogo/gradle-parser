# Getting Started

Welcome to Gradle Parser! This guide will help you get up and running with parsing Gradle build files in Go.

## Installation

Install Gradle Parser using Go modules:

```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

## Quick Start

Here's a simple example to get you started:

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    // Parse a Gradle file
    result, err := api.ParseFile("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    // Print project information
    project := result.Project
    fmt.Printf("Project: %s\n", project.Name)
    fmt.Printf("Group: %s\n", project.Group)
    fmt.Printf("Version: %s\n", project.Version)
    fmt.Printf("Description: %s\n", project.Description)
}
```

## Basic Parsing

### Parse from File

The most common way to parse a Gradle file:

```go
result, err := api.ParseFile("path/to/build.gradle")
if err != nil {
    log.Fatal(err)
}

// Access the parsed project
project := result.Project
```

### Parse from String

If you have Gradle content as a string:

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

### Parse from Reader

For streaming or other I/O sources:

```go
import "strings"

reader := strings.NewReader(gradleContent)
result, err := api.ParseReader(reader)
if err != nil {
    log.Fatal(err)
}
```

## Understanding the Result

The `ParseResult` contains:

- **Project**: The main project information and components
- **RawText**: Original file content (if enabled)
- **Errors**: Any parsing errors encountered
- **Warnings**: Non-fatal parsing warnings
- **ParseTime**: Time taken to parse the file

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Parse time: %s\n", result.ParseTime)
fmt.Printf("Warnings: %d\n", len(result.Warnings))

// Access project components
project := result.Project
fmt.Printf("Dependencies: %d\n", len(project.Dependencies))
fmt.Printf("Plugins: %d\n", len(project.Plugins))
fmt.Printf("Repositories: %d\n", len(project.Repositories))
```

## Error Handling

Gradle Parser provides detailed error information:

```go
result, err := api.ParseFile("build.gradle")
if err != nil {
    log.Printf("Failed to parse file: %v", err)
    return
}

// Check for parsing warnings
if len(result.Warnings) > 0 {
    fmt.Println("Parsing warnings:")
    for _, warning := range result.Warnings {
        fmt.Printf("  - %s\n", warning)
    }
}
```

## Next Steps

Now that you have the basics, explore more features:

- [Basic Usage](./basic-usage.md) - Learn about extracting specific information
- [Advanced Features](./advanced-features.md) - Discover powerful parsing capabilities
- [Structured Editing](./structured-editing.md) - Modify Gradle files programmatically
- [Configuration](./configuration.md) - Customize parser behavior

## Requirements

- Go 1.19 or later
- No external dependencies required

## Supported Formats

- Groovy DSL (`build.gradle`)
- Kotlin DSL (`build.gradle.kts`) - Basic support
- Single and multi-module projects
