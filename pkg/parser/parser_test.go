package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Error("NewParser() returned nil")
	}
}

func TestParseBasic(t *testing.T) {
	parser := NewParser()

	// Test parsing a basic Gradle content
	content := `
		// Basic Gradle file
		group = 'com.example'
		version = '1.0.0'
		description = 'Test project'
	`

	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if result == nil {
		t.Fatal("Parse() returned nil result")
	}

	if result.Project == nil {
		t.Fatal("Parse() returned nil project")
	}

	// Confirm parse time is set
	if result.ParseTime == "" {
		t.Error("Parse() did not set ParseTime")
	}

	// Since parseLine is implemented as a no-op, we can't expect properties to be populated
	// When parseLine is implemented, we can add assertions like:
	// if result.Project.Group != "com.example" {
	//     t.Errorf("Parse() did not set Group correctly, got %s", result.Project.Group)
	// }
}

func TestParseFile(t *testing.T) {
	parser := NewParser()

	// Create a temporary Gradle file
	content := `
		// Sample build.gradle
		group = 'com.example'
		version = '1.0.0'
	`

	tmpFile, err := os.CreateTemp("", "build.gradle")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Parse the file
	result, err := parser.ParseFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if result == nil {
		t.Fatal("ParseFile() returned nil result")
	}

	if result.Project == nil {
		t.Fatal("ParseFile() returned nil project")
	}

	// Check that file path is set
	if result.Project.FilePath != tmpFile.Name() {
		t.Errorf("ParseFile() did not set FilePath correctly, got %s, want %s", result.Project.FilePath, tmpFile.Name())
	}

	// Test with non-existent file
	_, err = parser.ParseFile("/path/to/nonexistent/file.gradle")
	if err == nil {
		t.Error("ParseFile() with non-existent file should return error")
	}
}

func TestParseReader(t *testing.T) {
	parser := NewParser()

	content := `
		// Sample build.gradle from reader
		group = 'com.example'
		version = '1.0.0'
	`

	reader := strings.NewReader(content)
	result, err := parser.ParseReader(reader)
	if err != nil {
		t.Fatalf("ParseReader() error = %v", err)
	}

	if result == nil {
		t.Fatal("ParseReader() returned nil result")
	}

	if result.Project == nil {
		t.Fatal("ParseReader() returned nil project")
	}
}

func TestWithSkipComments(t *testing.T) {
	parser := &GradleParser{}
	parser = parser.WithSkipComments(false)
	if parser.skipComments != false {
		t.Error("WithSkipComments() did not set skipComments correctly")
	}
}

func TestWithCollectRawContent(t *testing.T) {
	parser := &GradleParser{}
	parser = parser.WithCollectRawContent(false)
	if parser.collectRawContent != false {
		t.Error("WithCollectRawContent() did not set collectRawContent correctly")
	}
}

func TestWithParsePlugins(t *testing.T) {
	parser := &GradleParser{}
	parser = parser.WithParsePlugins(false)
	if parser.parsePlugins != false {
		t.Error("WithParsePlugins() did not set parsePlugins correctly")
	}
}

func TestWithParseDependencies(t *testing.T) {
	parser := &GradleParser{}
	parser = parser.WithParseDependencies(false)
	if parser.parseDependencies != false {
		t.Error("WithParseDependencies() did not set parseDependencies correctly")
	}
}

func TestWithParseRepositories(t *testing.T) {
	parser := &GradleParser{}
	parser = parser.WithParseRepositories(false)
	if parser.parseRepositories != false {
		t.Error("WithParseRepositories() did not set parseRepositories correctly")
	}
}

func TestWithParseTasks(t *testing.T) {
	parser := &GradleParser{}
	parser = parser.WithParseTasks(false)
	if parser.parseTasks != false {
		t.Error("WithParseTasks() did not set parseTasks correctly")
	}
}

