---
layout: home

hero:
  name: "Gradle Parser"
  text: "Powerful Gradle Build File Parser"
  tagline: "Parse, analyze, and edit Gradle build files with ease in Go"
  image:
    src: /logo.svg
    alt: Gradle Parser
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/scagogogo/gradle-parser

features:
  - icon: 🚀
    title: Fast & Reliable
    details: High-performance parsing of Gradle build files with comprehensive error handling and validation.
  
  - icon: 🔍
    title: Deep Analysis
    details: Extract detailed information about dependencies, plugins, repositories, and project configuration.
  
  - icon: ✏️
    title: Structured Editing
    details: Modify Gradle files programmatically while preserving formatting and minimizing diffs.
  
  - icon: 📍
    title: Source Mapping
    details: Track precise source locations for every parsed element, enabling accurate modifications.
  
  - icon: 🌐
    title: Multi-Format Support
    details: Support for both Groovy DSL and Kotlin DSL syntax in Gradle build files.
  
  - icon: 🛠️
    title: Highly Configurable
    details: Customize parsing behavior with flexible options to meet your specific requirements.
---

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    // Parse a Gradle file
    result, err := api.ParseFile("build.gradle")
    if err != nil {
        panic(err)
    }

    // Access project information
    fmt.Printf("Project: %s\n", result.Project.Name)
    fmt.Printf("Version: %s\n", result.Project.Version)
    
    // List dependencies
    for _, dep := range result.Project.Dependencies {
        fmt.Printf("Dependency: %s:%s:%s (%s)\n", 
            dep.Group, dep.Name, dep.Version, dep.Scope)
    }
}
```

## Installation

```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

## Key Features

### 🔍 **Comprehensive Parsing**
- Extract project metadata (group, name, version, description)
- Parse dependencies with scope classification
- Analyze plugin configurations and detect project types
- Process repository configurations including custom repositories

### ✏️ **Structured Editing**
- Update dependency versions programmatically
- Modify plugin versions and configurations
- Edit project properties while preserving formatting
- Add new dependencies with proper placement

### 📍 **Source Location Tracking**
- Precise line and column information for every element
- Enable accurate modifications with minimal diff
- Support for complex editing scenarios

### 🌐 **Multi-Language Support**
- Full support for Groovy DSL syntax
- Kotlin DSL compatibility
- Handle both single and multi-module projects

## Use Cases

- **Build Tool Integration**: Integrate Gradle parsing into your build tools and IDEs
- **Dependency Management**: Analyze and update project dependencies programmatically  
- **Project Analysis**: Extract insights from Gradle projects for reporting and analytics
- **Automated Maintenance**: Bulk update dependencies and configurations across projects
- **Migration Tools**: Build tools to migrate between Gradle versions or configurations

## Community

- [GitHub Repository](https://github.com/scagogogo/gradle-parser)
- [Issue Tracker](https://github.com/scagogogo/gradle-parser/issues)
- [Discussions](https://github.com/scagogogo/gradle-parser/discussions)

## License

Gradle Parser is released under the [MIT License](https://github.com/scagogogo/gradle-parser/blob/main/LICENSE).
