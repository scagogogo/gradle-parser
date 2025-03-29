package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
