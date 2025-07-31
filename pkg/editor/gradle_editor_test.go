package editor

import (
	"strings"
	"testing"

	"github.com/scagogogo/gradle-parser/pkg/model"
	"github.com/scagogogo/gradle-parser/pkg/parser"
)

// 测试用的示例Gradle文件。
const testGradleContent = `// Test Gradle file。
plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
}

group = 'com.example'
version = '0.1.0-SNAPSHOT'
description = 'Test project'

sourceCompatibility = '11'
targetCompatibility = '11'

repositories {
    mavenCentral()
    google()
    maven { url 'https://jitpack.io' }
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
    implementation 'mysql:mysql-connector-java:8.0.29'
    implementation 'com.google.guava:guava:31.0-jre'
    
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
    testImplementation 'org.junit.jupiter:junit-jupiter-api:5.8.2'
    testRuntimeOnly 'org.junit.jupiter:junit-jupiter-engine:5.8.2'
}

task customTask {
    group = 'custom'
    description = 'A custom task'
    doLast {
        println 'Hello from custom task'
    }
}
`

func createTestEditor(t *testing.T) *GradleEditor {
	// 创建位置感知解析器。
	sourceAwareParser := parser.NewSourceAwareParser()
	result, err := sourceAwareParser.ParseWithSourceMapping(testGradleContent)
	if err != nil {
		t.Fatalf("Failed to parse test content: %v", err)
	}

	return NewGradleEditor(result.SourceMappedProject)
}

func TestGradleEditor_UpdateDependencyVersion(t *testing.T) {
	tests := []struct {
		name        string
		group       string
		artifact    string
		newVersion  string
		expectError bool
		expectDiff  bool
	}{
		{
			name:        "Update existing dependency with version",
			group:       "mysql",
			artifact:    "mysql-connector-java",
			newVersion:  "8.0.30",
			expectError: false,
			expectDiff:  true,
		},
		{
			name:        "Add version to dependency without version",
			group:       "org.springframework.boot",
			artifact:    "spring-boot-starter-web",
			newVersion:  "2.7.1",
			expectError: false,
			expectDiff:  true,
		},
		{
			name:        "Update to same version (no change)",
			group:       "mysql",
			artifact:    "mysql-connector-java",
			newVersion:  "8.0.29",
			expectError: false,
			expectDiff:  false,
		},
		{
			name:        "Update non-existent dependency",
			group:       "non.existent",
			artifact:    "artifact",
			newVersion:  "1.0.0",
			expectError: true,
			expectDiff:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := createTestEditor(t)

			err := editor.UpdateDependencyVersion(tt.group, tt.artifact, tt.newVersion)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			modifications := editor.GetModifications()
			hasDiff := len(modifications) > 0

			if tt.expectDiff != hasDiff {
				t.Errorf("Expected diff: %v, got: %v (modifications: %d)", tt.expectDiff, hasDiff, len(modifications))
			}

			if hasDiff {
				// 验证修改的内容。
				mod := modifications[0]
				if mod.Type != ModificationTypeReplace {
					t.Errorf("Expected replace modification, got: %s", mod.Type)
				}

				if !strings.Contains(mod.NewText, tt.newVersion) {
					t.Errorf("New text should contain version %s, got: %s", tt.newVersion, mod.NewText)
				}
			}
		})
	}
}

func TestGradleEditor_UpdatePluginVersion(t *testing.T) {
	tests := []struct {
		name        string
		pluginId    string
		newVersion  string
		expectError bool
		expectDiff  bool
	}{
		{
			name:        "Update existing plugin version",
			pluginId:    "org.springframework.boot",
			newVersion:  "2.7.1",
			expectError: false,
			expectDiff:  true,
		},
		{
			name:        "Update to same version (no change)",
			pluginId:    "org.springframework.boot",
			newVersion:  "2.7.0",
			expectError: false,
			expectDiff:  false,
		},
		{
			name:        "Update non-existent plugin",
			pluginId:    "non.existent.plugin",
			newVersion:  "1.0.0",
			expectError: true,
			expectDiff:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := createTestEditor(t)

			err := editor.UpdatePluginVersion(tt.pluginId, tt.newVersion)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			modifications := editor.GetModifications()
			hasDiff := len(modifications) > 0

			if tt.expectDiff != hasDiff {
				t.Errorf("Expected diff: %v, got: %v", tt.expectDiff, hasDiff)
			}
		})
	}
}

