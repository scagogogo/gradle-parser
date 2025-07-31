# Gradle Parser

[![CI](https://github.com/scagogogo/gradle-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/gradle-parser/actions/workflows/ci.yml)
[![Quality](https://github.com/scagogogo/gradle-parser/actions/workflows/quality.yml/badge.svg)](https://github.com/scagogogo/gradle-parser/actions/workflows/quality.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/gradle-parser)](https://goreportcard.com/report/github.com/scagogogo/gradle-parser)
[![GoDoc](https://godoc.org/github.com/scagogogo/gradle-parser?status.svg)](https://pkg.go.dev/github.com/scagogogo/gradle-parser)
[![codecov](https://codecov.io/gh/scagogogo/gradle-parser/branch/main/graph/badge.svg)](https://codecov.io/gh/scagogogo/gradle-parser)
[![License](https://img.shields.io/github/license/scagogogo/gradle-parser)](/LICENSE)
[![Release](https://img.shields.io/github/v/release/scagogogo/gradle-parser)](https://github.com/scagogogo/gradle-parser/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/scagogogo/gradle-parser)](https://golang.org/)

A powerful Go library for parsing and manipulating Gradle build files. Extract dependencies, plugins, repositories, and other configuration data with ease. Features structured editing capabilities for programmatic modifications while preserving original formatting.

## 📚 Documentation

> **🌟 Complete Documentation: [https://scagogogo.github.io/gradle-parser/](https://scagogogo.github.io/gradle-parser/)**
>
> 📖 [中文文档](README_zh.md) | 🚀 [Getting Started](https://scagogogo.github.io/gradle-parser/guide/getting-started.html) | 📋 [API Reference](https://scagogogo.github.io/gradle-parser/api/) | 💡 [Examples](https://scagogogo.github.io/gradle-parser/examples/)

## ✨ Features

### 🔍 **Comprehensive Parsing**
- Parse Gradle build files (supports both `build.gradle` and `build.gradle.kts`)
- Extract project metadata (group, name, version, description)
- Parse and categorize dependencies with scope classification
- Analyze plugin configurations and detect project types (Android/Kotlin/Spring Boot)
- Process repository configurations including custom repositories and authentication

### ✏️ **Structured Editing**
- **Precise modifications**: Update dependency versions, plugin versions, project properties
- **Minimal diff**: Preserve original formatting, modify only necessary parts
- **Source location tracking**: Record exact positions of elements in source files
- **Batch operations**: Apply multiple modifications in a single operation

### 🛠️ **Advanced Capabilities**
- Support for both Groovy DSL and Kotlin DSL syntax
- Customizable parser configuration for different requirements
- Multi-module Gradle project support
- Comprehensive error handling and validation

## 🚀 Quick Start

### Installation

```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    // Parse a Gradle file
    result, err := api.ParseFile("path/to/build.gradle")
    if err != nil {
        panic(err)
    }

    // Access project information
    project := result.Project
    fmt.Printf("Project: %s\n", project.Name)
    fmt.Printf("Group: %s\n", project.Group)
    fmt.Printf("Version: %s\n", project.Version)

    // List dependencies
    for _, dep := range project.Dependencies {
        fmt.Printf("Dependency: %s:%s:%s (%s)\n",
            dep.Group, dep.Name, dep.Version, dep.Scope)
    }

    // List plugins
    for _, plugin := range project.Plugins {
        fmt.Printf("Plugin: %s", plugin.ID)
        if plugin.Version != "" {
            fmt.Printf(" v%s", plugin.Version)
        }
        fmt.Println()
    }

    // List repositories
    for _, repo := range project.Repositories {
        fmt.Printf("Repository: %s (%s)\n", repo.Name, repo.URL)
    }
}
```

## 📖 Key Features

### 🔍 **Dependency Analysis**

```go
// Extract dependencies directly
dependencies, err := api.GetDependencies("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Group dependencies by scope
dependencySets := api.DependenciesByScope(dependencies)
for _, set := range dependencySets {
    fmt.Printf("Scope: %s\n", set.Scope)
    for _, dep := range set.Dependencies {
        fmt.Printf("  %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
    }
}
```

### 🔌 **Plugin Detection**

```go
// Extract plugin information
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Detect project types
if api.IsAndroidProject(plugins) {
    fmt.Println("Android project detected")
}
if api.IsKotlinProject(plugins) {
    fmt.Println("Kotlin project detected")
}
if api.IsSpringBootProject(plugins) {
    fmt.Println("Spring Boot project detected")
}
```

### 📝 **Repository Configuration**

```go
// Extract repository configurations
repos, err := api.GetRepositories("build.gradle")
if err != nil {
    log.Fatal(err)
}

for _, repo := range repos {
    fmt.Printf("Repository: %s (%s)\n", repo.Name, repo.URL)
}
```

### ✏️ **Structured Editing**

```go
// Simple version updates
newText, err := api.UpdateDependencyVersion("build.gradle", "mysql", "mysql-connector-java", "8.0.31")
if err != nil {
    log.Fatal(err)
}

newText, err = api.UpdatePluginVersion("build.gradle", "org.springframework.boot", "2.7.2")
if err != nil {
    log.Fatal(err)
}
```

### 🛠️ **Advanced Editing**

```go
// Create an editor for batch modifications
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// Perform multiple modifications
editor.UpdateProperty("version", "1.0.0")
editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
editor.UpdateDependencyVersion("com.google.guava", "guava", "31.1-jre")
editor.AddDependency("org.apache.commons", "commons-text", "1.9", "implementation")

// Apply all modifications
serializer := editor.NewGradleSerializer(editor.GetSourceMappedProject().OriginalText)
finalText, err := serializer.ApplyModifications(editor.GetModifications())
if err != nil {
    log.Fatal(err)
}

// Generate diff for review
diffLines := serializer.GenerateDiff(editor.GetModifications())
for _, diffLine := range diffLines {
    fmt.Println(diffLine.String())
}
```

### ⚙️ **Custom Parser Configuration**

```go
// Create custom parser options
options := api.DefaultOptions()
options.SkipComments = true
options.ParsePlugins = true
options.ParseDependencies = true
options.ParseRepositories = true

// Create custom parser
parser := api.NewParser(options)
result, err := parser.ParseFile("build.gradle")
if err != nil {
    log.Fatal(err)
}
```

## 🏗️ Project Structure

The project uses a modular design with clean separation of concerns:

```
├── pkg/                  # Core packages
│   ├── api/              # Main API interface
│   ├── config/           # Configuration parsing
│   ├── dependency/       # Dependency parsing
│   ├── editor/           # Structured editor
│   ├── model/            # Data models
│   ├── parser/           # Parser core
│   └── util/             # Utility functions
└── examples/             # Example code
    ├── 01_basic/         # Basic usage
    ├── 02_dependencies/  # Dependency extraction
    ├── 03_plugins/       # Plugin extraction
    ├── 04_repositories/  # Repository extraction
    ├── 05_complete/      # Complete features
    ├── 06_editor/        # Structured editing
    └── sample_files/     # Sample Gradle files
```

## 📚 Resources

### 📖 **Examples**
Explore the [examples](examples/) directory for comprehensive code samples demonstrating different features:

- **Basic Usage**: Simple parsing and data extraction
- **Advanced Features**: Complex parsing scenarios and customization
- **Structured Editing**: Programmatic file modifications
- **Project Analysis**: Complete project inspection and reporting

### 🧪 **Testing**
Run the comprehensive test suite:

```bash
# Run all tests
go test ./...

# Run test suite with coverage
cd test && ./scripts/run-tests.sh

# Run examples
cd examples && ./run-all-examples.sh
```

**Test Coverage**: Target >90%

### 🔄 **Continuous Integration**
This project uses GitHub Actions for comprehensive quality assurance:

- **🔄 CI**: Multi-version Go testing, code validation, example verification
- **📊 Quality**: Code coverage, security scanning, complexity analysis
- **📚 Docs**: Documentation building and deployment
- **🚀 Release**: Automated releases and asset building

**Quality Standards:**
- ✅ Unit and integration tests
- 🔍 Code quality checks (golangci-lint)
- 🛡️ Security vulnerability scanning
- 📈 Performance benchmarking
- 📝 Documentation link validation

## 👥 Contributors

Thanks to all the amazing people who have contributed to this project! 🎉

<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%">
        <a href="https://github.com/CC11001100">
          <img src="https://avatars.githubusercontent.com/u/12819457?v=4" width="100px;" alt="CC11001100"/>
          <br />
          <sub><b>CC11001100</b></sub>
        </a>
        <br />
        <sub>💻 📖 🎨 🚧</sub>
        <br />
        <sub>27 commits</sub>
      </td>
      <td align="center" valign="top" width="14.28%">
        <a href="https://github.com/AdamKorcz">
          <img src="https://avatars.githubusercontent.com/u/44787359?v=4" width="100px;" alt="AdamKorcz"/>
          <br />
          <sub><b>Adam Korczynski</b></sub>
        </a>
        <br />
        <sub>🐛 🔧</sub>
        <br />
        <sub>1 commit</sub>
      </td>
    </tr>
  </tbody>
</table>

### 🏆 Contribution Types

- 💻 Code
- 📖 Documentation
- 🎨 Design
- 🚧 Maintenance
- 🐛 Bug fixes
- 🔧 Tools

## 🤝 Contributing

Contributions are welcome! Please see the [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/scagogogo/gradle-parser.git
cd gradle-parser

# Install dependencies
go mod download

# Run tests
go test ./...

# Try examples
cd examples/01_basic && go run main.go
```

### Reporting Issues

Found a bug or have a feature request? Please report it in [GitHub Issues](https://github.com/scagogogo/gradle-parser/issues).

## 🗺️ Roadmap

- [ ] Enhanced Kotlin DSL support
- [ ] Performance optimizations
- [ ] More editing capabilities
- [ ] IDE plugin support
- [ ] Additional Gradle DSL syntax support

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details.

---

<div align="center">

**⭐ If this project helps you, please give it a star!**

Made with ❤️ by [scagogogo](https://github.com/scagogogo)

</div>