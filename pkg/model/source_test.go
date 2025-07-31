package model

import (
	"testing"
)

func TestSourcePosition(t *testing.T) {
	pos := SourcePosition{
		Line:     10,
		Column:   5,
		StartPos: 100,
		EndPos:   110,
		Length:   10,
	}

	expected := "line 10, col 5"
	if pos.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, pos.String())
	}
}

func TestSourceRange(t *testing.T) {
	start := SourcePosition{Line: 10, Column: 5}
	end := SourcePosition{Line: 10, Column: 15}

	sourceRange := SourceRange{
		Start: start,
		End:   end,
	}

	expected := "line 10, col 5 - line 10, col 15"
	if sourceRange.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, sourceRange.String())
	}
}

func TestSourceMappedProject_GetLineText(t *testing.T) {
	lines := []string{
		"line 1",
		"line 2",
		"line 3",
	}

	project := &SourceMappedProject{
		Lines: lines,
	}

	tests := []struct {
		lineNumber int
		expected   string
	}{
		{1, "line 1"},
		{2, "line 2"},
		{3, "line 3"},
		{0, ""},  // invalid line number。
		{4, ""},  // line number out of range。
		{-1, ""}, // negative line number。
	}

	for _, tt := range tests {
		result := project.GetLineText(tt.lineNumber)
		if result != tt.expected {
			t.Errorf("GetLineText(%d): expected '%s', got '%s'", tt.lineNumber, tt.expected, result)
		}
	}
}

