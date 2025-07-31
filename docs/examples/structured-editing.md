# Structured Editing Examples

Examples of programmatically modifying Gradle build files while preserving formatting.

## Quick Updates

### Update Dependency Version

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
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

    fmt.Println("✅ Successfully updated MySQL connector to v8.0.31")
}
```

## Advanced Editing

### Batch Updates with Editor

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    editor, err := api.CreateGradleEditor("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    // Update multiple dependencies
    updates := map[string]string{
        "mysql:mysql-connector-java":        "8.0.31",
        "org.springframework:spring-core":   "5.3.21",
        "com.google.guava:guava":           "31.1-jre",
    }

    for depKey, version := range updates {
        parts := strings.Split(depKey, ":")
        if len(parts) == 2 {
            err = editor.UpdateDependencyVersion(parts[0], parts[1], version)
            if err != nil {
                log.Printf("Failed to update %s: %v", depKey, err)
            } else {
                fmt.Printf("✅ Updated %s to %s\n", depKey, version)
            }
        }
    }

    // Apply all changes
    newContent, err := editor.ApplyModifications()
    if err != nil {
        log.Fatal(err)
    }

    // Save to file
    err = os.WriteFile("build.gradle", []byte(newContent), 0644)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("🎉 All updates applied successfully!")
}
```

## Next Steps

- [Custom Parser Examples](./custom-parser.md)
- [Basic Parsing Examples](./basic-parsing.md)