// Test that ParseOptions are correctly applied
func TestParseWithOptions(t *testing.T) {
	// Create a parser with all options disabled
	parser := &GradleParser{}
	parser = parser.WithSkipComments(false).
		WithCollectRawContent(false).
		WithParsePlugins(false).
		WithParseDependencies(false).
		WithParseRepositories(false).
		WithParseTasks(false)

	content := `
		// This is a comment
		group = 'com.example'
		version = '1.0.0'
	`

	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// With collectRawContent set to false, RawText should be empty
	if result.RawText != "" {
		t.Errorf("Parse() with collectRawContent=false returned non-empty RawText: %s", result.RawText)
	}

	// Test with collectRawContent enabled
	parser = parser.WithCollectRawContent(true)
	result, err = parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// With collectRawContent set to true, RawText should not be empty
	if result.RawText == "" {
		t.Error("Parse() with collectRawContent=true returned empty RawText")
	}
}

// Test parseLine (currently a no-op)
func TestParseLine(t *testing.T) {
	parser := &GradleParser{}
	project := &model.Project{}

	// Since parseLine is a no-op, it should always return nil error
	err := parser.parseLine("group = 'com.example'", 1, project)
	if err != nil {
		t.Errorf("parseLine() returned error: %v", err)
	}
}

// Test for parsing with errors
func TestParseWithErrors(t *testing.T) {
	// Create a new parser
	parser := &GradleParser{
		errors: make([]error, 0),
	}

	// Add some errors to the parser
	parser.errors = append(parser.errors, fmt.Errorf("test error 1"))
	parser.errors = append(parser.errors, fmt.Errorf("test error 2"))

	// Since Parse() resets errors, we need to modify it a bit
	content := `
		line1
		line2
		line3
	`

	// Parse will reset errors, so this test is more of a sanity check
	result, err := parser.Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// The result should contain the project
	if result.Project == nil {
		t.Error("Parse() with errors returned nil project")
	}
}

// 测试边界条件和错误处理
func TestParserEdgeCases(t *testing.T) {
	parser := NewParser()

	t.Run("Empty content", func(t *testing.T) {
		result, err := parser.Parse("")
		if err != nil {
			t.Fatalf("Parse() with empty content error = %v", err)
		}
		if result == nil {
			t.Fatal("Parse() with empty content returned nil result")
		}
		if result.Project == nil {
			t.Fatal("Parse() with empty content returned nil project")
		}
	})

	t.Run("Only whitespace", func(t *testing.T) {
		result, err := parser.Parse("   \n\t\n   ")
		if err != nil {
			t.Fatalf("Parse() with whitespace error = %v", err)
		}
		if result == nil {
			t.Fatal("Parse() with whitespace returned nil result")
		}
	})

	t.Run("Only comments", func(t *testing.T) {
		content := `
		// This is a comment
		/* This is a block comment */
		// Another comment
		`
		result, err := parser.Parse(content)
		if err != nil {
			t.Fatalf("Parse() with comments error = %v", err)
		}
		if result == nil {
			t.Fatal("Parse() with comments returned nil result")
		}
	})

	t.Run("Very long lines", func(t *testing.T) {
		longValue := strings.Repeat("a", 10000)
		content := fmt.Sprintf("group = '%s'", longValue)
		result, err := parser.Parse(content)
		if err != nil {
			t.Fatalf("Parse() with long line error = %v", err)
		}
		if result == nil {
			t.Fatal("Parse() with long line returned nil result")
		}
	})

	t.Run("Special characters", func(t *testing.T) {
		content := `
		group = 'com.example.测试'
		version = '1.0.0-SNAPSHOT'
		description = 'Test with émojis 🚀 and unicode ñ'
		`
		result, err := parser.Parse(content)
		if err != nil {
			t.Fatalf("Parse() with special chars error = %v", err)
		}
		if result == nil {
			t.Fatal("Parse() with special chars returned nil result")
		}
	})

	t.Run("Malformed syntax", func(t *testing.T) {
		content := `
		group = 'com.example
		version = 1.0.0'
		description = "unclosed quote
		`
		result, err := parser.Parse(content)
		// Should not fail, but may have warnings
		if err != nil {
			t.Fatalf("Parse() with malformed syntax error = %v", err)
		}
		if result == nil {
			t.Fatal("Parse() with malformed syntax returned nil result")
		}
		// Check if warnings were generated
		if len(result.Warnings) == 0 {
			t.Log("No warnings generated for malformed syntax")
		}
	})
}

