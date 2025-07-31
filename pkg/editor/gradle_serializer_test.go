package editor

import (
	"strings"
	"testing"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

const testSerializerContent = `plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
}

group = 'com.example'
version = '0.1.0-SNAPSHOT'

dependencies {
    implementation 'mysql:mysql-connector-java:8.0.29'
    implementation 'com.google.guava:guava:31.0-jre'
}
`

func TestGradleSerializer_ApplyModifications(t *testing.T) {
	tests := []struct {
		name          string
		modifications []Modification
		expectError   bool
		validateFunc  func(t *testing.T, original, result string)
	}{
		{
			name: "Single replace modification",
			modifications: []Modification{
				{
					Type: ModificationTypeReplace,
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{StartPos: strings.Index(testSerializerContent, "version = '0.1.0-SNAPSHOT'")},
						End:   model.SourcePosition{StartPos: strings.Index(testSerializerContent, "version = '0.1.0-SNAPSHOT'") + len("version = '0.1.0-SNAPSHOT'")},
					},
					OldText: "version = '0.1.0-SNAPSHOT'",
					NewText: "version = '1.0.0'",
				},
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, result string) {
				if !strings.Contains(result, "version = '1.0.0'") {
					t.Errorf("Result should contain updated version")
				}
				if strings.Contains(result, "0.1.0-SNAPSHOT") {
					t.Errorf("Result should not contain old version")
				}
				// 验证其他内容保持不变。
				if !strings.Contains(result, "id 'java'") {
					t.Errorf("Other content should remain unchanged")
				}
			},
		},
		{
			name: "Multiple replace modifications",
			modifications: []Modification{
				{
					Type: ModificationTypeReplace,
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{StartPos: strings.Index(testSerializerContent, "version = '0.1.0-SNAPSHOT'")},
						End:   model.SourcePosition{StartPos: strings.Index(testSerializerContent, "version = '0.1.0-SNAPSHOT'") + len("version = '0.1.0-SNAPSHOT'")},
					},
					OldText: "version = '0.1.0-SNAPSHOT'",
					NewText: "version = '1.0.0'",
				},
				{
					Type: ModificationTypeReplace,
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{StartPos: strings.Index(testSerializerContent, "group = 'com.example'")},
						End:   model.SourcePosition{StartPos: strings.Index(testSerializerContent, "group = 'com.example'") + len("group = 'com.example'")},
					},
					OldText: "group = 'com.example'",
					NewText: "group = 'org.example'",
				},
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, result string) {
				if !strings.Contains(result, "version = '1.0.0'") {
					t.Errorf("Result should contain updated version")
				}
				if !strings.Contains(result, "group = 'org.example'") {
					t.Errorf("Result should contain updated group")
				}
			},
		},
		{
			name: "Insert modification",
			modifications: []Modification{
				{
					Type: ModificationTypeInsert,
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{StartPos: len(testSerializerContent) - 1},
						End:   model.SourcePosition{StartPos: len(testSerializerContent) - 1},
					},
					OldText: "",
					NewText: "    implementation 'org.apache.commons:commons-text:1.9'\n",
				},
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, result string) {
				if !strings.Contains(result, "commons-text") {
					t.Errorf("Result should contain inserted dependency")
				}
				// 验证原始内容仍然存在。
				if !strings.Contains(result, "mysql:mysql-connector-java") {
					t.Errorf("Original content should be preserved")
				}
			},
		},
		{
			name: "Delete modification",
			modifications: []Modification{
				{
					Type: ModificationTypeDelete,
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{StartPos: strings.Index(testSerializerContent, "version = '0.1.0-SNAPSHOT'")},
						End:   model.SourcePosition{StartPos: strings.Index(testSerializerContent, "version = '0.1.0-SNAPSHOT'") + len("version = '0.1.0-SNAPSHOT'")},
					},
					OldText: "version = '0.1.0-SNAPSHOT'",
					NewText: "",
				},
			},
			expectError: false,
			validateFunc: func(t *testing.T, original, result string) {
				if strings.Contains(result, "version = '0.1.0-SNAPSHOT'") {
					t.Errorf("Deleted content should not be present")
				}
				// 验证其他内容仍然存在。
				if !strings.Contains(result, "group = 'com.example'") {
					t.Errorf("Other content should be preserved")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serializer := NewGradleSerializer(testSerializerContent)

			result, err := serializer.ApplyModifications(tt.modifications)

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

			if tt.validateFunc != nil {
				tt.validateFunc(t, testSerializerContent, result)
			}
		})
	}
}

func TestGradleSerializer_MinimalDiff(t *testing.T) {
	// 测试最小diff的核心功能。
	t.Run("Minimal diff validation", func(t *testing.T) {
		original := testSerializerContent
		serializer := NewGradleSerializer(original)

		// 创建一个简单的替换修改。
		versionStart := strings.Index(original, "version = '0.1.0-SNAPSHOT'")
		versionEnd := versionStart + len("version = '0.1.0-SNAPSHOT'")
		modifications := []Modification{
			{
				Type: ModificationTypeReplace,
				SourceRange: model.SourceRange{
					Start: model.SourcePosition{StartPos: versionStart},
					End:   model.SourcePosition{StartPos: versionEnd},
				},
				OldText: "version = '0.1.0-SNAPSHOT'",
				NewText: "version = '1.0.0'",
			},
		}

		result, err := serializer.ApplyModifications(modifications)
		if err != nil {
			t.Fatalf("Failed to apply modifications: %v", err)
		}

		// 计算实际的diff。
		originalLines := strings.Split(original, "\n")
		resultLines := strings.Split(result, "\n")

		if len(originalLines) != len(resultLines) {
			t.Errorf("Line count should remain the same. Original: %d, Result: %d", len(originalLines), len(resultLines))
		}

		// 计算变更的行数。
		changedLines := 0
		for i := 0; i < len(originalLines) && i < len(resultLines); i++ {
			if originalLines[i] != resultLines[i] {
				changedLines++
			}
		}

		// 应该只有一行发生变更。
		if changedLines != 1 {
			t.Errorf("Expected exactly 1 changed line, got %d", changedLines)
		}

		// 验证变更的行包含新版本。
		foundNewVersion := false
		for _, line := range resultLines {
			if strings.Contains(line, "version = '1.0.0'") {
				foundNewVersion = true
				break
			}
		}

		if !foundNewVersion {
			t.Errorf("New version not found in result")
		}
	})
}

