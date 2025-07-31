# Dependency Analysis Examples

This page provides examples of analyzing dependencies in Gradle projects.

## Basic Dependency Extraction

### List All Dependencies

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    deps, err := api.GetDependencies("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d dependencies:\n\n", len(deps))
    
    for i, dep := range deps {
        fmt.Printf("%d. %s:%s:%s (%s)\n", 
            i+1, dep.Group, dep.Name, dep.Version, dep.Scope)
    }
}
```

## Dependency Grouping

### Group by Scope

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    deps, err := api.GetDependencies("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    depSets := api.DependenciesByScope(deps)
    
    for _, depSet := range depSets {
        fmt.Printf("\n=== %s Dependencies (%d) ===\n", 
            depSet.Scope, len(depSet.Dependencies))
        
        for _, dep := range depSet.Dependencies {
            fmt.Printf("  %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
        }
    }
}
```

### Group by Organization

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    deps, err := api.GetDependencies("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    // Group by organization (first part of group)
    orgGroups := make(map[string][]*api.Dependency)
    
    for _, dep := range deps {
        org := strings.Split(dep.Group, ".")[0]
        orgGroups[org] = append(orgGroups[org], dep)
    }

    fmt.Println("Dependencies by Organization:")
    for org, orgDeps := range orgGroups {
        fmt.Printf("\n=== %s (%d dependencies) ===\n", org, len(orgDeps))
        for _, dep := range orgDeps {
            fmt.Printf("  %s:%s:%s (%s)\n", 
                dep.Group, dep.Name, dep.Version, dep.Scope)
        }
    }
}
```

## Dependency Analysis

### Find Outdated Dependencies

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

// Mock function to check if a version is outdated
func isOutdated(dep *api.Dependency) (bool, string) {
    // This is a simplified example
    // In real usage, you'd check against a dependency database
    outdatedVersions := map[string]string{
        "mysql:mysql-connector-java": "8.0.31",
        "org.springframework:spring-core": "5.3.23",
        "com.google.guava:guava": "31.1-jre",
    }
    
    key := fmt.Sprintf("%s:%s", dep.Group, dep.Name)
    if latestVersion, exists := outdatedVersions[key]; exists {
        if dep.Version != latestVersion {
            return true, latestVersion
        }
    }
    
    return false, ""
}

func main() {
    deps, err := api.GetDependencies("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Checking for outdated dependencies...\n")
    
    outdatedCount := 0
    for _, dep := range deps {
        if outdated, latestVersion := isOutdated(dep); outdated {
            fmt.Printf("⚠️  %s:%s\n", dep.Group, dep.Name)
            fmt.Printf("   Current: %s → Latest: %s\n", dep.Version, latestVersion)
            outdatedCount++
        }
    }
    
    if outdatedCount == 0 {
        fmt.Println("✅ All dependencies are up to date!")
    } else {
        fmt.Printf("\n📊 Found %d outdated dependencies\n", outdatedCount)
    }
}
```

### Security Vulnerability Check

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

// Mock vulnerability database
var vulnerabilities = map[string][]string{
    "mysql:mysql-connector-java": {"8.0.28", "8.0.27"},
    "org.apache.commons:commons-text": {"1.8", "1.7"},
}

func checkVulnerabilities(dep *api.Dependency) []string {
    key := fmt.Sprintf("%s:%s", dep.Group, dep.Name)
    if vulnVersions, exists := vulnerabilities[key]; exists {
        for _, vulnVersion := range vulnVersions {
            if dep.Version == vulnVersion {
                return vulnVersions
            }
        }
    }
    return nil
}

