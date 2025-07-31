package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestProject(t *testing.T) {
	// Test that we can create and use a Project。
	project := &Project{
		Group:       "com.example",
		Name:        "test-project",
		Version:     "1.0.0",
		Description: "Test project",
		Properties: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		Plugins:      make([]*Plugin, 0),
		Dependencies: make([]*Dependency, 0),
		Repositories: make([]*Repository, 0),
		SubProjects:  make([]*Project, 0),
		Tasks:        make([]*Task, 0),
		Extensions:   make(map[string]any),
		FilePath:     "/path/to/build.gradle",
	}

	// Verify the project fields。
	if project.Group != "com.example" {
		t.Errorf("Project Group = %s, want com.example", project.Group)
	}
	if project.Name != "test-project" {
		t.Errorf("Project Name = %s, want test-project", project.Name)
	}
	if len(project.Properties) != 2 {
		t.Errorf("Project Properties has %d items, want 2", len(project.Properties))
	}
}

func TestDependency(t *testing.T) {
	// Test that we can create and use a Dependency。
	dep := &Dependency{
		Group:      "org.springframework",
		Name:       "spring-core",
		Version:    "5.3.10",
		Scope:      "implementation",
		Transitive: true,
		Raw:        "org.springframework:spring-core:5.3.10",
	}

	// Verify the dependency fields。
	if dep.Group != "org.springframework" {
		t.Errorf("Dependency Group = %s, want org.springframework", dep.Group)
	}
	if dep.Name != "spring-core" {
		t.Errorf("Dependency Name = %s, want spring-core", dep.Name)
	}
	if dep.Scope != "implementation" {
		t.Errorf("Dependency Scope = %s, want implementation", dep.Scope)
	}
}

func TestPlugin(t *testing.T) {
	// Test that we can create and use a Plugin。
	plugin := &Plugin{
		ID:      "java",
		Version: "1.0.0",
		Apply:   true,
		Config:  make(map[string]interface{}),
	}

	// Verify the plugin fields。
	if plugin.ID != "java" {
		t.Errorf("Plugin ID = %s, want java", plugin.ID)
	}
	if plugin.Version != "1.0.0" {
		t.Errorf("Plugin Version = %s, want 1.0.0", plugin.Version)
	}
	if !plugin.Apply {
		t.Error("Plugin Apply = false, want true")
	}
}

func TestRepository(t *testing.T) {
	// Test that we can create and use a Repository。
	repo := &Repository{
		Name:     "mavenCentral",
		URL:      "https://repo.maven.apache.org/maven2/",
		Type:     "maven",
		Config:   make(map[string]interface{}),
		Username: "user",
		Password: "pass",
	}

	// Verify the repository fields。
	if repo.Name != "mavenCentral" {
		t.Errorf("Repository Name = %s, want mavenCentral", repo.Name)
	}
	if repo.URL != "https://repo.maven.apache.org/maven2/" {
		t.Errorf("Repository URL = %s, want https://repo.maven.apache.org/maven2/", repo.URL)
	}
	if repo.Type != "maven" {
		t.Errorf("Repository Type = %s, want maven", repo.Type)
	}
}

func TestTask(t *testing.T) {
	// Test that we can create and use a Task。
	task := &Task{
		Name:        "test",
		Type:        "Test",
		Description: "Run tests",
		Group:       "verification",
		DependsOn:   []string{"compile"},
		Config:      make(map[string]interface{}),
	}

	// Verify the task fields。
	if task.Name != "test" {
		t.Errorf("Task Name = %s, want test", task.Name)
	}
	if task.Type != "Test" {
		t.Errorf("Task Type = %s, want Test", task.Type)
	}
	if len(task.DependsOn) != 1 || task.DependsOn[0] != "compile" {
		t.Errorf("Task DependsOn = %v, want [compile]", task.DependsOn)
	}
}

func TestScriptBlock(t *testing.T) {
	// Test that we can create and use a ScriptBlock。
	parent := &ScriptBlock{
		Name:     "parent",
		Children: make([]*ScriptBlock, 0),
		Values:   make(map[string]interface{}),
		Closures: make(map[string][]*ScriptBlock),
	}

	child := &ScriptBlock{
		Name:     "child",
		Parent:   parent,
		Children: make([]*ScriptBlock, 0),
		Values: map[string]interface{}{
			"key": "value",
		},
		Closures: make(map[string][]*ScriptBlock),
	}

	parent.Children = append(parent.Children, child)

	// Verify the relationship。
	if len(parent.Children) != 1 {
		t.Errorf("Parent has %d children, want 1", len(parent.Children))
	}
	if parent.Children[0] != child {
		t.Error("Parent's child is not correct")
	}
	if child.Parent != parent {
		t.Error("Child's parent is not correct")
	}
}

func TestDependencySet(t *testing.T) {
	// Test that we can create and use a DependencySet。
	deps := []*Dependency{
		{Group: "org.springframework", Name: "spring-core", Version: "5.3.10"},
		{Group: "com.google.guava", Name: "guava", Version: "31.1-jre"},
	}

	set := &DependencySet{
		Scope:        "implementation",
		Dependencies: deps,
	}

	// Verify the dependency set。
	if set.Scope != "implementation" {
		t.Errorf("DependencySet Scope = %s, want implementation", set.Scope)
	}
	if len(set.Dependencies) != 2 {
		t.Errorf("DependencySet has %d dependencies, want 2", len(set.Dependencies))
	}
}

func TestParseResult(t *testing.T) {
	// Test that we can create and use a ParseResult。
	project := &Project{
		Group:   "com.example",
		Name:    "test-project",
		Version: "1.0.0",
	}

	result := &ParseResult{
		Project:   project,
		RawText:   "sample raw text",
		Errors:    []error{nil},
		Warnings:  []string{"warning1", "warning2"},
		ParseTime: "100ms",
	}

	// Verify the parse result。
	if result.Project != project {
		t.Error("ParseResult Project is not correct")
	}
	if result.RawText != "sample raw text" {
		t.Errorf("ParseResult RawText = %s, want sample raw text", result.RawText)
	}
	if len(result.Warnings) != 2 {
		t.Errorf("ParseResult has %d warnings, want 2", len(result.Warnings))
	}
	if result.ParseTime != "100ms" {
		t.Errorf("ParseResult ParseTime = %s, want 100ms", result.ParseTime)
	}
}

func TestModelJSON(t *testing.T) {
	// Test that models can be marshalled to JSON。
	project := &Project{
		Group:       "com.example",
		Name:        "test-project",
		Version:     "1.0.0",
		Description: "Test project",
		Plugins: []*Plugin{
			{ID: "java", Version: "1.0.0"},
		},
		Dependencies: []*Dependency{
			{Group: "org.springframework", Name: "spring-core", Version: "5.3.10", Scope: "implementation"},
		},
		Repositories: []*Repository{
			{Name: "mavenCentral", URL: "https://repo.maven.apache.org/maven2/", Type: "maven"},
		},
	}

	// Marshal to JSON。
	jsonData, err := json.Marshal(project)
	if err != nil {
		t.Fatalf("Failed to marshal Project to JSON: %v", err)
	}

	// Verify that the JSON contains expected fields。
	jsonStr := string(jsonData)
	expectedFields := []string{"group", "name", "version", "description", "plugins", "dependencies", "repositories"}
	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, "\""+field+"\"") {
			t.Errorf("JSON does not contain field %s", field)
		}
	}
}

// Helper function to check if a string contains a substring。
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