func TestGradleSerializer_ValidateModifications(t *testing.T) {
	serializer := NewGradleSerializer(testSerializerContent)

	tests := []struct {
		name          string
		modifications []Modification
		expectErrors  int
	}{
		{
			name: "Valid modifications",
			modifications: []Modification{
				{
					Type: ModificationTypeReplace,
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{StartPos: strings.Index(testSerializerContent, "version = '0.1.0-SNAPSHOT'")},
						End:   model.SourcePosition{StartPos: strings.Index(testSerializerContent, "version = '0.1.0-SNAPSHOT'") + len("version = '0.1.0-SNAPSHOT'")},
					},
					OldText: "version = '0.1.0-SNAPSHOT'",
					NewText: "version = '1.0.0'",
				},
			},
			expectErrors: 0,
		},
		{
			name: "Invalid start position",
			modifications: []Modification{
				{
					Type: ModificationTypeReplace,
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{StartPos: -1, EndPos: -1},
						End:   model.SourcePosition{StartPos: 10, EndPos: 10},
					},
					OldText: "test",
					NewText: "new",
				},
			},
			expectErrors: 1,
		},
		{
			name: "End position exceeds text length",
			modifications: []Modification{
				{
					Type: ModificationTypeReplace,
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{StartPos: 0, EndPos: 0},
						End:   model.SourcePosition{StartPos: 99999, EndPos: 99999},
					},
					OldText: "test",
					NewText: "new",
				},
			},
			expectErrors: 1,
		},
		{
			name: "Start position greater than end position",
			modifications: []Modification{
				{
					Type: ModificationTypeReplace,
					SourceRange: model.SourceRange{
						Start: model.SourcePosition{StartPos: 100, EndPos: 100},
						End:   model.SourcePosition{StartPos: 50, EndPos: 50},
					},
					OldText: "test",
					NewText: "new",
				},
			},
			expectErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := serializer.ValidateModifications(tt.modifications)

			if len(errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d: %v", tt.expectErrors, len(errors), errors)
			}
		})
	}
}

func TestGradleSerializer_GenerateDiff(t *testing.T) {
	serializer := NewGradleSerializer(testSerializerContent)

	modifications := []Modification{
		{
			Type: ModificationTypeReplace,
			SourceRange: model.SourceRange{
				Start: model.SourcePosition{Line: 6, StartPos: 85},
				End:   model.SourcePosition{Line: 6, StartPos: 106},
			},
			OldText:     "version = '0.1.0-SNAPSHOT'",
			NewText:     "version = '1.0.0'",
			Description: "Update version",
		},
		{
			Type: ModificationTypeInsert,
			SourceRange: model.SourceRange{
				Start: model.SourcePosition{Line: 10, StartPos: 200},
				End:   model.SourcePosition{Line: 10, StartPos: 200},
			},
			OldText:     "",
			NewText:     "    implementation 'new:dependency:1.0'",
			Description: "Add new dependency",
		},
	}

	diffLines := serializer.GenerateDiff(modifications)

	if len(diffLines) != 3 { // 1 remove + 1 add + 1 add。
		t.Errorf("Expected 3 diff lines, got %d", len(diffLines))
	}

	// 验证diff内容。
	foundRemove := false
	foundAdd := false

	for _, diffLine := range diffLines {
		if diffLine.Type == DiffTypeRemove && strings.Contains(diffLine.Content, "0.1.0-SNAPSHOT") {
			foundRemove = true
		}
		if diffLine.Type == DiffTypeAdd && strings.Contains(diffLine.Content, "1.0.0") {
			foundAdd = true
		}
	}

	if !foundRemove {
		t.Errorf("Should find remove diff for old version")
	}

	if !foundAdd {
		t.Errorf("Should find add diff for new version")
	}
}

func TestGradleSerializer_GetModificationSummary(t *testing.T) {
	serializer := NewGradleSerializer(testSerializerContent)

	modifications := []Modification{
		{Type: ModificationTypeReplace, Description: "Update version"},
		{Type: ModificationTypeReplace, Description: "Update group"},
		{Type: ModificationTypeInsert, Description: "Add dependency"},
		{Type: ModificationTypeDelete, Description: "Remove comment"},
	}

	summary := serializer.GetModificationSummary(modifications)

	if summary.TotalModifications != 4 {
		t.Errorf("Expected 4 total modifications, got %d", summary.TotalModifications)
	}

	if summary.ByType[ModificationTypeReplace] != 2 {
		t.Errorf("Expected 2 replace modifications, got %d", summary.ByType[ModificationTypeReplace])
	}

	if summary.ByType[ModificationTypeInsert] != 1 {
		t.Errorf("Expected 1 insert modification, got %d", summary.ByType[ModificationTypeInsert])
	}

	if summary.ByType[ModificationTypeDelete] != 1 {
		t.Errorf("Expected 1 delete modification, got %d", summary.ByType[ModificationTypeDelete])
	}

	if len(summary.Descriptions) != 4 {
		t.Errorf("Expected 4 descriptions, got %d", len(summary.Descriptions))
	}
}