func main() {
    deps, err := api.GetDependencies("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("🔍 Scanning for security vulnerabilities...\n")
    
    vulnerableCount := 0
    for _, dep := range deps {
        if vulnVersions := checkVulnerabilities(dep); vulnVersions != nil {
            fmt.Printf("🚨 VULNERABILITY FOUND\n")
            fmt.Printf("   Dependency: %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
            fmt.Printf("   Vulnerable versions: %v\n", vulnVersions)
            fmt.Printf("   Recommendation: Update to latest version\n\n")
            vulnerableCount++
        }
    }
    
    if vulnerableCount == 0 {
        fmt.Println("✅ No known vulnerabilities found!")
    } else {
        fmt.Printf("⚠️  Found %d vulnerable dependencies\n", vulnerableCount)
    }
}
```

## Advanced Analysis

### Dependency Tree Analysis

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func analyzeDependencyTree(filePath string) {
    result, err := api.ParseFile(filePath)
    if err != nil {
        log.Fatal(err)
    }

    dependencies := result.Project.Dependencies
    
    fmt.Printf("📊 Dependency Analysis for %s\n", filePath)
    fmt.Printf("=====================================\n\n")

    // Group by scope
    depSets := api.DependenciesByScope(dependencies)
    
    for _, depSet := range depSets {
        fmt.Printf("📦 %s Dependencies (%d)\n", depSet.Scope, len(depSet.Dependencies))
        fmt.Println(strings.Repeat("-", 40))
        
        // Group by organization within scope
        orgGroups := make(map[string][]*api.Dependency)
        for _, dep := range depSet.Dependencies {
            org := strings.Split(dep.Group, ".")[0]
            orgGroups[org] = append(orgGroups[org], dep)
        }
        
        for org, orgDeps := range orgGroups {
            fmt.Printf("  🏢 %s (%d)\n", org, len(orgDeps))
            for _, dep := range orgDeps {
                fmt.Printf("    └── %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
            }
        }
        fmt.Println()
    }
}

func main() {
    analyzeDependencyTree("build.gradle")
}
```

### License Analysis

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

// Mock license database
var licenses = map[string]string{
    "org.springframework:spring-core": "Apache-2.0",
    "mysql:mysql-connector-java": "GPL-2.0",
    "com.google.guava:guava": "Apache-2.0",
    "org.junit.jupiter:junit-jupiter-api": "EPL-2.0",
}

func getLicense(dep *api.Dependency) string {
    key := fmt.Sprintf("%s:%s", dep.Group, dep.Name)
    if license, exists := licenses[key]; exists {
        return license
    }
    return "Unknown"
}

func main() {
    deps, err := api.GetDependencies("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("📄 License Analysis")
    fmt.Println("===================\n")

    licenseCount := make(map[string]int)
    licenseGroups := make(map[string][]*api.Dependency)

    for _, dep := range deps {
        license := getLicense(dep)
        licenseCount[license]++
        licenseGroups[license] = append(licenseGroups[license], dep)
    }

    // Summary
    fmt.Println("License Summary:")
    for license, count := range licenseCount {
        fmt.Printf("  %s: %d dependencies\n", license, count)
    }

    // Detailed breakdown
    fmt.Println("\nDetailed Breakdown:")
    for license, deps := range licenseGroups {
        fmt.Printf("\n📋 %s License:\n", license)
        for _, dep := range deps {
            fmt.Printf("  • %s:%s:%s\n", dep.Group, dep.Name, dep.Version)
        }
    }
}
```

## Comparison Analysis

### Compare Dependencies Between Projects

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func compareDependencies(file1, file2 string) {
    deps1, err := api.GetDependencies(file1)
    if err != nil {
        log.Fatal(err)
    }

    deps2, err := api.GetDependencies(file2)
    if err != nil {
        log.Fatal(err)
    }

    // Create maps for easier comparison
    depMap1 := make(map[string]*api.Dependency)
    depMap2 := make(map[string]*api.Dependency)

    for _, dep := range deps1 {
        key := fmt.Sprintf("%s:%s", dep.Group, dep.Name)
        depMap1[key] = dep
    }

    for _, dep := range deps2 {
        key := fmt.Sprintf("%s:%s", dep.Group, dep.Name)
        depMap2[key] = dep
    }

    fmt.Printf("🔍 Comparing %s vs %s\n", file1, file2)
    fmt.Println("=====================================\n")

    // Find common dependencies
    common := []string{}
    versionDiffs := []string{}

    for key, dep1 := range depMap1 {
        if dep2, exists := depMap2[key]; exists {
            common = append(common, key)
            if dep1.Version != dep2.Version {
                versionDiffs = append(versionDiffs, 
                    fmt.Sprintf("%s: %s vs %s", key, dep1.Version, dep2.Version))
            }
        }
    }

    // Find unique dependencies
    unique1 := []string{}
    unique2 := []string{}

    for key := range depMap1 {
        if _, exists := depMap2[key]; !exists {
            unique1 = append(unique1, key)
        }
    }

    for key := range depMap2 {
        if _, exists := depMap1[key]; !exists {
            unique2 = append(unique2, key)
        }
    }

    // Report results
    fmt.Printf("📊 Common dependencies: %d\n", len(common))
    fmt.Printf("🔄 Version differences: %d\n", len(versionDiffs))
    fmt.Printf("📦 Unique to %s: %d\n", file1, len(unique1))
    fmt.Printf("📦 Unique to %s: %d\n", file2, len(unique2))

    if len(versionDiffs) > 0 {
        fmt.Println("\n🔄 Version Differences:")
        for _, diff := range versionDiffs {
            fmt.Printf("  • %s\n", diff)
        }
    }

    if len(unique1) > 0 {
        fmt.Printf("\n📦 Unique to %s:\n", file1)
        for _, dep := range unique1 {
            fmt.Printf("  • %s\n", dep)
        }
    }

    if len(unique2) > 0 {
        fmt.Printf("\n📦 Unique to %s:\n", file2)
        for _, dep := range unique2 {
            fmt.Printf("  • %s\n", dep)
        }
    }
}

func main() {
    compareDependencies("build.gradle", "app/build.gradle")
}
```

## Next Steps

- [Plugin Detection Examples](./plugin-detection.md)
- [Repository Parsing Examples](./repository-parsing.md)
- [Structured Editing Examples](./structured-editing.md)