func TestSourceMappedProject_GetTextRange(t *testing.T) {
	originalText := "Hello, World!\nThis is line 2\nAnd line 3"

	project := &SourceMappedProject{
		OriginalText: originalText,
	}

	tests := []struct {
		name        string
		sourceRange SourceRange
		expected    string
	}{
		{
			name: "Valid range",
			sourceRange: SourceRange{
				Start: SourcePosition{StartPos: 0},
				End:   SourcePosition{EndPos: 5},
			},
			expected: "Hello",
		},
		{
			name: "Full first line",
			sourceRange: SourceRange{
				Start: SourcePosition{StartPos: 0},
				End:   SourcePosition{EndPos: 13},
			},
			expected: "Hello, World!",
		},
		{
			name: "Invalid range - negative start",
			sourceRange: SourceRange{
				Start: SourcePosition{StartPos: -1},
				End:   SourcePosition{EndPos: 5},
			},
			expected: "",
		},
		{
			name: "Invalid range - end exceeds text length",
			sourceRange: SourceRange{
				Start: SourcePosition{StartPos: 0},
				End:   SourcePosition{EndPos: 1000},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := project.GetTextRange(tt.sourceRange)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSourceMappedProject_FindDependencyByPosition(t *testing.T) {
	// 创建测试依赖。
	dep1 := &SourceMappedDependency{
		Dependency: &Dependency{Group: "com.example", Name: "lib1"},
		SourceRange: SourceRange{
			Start: SourcePosition{Line: 5, Column: 10},
			End:   SourcePosition{Line: 5, Column: 30},
		},
	}

	dep2 := &SourceMappedDependency{
		Dependency: &Dependency{Group: "org.test", Name: "lib2"},
		SourceRange: SourceRange{
			Start: SourcePosition{Line: 8, Column: 5},
			End:   SourcePosition{Line: 8, Column: 25},
		},
	}

	project := &SourceMappedProject{
		SourceMappedDependencies: []*SourceMappedDependency{dep1, dep2},
	}

	tests := []struct {
		name     string
		line     int
		column   int
		expected *SourceMappedDependency
	}{
		{
			name:     "Find first dependency",
			line:     5,
			column:   15,
			expected: dep1,
		},
		{
			name:     "Find second dependency",
			line:     8,
			column:   10,
			expected: dep2,
		},
		{
			name:     "Position before first dependency",
			line:     5,
			column:   5,
			expected: nil,
		},
		{
			name:     "Position after first dependency",
			line:     5,
			column:   35,
			expected: nil,
		},
		{
			name:     "Position on different line",
			line:     6,
			column:   15,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := project.FindDependencyByPosition(tt.line, tt.column)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSourceMappedProject_FindPluginByPosition(t *testing.T) {
	// 创建测试插件。
	plugin1 := &SourceMappedPlugin{
		Plugin: &Plugin{ID: "java"},
		SourceRange: SourceRange{
			Start: SourcePosition{Line: 2, Column: 5},
			End:   SourcePosition{Line: 2, Column: 15},
		},
	}

	plugin2 := &SourceMappedPlugin{
		Plugin: &Plugin{ID: "org.springframework.boot"},
		SourceRange: SourceRange{
			Start: SourcePosition{Line: 3, Column: 5},
			End:   SourcePosition{Line: 3, Column: 45},
		},
	}

	project := &SourceMappedProject{
		SourceMappedPlugins: []*SourceMappedPlugin{plugin1, plugin2},
	}

	tests := []struct {
		name     string
		line     int
		column   int
		expected *SourceMappedPlugin
	}{
		{
			name:     "Find first plugin",
			line:     2,
			column:   10,
			expected: plugin1,
		},
		{
			name:     "Find second plugin",
			line:     3,
			column:   20,
			expected: plugin2,
		},
		{
			name:     "Position not in any plugin",
			line:     4,
			column:   10,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := project.FindPluginByPosition(tt.line, tt.column)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSourceMappedProject_FindPropertyByKey(t *testing.T) {
	// 创建测试属性。
	prop1 := &SourceMappedProperty{
		Key:   "group",
		Value: "com.example",
	}

	prop2 := &SourceMappedProperty{
		Key:   "version",
		Value: "1.0.0",
	}

	project := &SourceMappedProject{
		SourceMappedProperties: []*SourceMappedProperty{prop1, prop2},
	}

	tests := []struct {
		name     string
		key      string
		expected *SourceMappedProperty
	}{
		{
			name:     "Find group property",
			key:      "group",
			expected: prop1,
		},
		{
			name:     "Find version property",
			key:      "version",
			expected: prop2,
		},
		{
			name:     "Property not found",
			key:      "description",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := project.FindPropertyByKey(tt.key)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSourceMappedDependency(t *testing.T) {
	dep := &Dependency{
		Group:   "com.example",
		Name:    "test-lib",
		Version: "1.0.0",
		Scope:   "implementation",
	}

	sourceRange := SourceRange{
		Start: SourcePosition{Line: 10, Column: 5, StartPos: 100},
		End:   SourcePosition{Line: 10, Column: 35, StartPos: 130},
	}

	sourceMappedDep := &SourceMappedDependency{
		Dependency:  dep,
		SourceRange: sourceRange,
		RawText:     "implementation 'com.example:test-lib:1.0.0'",
	}

	// 验证依赖信息。
	if sourceMappedDep.Group != "com.example" {
		t.Errorf("Expected group 'com.example', got '%s'", sourceMappedDep.Group)
	}

	if sourceMappedDep.Name != "test-lib" {
		t.Errorf("Expected name 'test-lib', got '%s'", sourceMappedDep.Name)
	}

	if sourceMappedDep.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", sourceMappedDep.Version)
	}

	// 验证源码位置信息。
	if sourceMappedDep.SourceRange.Start.Line != 10 {
		t.Errorf("Expected start line 10, got %d", sourceMappedDep.SourceRange.Start.Line)
	}

	if sourceMappedDep.RawText != "implementation 'com.example:test-lib:1.0.0'" {
		t.Errorf("Expected raw text to match, got '%s'", sourceMappedDep.RawText)
	}
}

func TestSourceMappedPlugin(t *testing.T) {
	plugin := &Plugin{
		ID:      "org.springframework.boot",
		Version: "2.7.0",
		Apply:   true,
	}

	sourceRange := SourceRange{
		Start: SourcePosition{Line: 3, Column: 5, StartPos: 50},
		End:   SourcePosition{Line: 3, Column: 45, StartPos: 90},
	}

	sourceMappedPlugin := &SourceMappedPlugin{
		Plugin:      plugin,
		SourceRange: sourceRange,
		RawText:     "id 'org.springframework.boot' version '2.7.0'",
	}

	// 验证插件信息。
	if sourceMappedPlugin.ID != "org.springframework.boot" {
		t.Errorf("Expected ID 'org.springframework.boot', got '%s'", sourceMappedPlugin.ID)
	}

	if sourceMappedPlugin.Version != "2.7.0" {
		t.Errorf("Expected version '2.7.0', got '%s'", sourceMappedPlugin.Version)
	}

	// 验证源码位置信息。
	if sourceMappedPlugin.SourceRange.Start.Line != 3 {
		t.Errorf("Expected start line 3, got %d", sourceMappedPlugin.SourceRange.Start.Line)
	}
}
