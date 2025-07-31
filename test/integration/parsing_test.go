package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/scagogogo/gradle-parser/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseRealGradleFiles tests parsing of real Gradle files
func TestParseRealGradleFiles(t *testing.T) {
	// Look for sample files in the examples directory
	sampleFilesDir := "../../examples/sample_files"

	if _, err := os.Stat(sampleFilesDir); os.IsNotExist(err) {
		t.Skip("Sample files directory not found, skipping integration tests")
	}

	testFiles := []struct {
		path        string
		description string
	}{
		{
			path:        filepath.Join(sampleFilesDir, "build.gradle"),
			description: "Main build.gradle",
		},
		{
			path:        filepath.Join(sampleFilesDir, "app", "build.gradle"),
			description: "App module build.gradle",
		},
		{
			path:        filepath.Join(sampleFilesDir, "common", "build.gradle"),
			description: "Common module build.gradle",
		},
	}

	for _, testFile := range testFiles {
		t.Run(testFile.description, func(t *testing.T) {
			if _, err := os.Stat(testFile.path); os.IsNotExist(err) {
				t.Skipf("Test file %s not found, skipping", testFile.path)
			}

			result, err := api.ParseFile(testFile.path)
			require.NoError(t, err, "Failed to parse %s", testFile.path)
			require.NotNil(t, result, "Parse result is nil for %s", testFile.path)
			require.NotNil(t, result.Project, "Project is nil for %s", testFile.path)

			// Basic validation
			project := result.Project
			assert.NotEmpty(t, project.FilePath, "FilePath should be set")

			// Log some basic info for debugging
			t.Logf("File: %s", testFile.path)
			t.Logf("  Group: %s", project.Group)
			t.Logf("  Name: %s", project.Name)
			t.Logf("  Version: %s", project.Version)
			t.Logf("  Dependencies: %d", len(project.Dependencies))
			t.Logf("  Plugins: %d", len(project.Plugins))
			t.Logf("  Repositories: %d", len(project.Repositories))
		})
	}
}

// TestCompleteWorkflow tests a complete parsing and analysis workflow
func TestCompleteWorkflow(t *testing.T) {
	sampleFile := "../../examples/sample_files/build.gradle"

	if _, err := os.Stat(sampleFile); os.IsNotExist(err) {
		t.Skip("Sample file not found, skipping workflow test")
	}

	// Step 1: Parse the file
	start := time.Now()
	result, err := api.ParseFile(sampleFile)
	parseTime := time.Since(start)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Project)

	t.Logf("Parse time: %v", parseTime)
	assert.Less(t, parseTime, 5*time.Second, "Parsing should be fast")

	project := result.Project

	// Step 2: Analyze project type
	plugins := project.Plugins

	projectTypes := []string{}
	if api.IsAndroidProject(plugins) {
		projectTypes = append(projectTypes, "Android")
	}
	if api.IsKotlinProject(plugins) {
		projectTypes = append(projectTypes, "Kotlin")
	}
	if api.IsSpringBootProject(plugins) {
		projectTypes = append(projectTypes, "Spring Boot")
	}

	t.Logf("Detected project types: %v", projectTypes)

	// Step 3: Analyze dependencies
	dependencies := project.Dependencies
	if len(dependencies) > 0 {
		depSets := api.DependenciesByScope(dependencies)

		t.Logf("Dependencies by scope:")
		for _, depSet := range depSets {
			t.Logf("  %s: %d dependencies", depSet.Scope, len(depSet.Dependencies))
		}

		// Check for common dependencies
		hasSpringBoot := false
		hasJUnit := false
		hasMySQL := false

		for _, dep := range dependencies {
			if dep.Group == "org.springframework.boot" {
				hasSpringBoot = true
			}
			if dep.Group == "junit" || dep.Group == "org.junit.jupiter" {
				hasJUnit = true
			}
			if dep.Group == "mysql" {
				hasMySQL = true
			}
		}

		t.Logf("Common dependencies found: Spring Boot=%v, JUnit=%v, MySQL=%v",
			hasSpringBoot, hasJUnit, hasMySQL)
	}

	// Step 4: Validate structure
	assert.NotEmpty(t, project.FilePath, "FilePath should be set")

	// If it's a Spring Boot project, it should have some dependencies
	if api.IsSpringBootProject(plugins) {
		assert.Greater(t, len(dependencies), 0, "Spring Boot project should have dependencies")
	}
}

