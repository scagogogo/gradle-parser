// Package editor 提供全面的编辑器功能测试
package editor

import (
	"strings"
	"testing"

	"github.com/scagogogo/gradle-parser/pkg/parser"
)

// 综合测试用例：验证完整的编辑器工作流程
func TestComprehensiveEditorWorkflow(t *testing.T) {
	// 测试用的Gradle文件内容
	gradleContent := `plugins {
    id 'java'
    id 'org.springframework.boot' version '2.6.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
}

group = 'com.example'
version = '0.1.0-SNAPSHOT'
description = 'Test project'

repositories {
    mavenCentral()
    google()
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'mysql:mysql-connector-java:8.0.28'
    implementation 'com.google.guava:guava:30.1-jre'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}`

	t.Run("Complete workflow test", func(t *testing.T) {
		// 1. 解析阶段
		sourceAwareParser := parser.NewSourceAwareParser()
		result, err := sourceAwareParser.ParseWithSourceMapping(gradleContent)
		if err != nil {
			t.Fatalf("Failed to parse content: %v", err)
		}

		// 验证解析结果
		if result.SourceMappedProject == nil {
			t.Fatal("SourceMappedProject should not be nil")
		}

		if len(result.SourceMappedProject.SourceMappedDependencies) < 3 {
			t.Errorf("Expected at least 3 dependencies, got %d", len(result.SourceMappedProject.SourceMappedDependencies))
		}

		if len(result.SourceMappedProject.SourceMappedPlugins) < 2 {
			t.Errorf("Expected at least 2 plugins, got %d", len(result.SourceMappedProject.SourceMappedPlugins))
		}

		// 2. 编辑阶段
		editor := NewGradleEditor(result.SourceMappedProject)

		// 执行多种类型的修改
		modifications := []struct {
			name   string
			action func() error
		}{
			{
				name: "Update project version",
				action: func() error {
					return editor.UpdateProperty("version", "1.0.0")
				},
			},
			{
				name: "Update Spring Boot plugin version",
				action: func() error {
					return editor.UpdatePluginVersion("org.springframework.boot", "2.7.0")
				},
			},
			{
				name: "Update MySQL dependency version",
				action: func() error {
					return editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.30")
				},
			},
			{
				name: "Update Guava dependency version",
				action: func() error {
					return editor.UpdateDependencyVersion("com.google.guava", "guava", "31.1-jre")
				},
			},
			{
				name: "Add new dependency",
				action: func() error {
					return editor.AddDependency("org.apache.commons", "commons-text", "1.9", "implementation")
				},
			},
		}

		// 执行所有修改
		for _, mod := range modifications {
			if err := mod.action(); err != nil {
				t.Errorf("Failed to execute %s: %v", mod.name, err)
			}
		}

		// 验证修改操作被正确记录
		allModifications := editor.GetModifications()
		expectedModifications := 5
		if len(allModifications) != expectedModifications {
			t.Errorf("Expected %d modifications, got %d", expectedModifications, len(allModifications))
		}

		// 3. 序列化阶段
		serializer := NewGradleSerializer(gradleContent)

		// 验证修改操作的有效性
		validationErrors := serializer.ValidateModifications(allModifications)
		if len(validationErrors) > 0 {
			t.Errorf("Validation errors found: %v", validationErrors)
		}

		// 应用修改
		finalText, err := serializer.ApplyModifications(allModifications)
		if err != nil {
			t.Fatalf("Failed to apply modifications: %v", err)
		}

		// 4. 验证最小diff
		originalLines := strings.Split(gradleContent, "\n")
		resultLines := strings.Split(finalText, "\n")

		changedLines := 0
		for i := 0; i < len(originalLines) && i < len(resultLines); i++ {
			if originalLines[i] != resultLines[i] {
				changedLines++
			}
		}

		// 应该有4行被修改 + 1行新增 = 5行变化
		expectedChanges := 5
		if changedLines != expectedChanges {
			t.Errorf("Expected %d changed lines, got %d", expectedChanges, changedLines)
		}

		// 5. 验证具体修改内容
		verifications := []struct {
			description string
			check       func(string) bool
		}{
			{
				description: "Version should be updated to 1.0.0",
				check: func(text string) bool {
					return strings.Contains(text, "version = '1.0.0'") && !strings.Contains(text, "0.1.0-SNAPSHOT")
				},
			},
			{
				description: "Spring Boot plugin should be updated to 2.7.0",
				check: func(text string) bool {
					return strings.Contains(text, "id 'org.springframework.boot' version '2.7.0'")
				},
			},
			{
				description: "MySQL version should be updated to 8.0.30",
				check: func(text string) bool {
					return strings.Contains(text, "mysql:mysql-connector-java:8.0.30")
				},
			},
			{
				description: "Guava version should be updated to 31.1-jre",
				check: func(text string) bool {
					return strings.Contains(text, "com.google.guava:guava:31.1-jre")
				},
			},
			{
				description: "New commons-text dependency should be added",
				check: func(text string) bool {
					return strings.Contains(text, "org.apache.commons:commons-text:1.9")
				},
			},
		}

		for _, verification := range verifications {
			if !verification.check(finalText) {
				t.Errorf("Verification failed: %s", verification.description)
			}
		}

		// 6. 生成和验证diff
		diffLines := serializer.GenerateDiff(allModifications)
		if len(diffLines) == 0 {
			t.Error("Should generate diff lines")
		}

		// 验证diff包含预期的修改类型
		hasReplace := false
		hasInsert := false
		for _, diffLine := range diffLines {
			if diffLine.Type == DiffTypeRemove || diffLine.Type == DiffTypeAdd {
				hasReplace = true
			}
			if diffLine.Type == DiffTypeAdd && strings.Contains(diffLine.Content, "commons-text") {
				hasInsert = true
			}
		}

		if !hasReplace {
			t.Error("Diff should contain replace operations")
		}

		if !hasInsert {
			t.Error("Diff should contain insert operations")
		}

		// 7. 获取修改摘要
		summary := serializer.GetModificationSummary(allModifications)
		if summary.TotalModifications != expectedModifications {
			t.Errorf("Summary should show %d modifications, got %d", expectedModifications, summary.TotalModifications)
		}

		if summary.ByType[ModificationTypeReplace] != 4 {
			t.Errorf("Expected 4 replace modifications, got %d", summary.ByType[ModificationTypeReplace])
		}

		if summary.ByType[ModificationTypeInsert] != 1 {
			t.Errorf("Expected 1 insert modification, got %d", summary.ByType[ModificationTypeInsert])
		}
	})
}