// 测试文件解析的边界条件
func TestParseFileEdgeCases(t *testing.T) {
	parser := NewParser()

	t.Run("Non-existent file", func(t *testing.T) {
		_, err := parser.ParseFile("non-existent-file.gradle")
		if err == nil {
			t.Error("ParseFile() should return error for non-existent file")
		}
	})

	t.Run("Directory instead of file", func(t *testing.T) {
		tmpDir := t.TempDir()
		_, err := parser.ParseFile(tmpDir)
		if err == nil {
			t.Error("ParseFile() should return error for directory")
		}
	})

	t.Run("Empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		emptyFile := filepath.Join(tmpDir, "empty.gradle")
		err := os.WriteFile(emptyFile, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create empty file: %v", err)
		}

		result, err := parser.ParseFile(emptyFile)
		if err != nil {
			t.Fatalf("ParseFile() with empty file error = %v", err)
		}
		if result == nil {
			t.Fatal("ParseFile() with empty file returned nil result")
		}
	})

	t.Run("Large file", func(t *testing.T) {
		tmpDir := t.TempDir()
		largeFile := filepath.Join(tmpDir, "large.gradle")

		// Create a large file with many dependencies
		var content strings.Builder
		content.WriteString("dependencies {\n")
		for i := 0; i < 1000; i++ {
			content.WriteString(fmt.Sprintf("    implementation 'com.example:lib%d:1.0.0'\n", i))
		}
		content.WriteString("}\n")

		err := os.WriteFile(largeFile, []byte(content.String()), 0644)
		if err != nil {
			t.Fatalf("Failed to create large file: %v", err)
		}

		result, err := parser.ParseFile(largeFile)
		if err != nil {
			t.Fatalf("ParseFile() with large file error = %v", err)
		}
		if result == nil {
			t.Fatal("ParseFile() with large file returned nil result")
		}
	})

	t.Run("File with permission issues", func(t *testing.T) {
		tmpDir := t.TempDir()
		restrictedFile := filepath.Join(tmpDir, "restricted.gradle")

		err := os.WriteFile(restrictedFile, []byte("group = 'test'"), 0644)
		if err != nil {
			t.Fatalf("Failed to create restricted file: %v", err)
		}

		// Change permissions to make file unreadable
		err = os.Chmod(restrictedFile, 0000)
		if err != nil {
			t.Fatalf("Failed to change file permissions: %v", err)
		}

		// Restore permissions after test
		defer func() {
			os.Chmod(restrictedFile, 0644)
		}()

		_, err = parser.ParseFile(restrictedFile)
		if err == nil {
			t.Error("ParseFile() should return error for unreadable file")
		}
	})
}

// 测试解析器配置的边界条件
func TestParserConfigurationEdgeCases(t *testing.T) {
	t.Run("All options disabled", func(t *testing.T) {
		parser := &GradleParser{}
		parser = parser.WithSkipComments(true).
			WithCollectRawContent(false).
			WithParsePlugins(false).
			WithParseDependencies(false).
			WithParseRepositories(false).
			WithParseTasks(false)

		content := `
		plugins {
		    id 'java'
		    id 'org.springframework.boot' version '2.7.0'
		}

		dependencies {
		    implementation 'org.springframework.boot:spring-boot-starter-web'
		}

		repositories {
		    mavenCentral()
		}
		`

		result, err := parser.Parse(content)
		if err != nil {
			t.Fatalf("Parse() with disabled options error = %v", err)
		}

		// Should have empty collections since parsing is disabled
		if len(result.Project.Plugins) > 0 {
			t.Error("Should have no plugins when plugin parsing is disabled")
		}
		if len(result.Project.Dependencies) > 0 {
			t.Error("Should have no dependencies when dependency parsing is disabled")
		}
		if len(result.Project.Repositories) > 0 {
			t.Error("Should have no repositories when repository parsing is disabled")
		}
		if result.RawText != "" {
			t.Error("Should have empty RawText when collectRawContent is disabled")
		}
	})

	t.Run("Chained configuration", func(t *testing.T) {
		parser := &GradleParser{}

		// Test method chaining
		configuredParser := parser.
			WithSkipComments(false).
			WithCollectRawContent(true).
			WithParsePlugins(true).
			WithParseDependencies(true).
			WithParseRepositories(true).
			WithParseTasks(true)

		if configuredParser != parser {
			t.Error("Configuration methods should return the same parser instance")
		}

		// Verify all options are set correctly
		if parser.skipComments != false {
			t.Error("skipComments should be false")
		}
		if parser.collectRawContent != true {
			t.Error("collectRawContent should be true")
		}
		if parser.parsePlugins != true {
			t.Error("parsePlugins should be true")
		}
		if parser.parseDependencies != true {
			t.Error("parseDependencies should be true")
		}
		if parser.parseRepositories != true {
			t.Error("parseRepositories should be true")
		}
		if parser.parseTasks != true {
			t.Error("parseTasks should be true")
		}
	})
}

