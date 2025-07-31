# Basic Parsing Examples

This page provides practical examples of basic Gradle file parsing operations.

## Simple File Parsing

### Parse and Display Project Information

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    // Parse the build.gradle file
    result, err := api.ParseFile("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    project := result.Project

    // Display basic project information
    fmt.Println("=== Project Information ===")
    fmt.Printf("Name: %s\n", project.Name)
    fmt.Printf("Group: %s\n", project.Group)
    fmt.Printf("Version: %s\n", project.Version)
    fmt.Printf("Description: %s\n", project.Description)

    // Display Java compatibility if available
    if project.SourceCompatibility != "" {
        fmt.Printf("Source Compatibility: %s\n", project.SourceCompatibility)
    }
    if project.TargetCompatibility != "" {
        fmt.Printf("Target Compatibility: %s\n", project.TargetCompatibility)
    }

    // Display parsing statistics
    fmt.Printf("\n=== Parsing Statistics ===\n")
    fmt.Printf("Dependencies: %d\n", len(project.Dependencies))
    fmt.Printf("Plugins: %d\n", len(project.Plugins))
    fmt.Printf("Repositories: %d\n", len(project.Repositories))
    fmt.Printf("Parse Time: %s\n", result.ParseTime)

    if len(result.Warnings) > 0 {
        fmt.Printf("Warnings: %d\n", len(result.Warnings))
    }
}
```

### Parse from String Content

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    gradleContent := `
plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
}

group = 'com.example'
version = '1.0.0'
description = 'Demo project for Spring Boot'

java {
    sourceCompatibility = '11'
}

repositories {
    mavenCentral()
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
    implementation 'mysql:mysql-connector-java'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}
`

    result, err := api.ParseString(gradleContent)
    if err != nil {
        log.Fatal(err)
    }

    project := result.Project
    fmt.Printf("Parsed project: %s v%s\n", project.Name, project.Version)
    fmt.Printf("Found %d dependencies\n", len(project.Dependencies))
}
```

## Error Handling Examples

### Robust Error Handling

```go
package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func parseWithErrorHandling(filePath string) {
    result, err := api.ParseFile(filePath)
    if err != nil {
        // Handle different types of errors
        if os.IsNotExist(err) {
            log.Printf("File not found: %s", filePath)
            return
        }
        
        if strings.Contains(err.Error(), "permission") {
            log.Printf("Permission denied: %s", filePath)
            return
        }
        
        log.Printf("Parse error: %v", err)
        return
    }

    // Check for parsing warnings
    if len(result.Warnings) > 0 {
        fmt.Printf("⚠️  Parsing completed with %d warnings:\n", len(result.Warnings))
        for i, warning := range result.Warnings {
            fmt.Printf("  %d. %s\n", i+1, warning)
        }
        fmt.Println()
    }

    // Validate the parsed result
    project := result.Project
    if project == nil {
        log.Println("❌ No project information found")
        return
    }

    fmt.Printf("✅ Successfully parsed: %s\n", filePath)
    fmt.Printf("   Project: %s v%s\n", project.Name, project.Version)
}

func main() {
    files := []string{
        "build.gradle",
        "app/build.gradle",
        "non-existent.gradle",
    }

    for _, file := range files {
        parseWithErrorHandling(file)
        fmt.Println()
    }
}
```

## Parsing Different Sources

### Parse from HTTP Response

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func parseFromURL(url string) {
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Fatalf("HTTP error: %d", resp.StatusCode)
    }

    result, err := api.ParseReader(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    project := result.Project
    fmt.Printf("Remote project: %s v%s\n", project.Name, project.Version)
    fmt.Printf("Dependencies: %d\n", len(project.Dependencies))
}

