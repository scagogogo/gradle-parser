# Examples

This section provides practical examples of using Gradle Parser in various scenarios. Each example includes complete, runnable code with explanations.

## Available Examples

### Basic Usage
- **[Basic Parsing](./basic-parsing.md)** - Simple file parsing and information extraction
- **[Dependency Analysis](./dependency-analysis.md)** - Working with project dependencies
- **[Plugin Detection](./plugin-detection.md)** - Analyzing plugins and project types
- **[Repository Parsing](./repository-parsing.md)** - Extracting repository configurations

### Advanced Features
- **[Structured Editing](./structured-editing.md)** - Modifying Gradle files programmatically
- **[Custom Parser](./custom-parser.md)** - Configuring parser options for specific needs

## Running Examples

All examples are designed to be self-contained and runnable. To run an example:

1. Create a new Go module:
```bash
mkdir gradle-parser-example
cd gradle-parser-example
go mod init example
```

2. Install Gradle Parser:
```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

3. Copy the example code into `main.go`

4. Create a sample `build.gradle` file (or use the provided samples)

5. Run the example:
```bash
go run main.go
```

## Sample Files

The examples use various sample Gradle files. You can find these in the project repository under `examples/sample_files/`:

- `build.gradle` - Standard Groovy DSL build file
- `build.gradle.kts` - Kotlin DSL build file  
- `app/build.gradle` - Android app module
- `common/build.gradle` - Library module
- `settings.gradle` - Multi-module project settings

## Example Categories

### 🔍 **Parsing & Analysis**
Learn how to parse Gradle files and extract information:
- Project metadata (name, version, group)
- Dependencies with scope analysis
- Plugin configurations
- Repository settings

### ✏️ **Editing & Modification**
Discover how to modify Gradle files programmatically:
- Update dependency versions
- Modify plugin configurations
- Add new dependencies
- Preserve formatting and minimize diffs

### 🛠️ **Advanced Usage**
Explore advanced features and customization:
- Custom parser configurations
- Source location tracking
- Error handling strategies
- Performance optimization

## Common Use Cases

### Build Tool Integration
```go
// Parse and analyze a project
result, err := api.ParseFile("build.gradle")
if err != nil {
    return err
}

// Check for outdated dependencies
for _, dep := range result.Project.Dependencies {
    if isOutdated(dep) {
        fmt.Printf("Outdated: %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
    }
}
```

### Dependency Management
```go
// Update all Spring Boot dependencies
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    return err
}

springBootVersion := "2.7.2"
err = editor.UpdatePluginVersion("org.springframework.boot", springBootVersion)
if err != nil {
    return err
}
```

### Project Analysis
```go
// Analyze project type and generate report
plugins, err := api.GetPlugins("build.gradle")
if err != nil {
    return err
}

if api.IsAndroidProject(plugins) {
    fmt.Println("Android project detected")
    // Android-specific analysis
}
```

## Tips for Examples

1. **Start Simple** - Begin with basic parsing examples before moving to advanced features
2. **Use Real Files** - Test examples with actual Gradle files from your projects
3. **Handle Errors** - Always include proper error handling in your code
4. **Experiment** - Modify the examples to explore different scenarios
5. **Check Output** - Verify that parsing results match your expectations

## Contributing Examples

Have a useful example to share? We welcome contributions! Please:

1. Follow the existing example format
2. Include complete, runnable code
3. Add clear explanations and comments
4. Test with various Gradle file formats
5. Submit a pull request with your example

## Getting Help

If you have questions about the examples or need help with specific use cases:

- Check the [API Reference](../api/) for detailed documentation
- Visit our [GitHub Discussions](https://github.com/scagogogo/gradle-parser/discussions)
- Open an [Issue](https://github.com/scagogogo/gradle-parser/issues) for bugs or feature requests