func TestGradleEditor_UpdateProperty(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		newValue    string
		expectError bool
		expectDiff  bool
	}{
		{
			name:        "Update existing property",
			key:         "version",
			newValue:    "1.0.0",
			expectError: false,
			expectDiff:  true,
		},
		{
			name:        "Update to same value (no change)",
			key:         "version",
			newValue:    "0.1.0-SNAPSHOT",
			expectError: false,
			expectDiff:  false,
		},
		{
			name:        "Update non-existent property",
			key:         "nonExistentProperty",
			newValue:    "value",
			expectError: true,
			expectDiff:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := createTestEditor(t)

			err := editor.UpdateProperty(tt.key, tt.newValue)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			modifications := editor.GetModifications()
			hasDiff := len(modifications) > 0

			if tt.expectDiff != hasDiff {
				t.Errorf("Expected diff: %v, got: %v", tt.expectDiff, hasDiff)
			}
		})
	}
}

func TestGradleEditor_AddDependency(t *testing.T) {
	tests := []struct {
		name        string
		group       string
		artifact    string
		version     string
		scope       string
		expectError bool
	}{
		{
			name:        "Add new dependency with version",
			group:       "org.apache.commons",
			artifact:    "commons-text",
			version:     "1.9",
			scope:       "implementation",
			expectError: false,
		},
		{
			name:        "Add new dependency without version",
			group:       "org.springframework.boot",
			artifact:    "spring-boot-starter-security",
			version:     "",
			scope:       "implementation",
			expectError: false,
		},
		{
			name:        "Add dependency with default scope",
			group:       "org.apache.commons",
			artifact:    "commons-lang3",
			version:     "3.12.0",
			scope:       "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := createTestEditor(t)

			err := editor.AddDependency(tt.group, tt.artifact, tt.version, tt.scope)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			modifications := editor.GetModifications()
			if len(modifications) == 0 {
				t.Errorf("Expected modification but got none")
				return
			}

			mod := modifications[0]
			if mod.Type != ModificationTypeInsert {
				t.Errorf("Expected insert modification, got: %s", mod.Type)
			}

			expectedScope := tt.scope
			if expectedScope == "" {
				expectedScope = "implementation"
			}

			if !strings.Contains(mod.NewText, expectedScope) {
				t.Errorf("New text should contain scope %s, got: %s", expectedScope, mod.NewText)
			}

			if !strings.Contains(mod.NewText, tt.group) {
				t.Errorf("New text should contain group %s, got: %s", tt.group, mod.NewText)
			}

			if !strings.Contains(mod.NewText, tt.artifact) {
				t.Errorf("New text should contain artifact %s, got: %s", tt.artifact, mod.NewText)
			}
		})
	}
}