func main() {
    // Example: Parse a build.gradle from GitHub
    url := "https://raw.githubusercontent.com/spring-projects/spring-boot/main/spring-boot-project/spring-boot-starters/spring-boot-starter/build.gradle"
    parseFromURL(url)
}
```

### Parse Multiple Files

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func parseMultipleFiles(files []string) {
    results := make(map[string]*api.ParseResult)
    
    for _, file := range files {
        result, err := api.ParseFile(file)
        if err != nil {
            log.Printf("Failed to parse %s: %v", file, err)
            continue
        }
        results[file] = result
    }

    // Summary report
    fmt.Println("=== Multi-File Parsing Summary ===")
    totalDeps := 0
    totalPlugins := 0
    
    for file, result := range results {
        project := result.Project
        deps := len(project.Dependencies)
        plugins := len(project.Plugins)
        
        fmt.Printf("%s:\n", filepath.Base(file))
        fmt.Printf("  Project: %s v%s\n", project.Name, project.Version)
        fmt.Printf("  Dependencies: %d\n", deps)
        fmt.Printf("  Plugins: %d\n", plugins)
        
        totalDeps += deps
        totalPlugins += plugins
    }
    
    fmt.Printf("\nTotals across %d files:\n", len(results))
    fmt.Printf("  Dependencies: %d\n", totalDeps)
    fmt.Printf("  Plugins: %d\n", totalPlugins)
}

func main() {
    files := []string{
        "build.gradle",
        "app/build.gradle",
        "lib/build.gradle",
    }
    
    parseMultipleFiles(files)
}
```

## Performance Examples

### Benchmarking Parse Performance

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func benchmarkParsing(filePath string, iterations int) {
    fmt.Printf("Benchmarking %s (%d iterations)...\n", filePath, iterations)
    
    var totalDuration time.Duration
    var successCount int
    
    for i := 0; i < iterations; i++ {
        start := time.Now()
        result, err := api.ParseFile(filePath)
        duration := time.Since(start)
        
        if err != nil {
            log.Printf("Iteration %d failed: %v", i+1, err)
            continue
        }
        
        totalDuration += duration
        successCount++
        
        if i == 0 {
            // Print details for first iteration
            project := result.Project
            fmt.Printf("  First parse: %s v%s (%d deps, %d plugins) in %v\n",
                project.Name, project.Version,
                len(project.Dependencies), len(project.Plugins),
                duration)
        }
    }
    
    if successCount > 0 {
        avgDuration := totalDuration / time.Duration(successCount)
        fmt.Printf("  Average: %v (%d/%d successful)\n", avgDuration, successCount, iterations)
    }
}

func main() {
    files := []string{
        "build.gradle",
        "app/build.gradle",
    }
    
    for _, file := range files {
        benchmarkParsing(file, 10)
        fmt.Println()
    }
}
```

## Validation Examples

### Project Structure Validation

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func validateProject(filePath string) {
    result, err := api.ParseFile(filePath)
    if err != nil {
        log.Fatal(err)
    }

    project := result.Project
    issues := []string{}

    // Check basic project information
    if project.Group == "" {
        issues = append(issues, "Missing project group")
    }
    
    if project.Version == "" {
        issues = append(issues, "Missing project version")
    }
    
    if project.Name == "" {
        issues = append(issues, "Missing project name")
    }

    // Check for essential plugins
    hasJavaPlugin := false
    for _, plugin := range project.Plugins {
        if plugin.ID == "java" || plugin.ID == "java-library" {
            hasJavaPlugin = true
            break
        }
    }
    
    if !hasJavaPlugin {
        issues = append(issues, "No Java plugin found")
    }

    // Check for dependencies
    if len(project.Dependencies) == 0 {
        issues = append(issues, "No dependencies found")
    }

    // Check for repositories
    if len(project.Repositories) == 0 {
        issues = append(issues, "No repositories configured")
    }

    // Report results
    fmt.Printf("Validating %s...\n", filePath)
    if len(issues) == 0 {
        fmt.Println("✅ Project structure is valid")
    } else {
        fmt.Printf("❌ Found %d issues:\n", len(issues))
        for i, issue := range issues {
            fmt.Printf("  %d. %s\n", i+1, issue)
        }
    }
}

func main() {
    validateProject("build.gradle")
}
```

## Sample build.gradle

For testing these examples, create a `build.gradle` file:

```gradle
plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
}

group = 'com.example'
version = '1.0.0'
description = 'Demo project for Spring Boot'

java {
    sourceCompatibility = '11'
}

repositories {
    mavenCentral()
    google()
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
    implementation 'mysql:mysql-connector-java:8.0.29'
    implementation 'com.google.guava:guava:31.1-jre'
    
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
    testImplementation 'org.junit.jupiter:junit-jupiter-api:5.8.2'
}

tasks.named('test') {
    useJUnitPlatform()
}
```

## Running the Examples

1. Create a new Go module:
```bash
mkdir gradle-parser-examples
cd gradle-parser-examples
go mod init examples
```

2. Install Gradle Parser:
```bash
go get github.com/scagogogo/gradle-parser/pkg/api
```

3. Create the sample `build.gradle` file above

4. Copy any example code into `main.go`

5. Run the example:
```bash
go run main.go
```
