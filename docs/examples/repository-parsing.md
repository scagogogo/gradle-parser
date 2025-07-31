# Repository Parsing Examples

Examples of parsing and analyzing repository configurations in Gradle projects.

## Basic Repository Extraction

### List All Repositories

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    repos, err := api.GetRepositories("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d repositories:\n\n", len(repos))
    
    for i, repo := range repos {
        fmt.Printf("%d. %s", i+1, repo.Name)
        if repo.URL != "" {
            fmt.Printf(" (%s)", repo.URL)
        }
        fmt.Printf(" [%s]\n", repo.Type)
    }
}
```

## Repository Analysis

### Analyze Repository Types

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
    repos, err := api.GetRepositories("build.gradle")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("📊 Repository Analysis")
    fmt.Println("======================\n")

    // Count by type
    typeCount := make(map[string]int)
    for _, repo := range repos {
        typeCount[repo.Type]++
    }

    fmt.Println("Repository Types:")
    for repoType, count := range typeCount {
        fmt.Printf("  %s: %d\n", repoType, count)
    }

    // Check for common repositories
    commonRepos := map[string]bool{
        "mavenCentral": false,
        "google":       false,
        "jcenter":      false,
    }

    for _, repo := range repos {
        if _, exists := commonRepos[repo.Name]; exists {
            commonRepos[repo.Name] = true
        }
    }

    fmt.Println("\nCommon Repositories:")
    for name, found := range commonRepos {
        status := "❌"
        if found {
            status = "✅"
        }
        fmt.Printf("  %s %s\n", status, name)
    }
}
```

## Next Steps

- [Structured Editing Examples](./structured-editing.md)
- [Custom Parser Examples](./custom-parser.md)