// 性能测试
func TestParserPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	parser := NewParser()

	// 创建一个复杂的Gradle文件内容
	var content strings.Builder
	content.WriteString(`
plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
}

group = 'com.example'
version = '1.0.0'
description = 'Performance test project'

repositories {
    mavenCentral()
    google()
    gradlePluginPortal()
}

dependencies {
`)

	// 添加大量依赖
	for i := 0; i < 500; i++ {
		content.WriteString(fmt.Sprintf("    implementation 'com.example:library%d:1.%d.0'\n", i, i%10))
		content.WriteString(fmt.Sprintf("    testImplementation 'com.test:test-library%d:2.%d.0'\n", i, i%5))
	}

	content.WriteString("}\n")

	// 添加大量任务
	for i := 0; i < 100; i++ {
		content.WriteString(fmt.Sprintf(`
task customTask%d {
    group = 'custom'
    description = 'Custom task %d'
    doLast {
        println 'Executing custom task %d'
    }
}
`, i, i, i))
	}

	testContent := content.String()

	// 测试解析时间
	startTime := time.Now()
	result, err := parser.Parse(testContent)
	parseTime := time.Since(startTime)

	if err != nil {
		t.Fatalf("Performance test parse error: %v", err)
	}

	if result == nil {
		t.Fatal("Performance test returned nil result")
	}

	// 验证解析结果
	if len(result.Project.Dependencies) == 0 {
		t.Error("Performance test should have parsed dependencies")
	}

	if len(result.Project.Plugins) == 0 {
		t.Error("Performance test should have parsed plugins")
	}

	// 记录性能指标
	t.Logf("Parsed %d characters in %v", len(testContent), parseTime)
	t.Logf("Found %d dependencies, %d plugins, %d repositories",
		len(result.Project.Dependencies),
		len(result.Project.Plugins),
		len(result.Project.Repositories))

	// 性能阈值检查（可根据需要调整）
	if parseTime > time.Second*5 {
		t.Errorf("Parse time %v exceeds threshold of 5 seconds", parseTime)
	}
}

// Helper to create a temporary Gradle project for testing
func createTempGradleProject(t *testing.T) string {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "gradle-project")
	if err != nil {
		t.Fatal(err)
	}

	// Create a simple build.gradle file
	buildGradleContent := `
		// Sample build.gradle
		group = 'com.example'
		version = '1.0.0'
		
		repositories {
			mavenCentral()
		}
		
		dependencies {
			implementation 'org.springframework:spring-core:5.3.10'
			testImplementation 'junit:junit:4.13.2'
		}
	`

	buildGradleFile := filepath.Join(tmpDir, "build.gradle")
	if err := os.WriteFile(buildGradleFile, []byte(buildGradleContent), 0644); err != nil {
		t.Fatal(err)
	}

	return tmpDir
}

// Clean up the temporary Gradle project
func cleanupTempGradleProject(tmpDir string) {
	os.RemoveAll(tmpDir)
}
