package dependency

import (
	"reflect"
	"testing"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

func TestNewDependencyParser(t *testing.T) {
	parser := NewDependencyParser()
	if parser == nil {
		t.Error("NewDependencyParser() returned nil")
	}
}

func TestParseDependencyBlock(t *testing.T) {
	parser := NewDependencyParser()

	// Test with nil block
	_, err := parser.ParseDependencyBlock(nil)
	if err == nil {
		t.Error("ParseDependencyBlock() should return error for nil block")
	}

	// Test with dependencies in standard scopes
	block := &model.ScriptBlock{
		Name: "dependencies",
		Closures: map[string][]*model.ScriptBlock{
			"implementation": {
				{
					Values: map[string]interface{}{
						"'org.springframework:spring-core:5.3.10'": "'org.springframework:spring-core:5.3.10'",
					},
				},
			},
			"testImplementation": {
				{
					Values: map[string]interface{}{
						"'junit:junit:4.13.2'": "'junit:junit:4.13.2'",
					},
				},
			},
		},
	}

	deps, err := parser.ParseDependencyBlock(block)
	if err != nil {
		t.Fatalf("ParseDependencyBlock() error = %v", err)
	}

	if len(deps) != 2 {
		t.Errorf("ParseDependencyBlock() returned %v dependencies, want 2", len(deps))
	}

	// Verify the dependencies
	var foundSpring, foundJunit bool
	for _, dep := range deps {
		if dep.Group == "org.springframework" && dep.Name == "spring-core" && dep.Version == "5.3.10" && dep.Scope == "implementation" {
			foundSpring = true
		}
		if dep.Group == "junit" && dep.Name == "junit" && dep.Version == "4.13.2" && dep.Scope == "testImplementation" {
			foundJunit = true
		}
	}

	if !foundSpring {
		t.Error("ParseDependencyBlock() did not find Spring dependency")
	}
	if !foundJunit {
		t.Error("ParseDependencyBlock() did not find JUnit dependency")
	}

	// Test with custom scope
	customBlock := &model.ScriptBlock{
		Name: "dependencies",
		Closures: map[string][]*model.ScriptBlock{
			"customScope": {
				{
					Values: map[string]interface{}{
						"'com.google.guava:guava:31.1-jre'": "'com.google.guava:guava:31.1-jre'",
					},
				},
			},
		},
	}

	deps, err = parser.ParseDependencyBlock(customBlock)
	if err != nil {
		t.Fatalf("ParseDependencyBlock() with custom scope error = %v", err)
	}

	if len(deps) != 1 {
		t.Errorf("ParseDependencyBlock() with custom scope returned %v dependencies, want 1", len(deps))
	}

	// Verify the custom scope dependency
	if deps[0].Group != "com.google.guava" || deps[0].Name != "guava" || deps[0].Version != "31.1-jre" || deps[0].Scope != "customScope" {
		t.Errorf("Custom scope dependency not parsed correctly: %+v", deps[0])
	}
}

func TestParseScopedDependencies(t *testing.T) {
	parser := NewDependencyParser()

	// Create a block with implementation dependencies
	block := &model.ScriptBlock{
		Name: "dependencies",
		Closures: map[string][]*model.ScriptBlock{
			"implementation": {
				{
					Values: map[string]interface{}{
						"'org.springframework:spring-core:5.3.10'": "'org.springframework:spring-core:5.3.10'",
						"'com.google.guava:guava:31.1-jre'":        "'com.google.guava:guava:31.1-jre'",
					},
				},
			},
		},
	}

	deps, err := parser.parseScopedDependencies(block, "implementation")
	if err != nil {
		t.Fatalf("parseScopedDependencies() error = %v", err)
	}

	if len(deps) != 2 {
		t.Errorf("parseScopedDependencies() returned %v dependencies, want 2", len(deps))
	}

	// Verify scope is set correctly
	for _, dep := range deps {
		if dep.Scope != "implementation" {
			t.Errorf("Dependency scope incorrect, got %s, want implementation", dep.Scope)
		}
	}

	// Test with non-existent scope
	deps, err = parser.parseScopedDependencies(block, "nonexistentScope")
	if err != nil {
		t.Fatalf("parseScopedDependencies() with non-existent scope error = %v", err)
	}

	if len(deps) != 0 {
		t.Errorf("parseScopedDependencies() with non-existent scope returned %v dependencies, want 0", len(deps))
	}
}

func TestParseCustomDependencies(t *testing.T) {
	parser := NewDependencyParser()

	// Create a block with custom dependencies
	block := &model.ScriptBlock{
		Values: map[string]interface{}{
			"'org.springframework:spring-core:5.3.10'": "'org.springframework:spring-core:5.3.10'",
			"project(':app')":                          "project(':app')",
		},
	}

	deps, err := parser.parseCustomDependencies(block, "customScope")
	if err != nil {
		t.Fatalf("parseCustomDependencies() error = %v", err)
	}

	if len(deps) != 2 {
		t.Errorf("parseCustomDependencies() returned %v dependencies, want 2", len(deps))
	}

	// Verify scope is set correctly and dependencies are parsed correctly
	var foundSpring, foundProject bool
	for _, dep := range deps {
		if dep.Scope != "customScope" {
			t.Errorf("Dependency scope incorrect, got %s, want customScope", dep.Scope)
		}

		if dep.Group == "org.springframework" && dep.Name == "spring-core" && dep.Version == "5.3.10" {
			foundSpring = true
		}
		if dep.Name == "app" && dep.Group == "" && dep.Version == "" && dep.Raw == "project(':app')" {
			foundProject = true
		}
	}

	if !foundSpring {
		t.Error("parseCustomDependencies() did not find Spring dependency")
	}
	if !foundProject {
		t.Error("parseCustomDependencies() did not find project dependency")
	}
}

func TestParseDependencyString(t *testing.T) {
	parser := NewDependencyParser()

	tests := []struct {
		name    string
		depStr  string
		scope   string
		want    *model.Dependency
		success bool
	}{
		{
			name:   "GAV format",
			depStr: "org.springframework:spring-core:5.3.10",
			scope:  "implementation",
			want: &model.Dependency{
				Group:   "org.springframework",
				Name:    "spring-core",
				Version: "5.3.10",
				Scope:   "implementation",
				Raw:     "org.springframework:spring-core:5.3.10",
			},
			success: true,
		},
		{
			name:   "quoted GAV format",
			depStr: "'org.springframework:spring-core:5.3.10'",
			scope:  "implementation",
			want: &model.Dependency{
				Group:   "org.springframework",
				Name:    "spring-core",
				Version: "5.3.10",
				Scope:   "implementation",
				Raw:     "'org.springframework:spring-core:5.3.10'",
			},
			success: true,
		},
		{
			name:   "dot name format",
			depStr: "org.springframework.boot:spring-boot-starter:2.5.5",
			scope:  "api",
			want: &model.Dependency{
				Group:   "org.springframework.boot",
				Name:    "spring-boot-starter",
				Version: "2.5.5",
				Scope:   "api",
				Raw:     "org.springframework.boot:spring-boot-starter:2.5.5",
			},
			success: true,
		},
		{
			name:   "project reference",
			depStr: "project(':app')",
			scope:  "implementation",
			want: &model.Dependency{
				Group:   "",
				Name:    "app",
				Version: "",
				Scope:   "implementation",
				Raw:     "project(':app')",
			},
			success: true,
		},
		{
			name:    "invalid format",
			depStr:  "invalid-dependency-format",
			scope:   "implementation",
			want:    nil,
			success: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parser.parseDependencyString(tt.depStr, tt.scope)
			if ok != tt.success {
				t.Errorf("parseDependencyString() success = %v, want %v", ok, tt.success)
				return
			}

			if tt.success && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDependencyString() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestExtractDependenciesFromText(t *testing.T) {
	parser := NewDependencyParser()

	// Test with empty text
	deps := parser.ExtractDependenciesFromText("")
	if len(deps) != 0 {
		t.Errorf("ExtractDependenciesFromText() with empty text returned %v dependencies, want 0", len(deps))
	}

	// Test with dependencies in text
	text := `dependencies {
		implementation 'org.springframework:spring-core:5.3.10'
		testImplementation 'junit:junit:4.13.2'
		api project(':app')
	}`

	deps = parser.ExtractDependenciesFromText(text)
	if len(deps) < 3 {
		t.Errorf("ExtractDependenciesFromText() returned %v dependencies, want at least 3", len(deps))
	}

	// Verify extraction of specific dependency types
	var foundSpring, foundJunit, foundProject bool
	for _, dep := range deps {
		if dep.Group == "org.springframework" && dep.Name == "spring-core" && dep.Version == "5.3.10" {
			foundSpring = true
		}
		if dep.Group == "junit" && dep.Name == "junit" && dep.Version == "4.13.2" {
			foundJunit = true
		}
		if dep.Name == "app" && dep.Raw == "project(':app')" {
			foundProject = true
		}
	}

	if !foundSpring {
		t.Error("ExtractDependenciesFromText() did not find Spring dependency")
	}
	if !foundJunit {
		t.Error("ExtractDependenciesFromText() did not find JUnit dependency")
	}
	if !foundProject {
		t.Error("ExtractDependenciesFromText() did not find project dependency")
	}
}

func TestGroupDependenciesByScope(t *testing.T) {
	parser := NewDependencyParser()

	// Create a list of dependencies with different scopes
	deps := []*model.Dependency{
		{Group: "org.springframework", Name: "spring-core", Version: "5.3.10", Scope: "implementation"},
		{Group: "junit", Name: "junit", Version: "4.13.2", Scope: "testImplementation"},
		{Group: "com.google.guava", Name: "guava", Version: "31.1-jre", Scope: "implementation"},
		{Group: "org.mockito", Name: "mockito-core", Version: "4.0.0", Scope: "testImplementation"},
		{Name: "app", Scope: "implementation", Raw: "project(':app')"},
		{Group: "no.scope", Name: "dep", Version: "1.0"}, // No scope
	}

	sets := parser.GroupDependenciesByScope(deps)

	// Verify we have 2 sets: implementation and testImplementation
	if len(sets) != 2 {
		t.Errorf("GroupDependenciesByScope() returned %v sets, want 2", len(sets))
		return
	}

	// Verify each scope's dependency count
	for _, set := range sets {
		switch set.Scope {
		case "implementation":
			// Should have 4 dependencies: 3 explicit + 1 default no-scope
			if len(set.Dependencies) != 4 {
				t.Errorf("Implementation scope set has %v dependencies, want 4", len(set.Dependencies))
			}
		case "testImplementation":
			// Should have 2 dependencies
			if len(set.Dependencies) != 2 {
				t.Errorf("TestImplementation scope set has %v dependencies, want 2", len(set.Dependencies))
			}
		default:
			t.Errorf("Unexpected scope found: %s", set.Scope)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		str   string
		want  bool
	}{
		{
			name:  "empty slice",
			slice: []string{},
			str:   "test",
			want:  false,
		},
		{
			name:  "string not in slice",
			slice: []string{"one", "two", "three"},
			str:   "four",
			want:  false,
		},
		{
			name:  "string in slice",
			slice: []string{"one", "two", "three"},
			str:   "two",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.slice, tt.str); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