// 测试编辑器的边界条件和错误处理。
func TestGradleEditorEdgeCases(t *testing.T) {
	t.Run("Empty project", func(t *testing.T) {
		// 创建空的源码映射项目。
		emptyProject := &model.SourceMappedProject{
			SourceMappedDependencies: []*model.SourceMappedDependency{},
			SourceMappedPlugins:      []*model.SourceMappedPlugin{},
			SourceMappedProperties:   []*model.SourceMappedProperty{},
			SourceMappedRepositories: []*model.SourceMappedRepository{},
			Lines:                    []string{},
			OriginalText:             "",
		}

		editor := NewGradleEditor(emptyProject)

		// 尝试更新不存在的依赖。
		err := editor.UpdateDependencyVersion("group", "name", "1.0.0")
		if err == nil {
			t.Error("Should return error for non-existent dependency in empty project")
		}

		// 尝试更新不存在的插件。
		err = editor.UpdatePluginVersion("plugin", "1.0.0")
		if err == nil {
			t.Error("Should return error for non-existent plugin in empty project")
		}

		// 尝试更新不存在的属性。
		err = editor.UpdateProperty("property", "value")
		if err == nil {
			t.Error("Should return error for non-existent property in empty project")
		}

		// 尝试添加依赖到不存在的dependencies块。
		err = editor.AddDependency("group", "name", "1.0.0", "implementation")
		if err == nil {
			t.Error("Should return error when dependencies block not found")
		}
	})

	t.Run("Nil project", func(t *testing.T) {
		// 测试nil项目的处理。
		editor := NewGradleEditor(nil)

		// 所有操作都应该返回错误或安全处理。
		err := editor.UpdateDependencyVersion("group", "name", "1.0.0")
		if err == nil {
			t.Error("Should handle nil project gracefully")
		}
	})

	t.Run("Duplicate dependencies", func(t *testing.T) {
		// 创建包含重复依赖的项目。
		duplicateProject := &model.SourceMappedProject{
			SourceMappedDependencies: []*model.SourceMappedDependency{
				{
					Dependency: &model.Dependency{
						Group:   "mysql",
						Name:    "mysql-connector-java",
						Version: "8.0.29",
						Scope:   "implementation",
					},
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{Line: 2, Column: 5},
						End:   model.SourcePosition{Line: 2, Column: 50},
					},
					RawText: "implementation 'mysql:mysql-connector-java:8.0.29'",
				},
				{
					Dependency: &model.Dependency{
						Group:   "mysql",
						Name:    "mysql-connector-java",
						Version: "8.0.28",
						Scope:   "testImplementation",
					},
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{Line: 3, Column: 5},
						End:   model.SourcePosition{Line: 3, Column: 55},
					},
					RawText: "testImplementation 'mysql:mysql-connector-java:8.0.28'",
				},
			},
			Lines: []string{
				"dependencies {",
				"    implementation 'mysql:mysql-connector-java:8.0.29'",
				"    testImplementation 'mysql:mysql-connector-java:8.0.28'",
				"}",
			},
			OriginalText: "dependencies {\n    implementation 'mysql:mysql-connector-java:8.0.29'\n    testImplementation 'mysql:mysql-connector-java:8.0.28'\n}",
		}

		editor := NewGradleEditor(duplicateProject)

		// 更新第一个匹配的依赖。
		err := editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.30")
		if err != nil {
			t.Fatalf("Should be able to update first matching dependency: %v", err)
		}

		modifications := editor.GetModifications()
		if len(modifications) != 1 {
			t.Errorf("Should create exactly 1 modification for duplicate dependencies, got %d", len(modifications))
		}
	})

	t.Run("Very long values", func(t *testing.T) {
		editor := createTestEditor(t)

		// 尝试设置非常长的版本号。
		longVersion := strings.Repeat("1.0.0-", 100) + "SNAPSHOT"
		err := editor.UpdateDependencyVersion("mysql", "mysql-connector-java", longVersion)
		if err != nil {
			t.Fatalf("Should handle long version values: %v", err)
		}

		modifications := editor.GetModifications()
		if len(modifications) != 1 {
			t.Errorf("Expected 1 modification for long version, got %d", len(modifications))
		}

		if !strings.Contains(modifications[0].NewText, longVersion) {
			t.Error("Modification should contain the long version")
		}
	})

	t.Run("Special characters in values", func(t *testing.T) {
		editor := createTestEditor(t)

		// 测试包含特殊字符的版本号。
		specialVersion := "1.0.0-测试版本-🚀-ñ"
		err := editor.UpdateDependencyVersion("mysql", "mysql-connector-java", specialVersion)
		if err != nil {
			t.Fatalf("Should handle special characters in version: %v", err)
		}

		modifications := editor.GetModifications()
		if len(modifications) != 1 {
			t.Errorf("Expected 1 modification for special version, got %d", len(modifications))
		}

		if !strings.Contains(modifications[0].NewText, specialVersion) {
			t.Error("Modification should contain the special version")
		}
	})

	t.Run("Empty values", func(t *testing.T) {
		editor := createTestEditor(t)

		// 测试空版本号。
		err := editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "")
		if err != nil {
			t.Fatalf("Should handle empty version: %v", err)
		}

		modifications := editor.GetModifications()
		if len(modifications) != 1 {
			t.Errorf("Expected 1 modification for empty version, got %d", len(modifications))
		}
	})

	t.Run("Case sensitivity", func(t *testing.T) {
		editor := createTestEditor(t)

		// 测试大小写敏感性。
		err := editor.UpdateDependencyVersion("MYSQL", "mysql-connector-java", "8.0.30")
		if err == nil {
			t.Error("Should be case sensitive for group names")
		}

		err = editor.UpdateDependencyVersion("mysql", "MYSQL-CONNECTOR-JAVA", "8.0.30")
		if err == nil {
			t.Error("Should be case sensitive for artifact names")
		}
	})
}
