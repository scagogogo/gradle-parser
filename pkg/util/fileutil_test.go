package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsBuildGradleFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{"build.gradle", "path/to/build.gradle", true},
		{"build.gradle.kts", "path/to/build.gradle.kts", true},
		{"other file", "path/to/other.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBuildGradleFile(tt.filePath); got != tt.want {
				t.Errorf("IsBuildGradleFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSettingsGradleFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{"settings.gradle", "path/to/settings.gradle", true},
		{"settings.gradle.kts", "path/to/settings.gradle.kts", true},
		{"other file", "path/to/other.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSettingsGradleFile(tt.filePath); got != tt.want {
				t.Errorf("IsSettingsGradleFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsKotlinDSL(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{"kotlin file", "path/to/build.gradle.kts", true},
		{"groovy file", "path/to/build.gradle", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsKotlinDSL(tt.filePath); got != tt.want {
				t.Errorf("IsKotlinDSL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindGradleFiles(t *testing.T) {
	// Create a temporary directory structure。
	tmpDir, err := os.MkdirTemp("", "gradle-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files。
	testFiles := []string{
		filepath.Join(tmpDir, "build.gradle"),
		filepath.Join(tmpDir, "settings.gradle"),
		filepath.Join(tmpDir, "subdir", "build.gradle.kts"),
	}

	// Create directories。
	if err := os.Mkdir(filepath.Join(tmpDir, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Create files。
	for _, file := range testFiles {
		f, err := os.Create(file)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	// Run the test。
	files, err := FindGradleFiles(tmpDir)
	if err != nil {
		t.Fatalf("FindGradleFiles() error = %v", err)
	}

	// Expected number of files。
	if len(files) != len(testFiles) {
		t.Errorf("FindGradleFiles() found %v files, want %v", len(files), len(testFiles))
	}
}

func TestFindProjectRoot(t *testing.T) {
	// Create a temporary directory structure。
	tmpDir, err := os.MkdirTemp("", "gradle-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create project structure。
	projectRoot := filepath.Join(tmpDir, "project")
	subDir := filepath.Join(projectRoot, "subdir")

	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create build.gradle in project root。
	buildFile := filepath.Join(projectRoot, "build.gradle")
	f, err := os.Create(buildFile)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Test finding project root from subdirectory。
	root, err := FindProjectRoot(subDir)
	if err != nil {
		t.Fatalf("FindProjectRoot() error = %v", err)
	}

	if root != projectRoot {
		t.Errorf("FindProjectRoot() = %v, want %v", root, projectRoot)
	}

	// Test finding project root when not in a project。
	_, err = FindProjectRoot(filepath.Join(tmpDir, "not-a-project"))
	if err == nil {
		t.Error("FindProjectRoot() should return error when not in a project")
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file。
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Create a temporary directory。
	tmpDir, err := os.MkdirTemp("", "test-dir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{"existing file", tmpFile.Name(), true},
		{"directory", tmpDir, false},
		{"non-existent file", "/path/to/nonexistent/file", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExists(tt.filePath); got != tt.want {
				t.Errorf("fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFileContent(t *testing.T) {
	// Create a temporary file with content。
	content := "test content"
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.WriteString(content)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Test getting content from existing file。
	got, err := GetFileContent(tmpFile.Name())
	if err != nil {
		t.Fatalf("GetFileContent() error = %v", err)
	}
	if got != content {
		t.Errorf("GetFileContent() = %v, want %v", got, content)
	}

	// Test getting content from non-existent file。
	_, err = GetFileContent("/path/to/nonexistent/file")
	if err == nil {
		t.Error("GetFileContent() should return error for non-existent file")
	}
}