// BenchmarkParsing benchmarks parsing performance
func BenchmarkParsing(b *testing.B) {
	sampleFile := "../../examples/sample_files/build.gradle"

	if _, err := os.Stat(sampleFile); os.IsNotExist(err) {
		b.Skip("Sample file not found, skipping benchmark")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := api.ParseFile(sampleFile)
		if err != nil {
			b.Fatalf("Parse error: %v", err)
		}
		if result == nil || result.Project == nil {
			b.Fatal("Invalid parse result")
		}
	}
}

// BenchmarkParsingWithOptions benchmarks parsing with different options
func BenchmarkParsingWithOptions(b *testing.B) {
	sampleFile := "../../examples/sample_files/build.gradle"

	if _, err := os.Stat(sampleFile); os.IsNotExist(err) {
		b.Skip("Sample file not found, skipping benchmark")
	}

	benchmarks := []struct {
		name    string
		options *api.Options
	}{
		{
			name: "Default",
			options: &api.Options{
				SkipComments:      true,
				CollectRawContent: true,
				ParsePlugins:      true,
				ParseDependencies: true,
				ParseRepositories: true,
				ParseTasks:        true,
			},
		},
		{
			name: "Fast",
			options: &api.Options{
				SkipComments:      true,
				CollectRawContent: false,
				ParsePlugins:      false,
				ParseDependencies: true,
				ParseRepositories: false,
				ParseTasks:        false,
			},
		},
		{
			name: "Memory Optimized",
			options: &api.Options{
				SkipComments:      true,
				CollectRawContent: false,
				ParsePlugins:      true,
				ParseDependencies: true,
				ParseRepositories: true,
				ParseTasks:        false,
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			parser := api.NewParser(bm.options)

			// Read file content once
			content, err := os.ReadFile(sampleFile)
			if err != nil {
				b.Fatalf("Failed to read file: %v", err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				result, err := parser.Parse(string(content))
				if err != nil {
					b.Fatalf("Parse error: %v", err)
				}
				if result == nil || result.Project == nil {
					b.Fatal("Invalid parse result")
				}
			}
		})
	}
}

// TestParsingMemoryUsage tests memory usage during parsing
func TestParsingMemoryUsage(t *testing.T) {
	sampleFile := "../../examples/sample_files/build.gradle"

	if _, err := os.Stat(sampleFile); os.IsNotExist(err) {
		t.Skip("Sample file not found, skipping memory test")
	}

	// Test with different options to see memory impact
	options := []*api.Options{
		{
			SkipComments:      true,
			CollectRawContent: true,
			ParsePlugins:      true,
			ParseDependencies: true,
			ParseRepositories: true,
			ParseTasks:        true,
		},
		{
			SkipComments:      true,
			CollectRawContent: false,
			ParsePlugins:      true,
			ParseDependencies: true,
			ParseRepositories: true,
			ParseTasks:        false,
		},
	}

	for i, opt := range options {
		t.Run(fmt.Sprintf("Config%d", i+1), func(t *testing.T) {
			parser := api.NewParser(opt)

			result, err := parser.ParseFile(sampleFile)
			require.NoError(t, err)
			require.NotNil(t, result)

			// Basic validation
			assert.NotNil(t, result.Project)

			if opt.CollectRawContent {
				assert.NotEmpty(t, result.RawText)
			} else {
				assert.Empty(t, result.RawText)
			}
		})
	}
}

// TestConcurrentParsing tests concurrent parsing safety
func TestConcurrentParsing(t *testing.T) {
	sampleFile := "../../examples/sample_files/build.gradle"

	if _, err := os.Stat(sampleFile); os.IsNotExist(err) {
		t.Skip("Sample file not found, skipping concurrent test")
	}

	const numGoroutines = 10
	const numIterations = 5

	results := make(chan error, numGoroutines*numIterations)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numIterations; j++ {
				result, err := api.ParseFile(sampleFile)
				if err != nil {
					results <- err
					return
				}
				if result == nil || result.Project == nil {
					results <- fmt.Errorf("invalid result")
					return
				}
			}
			results <- nil
		}()
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		assert.NoError(t, err, "Concurrent parsing failed")
	}
}