// 测试边界情况和错误处理
func TestEditorErrorHandling(t *testing.T) {
	gradleContent := `plugins {
    id 'java'
}

group = 'com.example'
version = '1.0.0'

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
}`

	t.Run("Error handling tests", func(t *testing.T) {
		sourceAwareParser := parser.NewSourceAwareParser()
		result, err := sourceAwareParser.ParseWithSourceMapping(gradleContent)
		if err != nil {
			t.Fatalf("Failed to parse content: %v", err)
		}

		editor := NewGradleEditor(result.SourceMappedProject)

		// 测试更新不存在的依赖
		err = editor.UpdateDependencyVersion("non.existent", "dependency", "1.0.0")
		if err == nil {
			t.Error("Should return error for non-existent dependency")
		}

		// 测试更新不存在的插件
		err = editor.UpdatePluginVersion("non.existent.plugin", "1.0.0")
		if err == nil {
			t.Error("Should return error for non-existent plugin")
		}

		// 测试更新不存在的属性
		err = editor.UpdateProperty("nonExistentProperty", "value")
		if err == nil {
			t.Error("Should return error for non-existent property")
		}

		// 测试相同版本更新（应该不产生修改）
		initialModCount := len(editor.GetModifications())
		err = editor.UpdateProperty("version", "1.0.0") // 相同版本
		if err != nil {
			t.Errorf("Should not return error for same version update: %v", err)
		}

		finalModCount := len(editor.GetModifications())
		if finalModCount != initialModCount {
			t.Error("Should not create modification for same version update")
		}
	})
}
