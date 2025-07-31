package parser

import (
	"strings"
	"testing"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

const testSourceAwareContent = `plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
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
    implementation 'mysql:mysql-connector-java:8.0.29'
    testImplementation 'org.junit.jupiter:junit-jupiter-api:5.8.2'
}
`

func TestNewSourceAwareParser(t *testing.T) {
	parser := NewSourceAwareParser()
	if parser == nil {
		t.Fatal("NewSourceAwareParser should not return nil")
	}

	if parser.GradleParser == nil {
		t.Fatal("SourceAwareParser should have a GradleParser")
	}
}

func TestSourceAwareParser_ParseWithSourceMapping(t *testing.T) {
	parser := NewSourceAwareParser()

	result, err := parser.ParseWithSourceMapping(testSourceAwareContent)
	if err != nil {
		t.Fatalf("ParseWithSourceMapping failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result should not be nil")
	}

	if result.SourceMappedProject == nil {
		t.Fatal("SourceMappedProject should not be nil")
	}

	// 验证原始文本被保存。
	if result.SourceMappedProject.OriginalText != testSourceAwareContent {
		t.Error("Original text should be preserved")
	}

	// 验证行分割。
	expectedLines := strings.Split(testSourceAwareContent, "\n")
	if len(result.SourceMappedProject.Lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(result.SourceMappedProject.Lines))
	}
}

func TestSourceAwareParser_ParseSourceMappedDependencies(t *testing.T) {
	parser := NewSourceAwareParser()

	result, err := parser.ParseWithSourceMapping(testSourceAwareContent)
	if err != nil {
		t.Fatalf("ParseWithSourceMapping failed: %v", err)
	}

	dependencies := result.SourceMappedProject.SourceMappedDependencies

	// 应该找到至少3个依赖。
	if len(dependencies) < 3 {
		t.Errorf("Expected at least 3 dependencies, got %d", len(dependencies))
	}

	// 验证依赖的源码位置信息。
	for i, dep := range dependencies {
		if dep.SourceRange.Start.Line <= 0 {
			t.Errorf("Dependency %d should have valid start line, got %d", i, dep.SourceRange.Start.Line)
		}

		if dep.SourceRange.Start.StartPos < 0 {
			t.Errorf("Dependency %d should have valid start position, got %d", i, dep.SourceRange.Start.StartPos)
		}

		if dep.RawText == "" {
			t.Errorf("Dependency %d should have raw text", i)
		}

		// 验证原始文本确实存在于源码中。
		if !strings.Contains(testSourceAwareContent, dep.RawText) {
			t.Errorf("Dependency %d raw text '%s' not found in source", i, dep.RawText)
		}
	}

	// 查找特定依赖并验证其信息。
	var mysqlDep *model.SourceMappedDependency
	for _, dep := range dependencies {
		if dep.Group == "mysql" && dep.Name == "mysql-connector-java" {
			mysqlDep = dep
			break
		}
	}

	if mysqlDep == nil {
		t.Error("Should find MySQL dependency")
	} else {
		if mysqlDep.Version != "8.0.29" {
			t.Errorf("MySQL dependency version should be '8.0.29', got '%s'", mysqlDep.Version)
		}

		if !strings.Contains(mysqlDep.RawText, "mysql:mysql-connector-java:8.0.29") {
			t.Errorf("MySQL dependency raw text should contain full coordinates, got '%s'", mysqlDep.RawText)
		}
	}
}

func TestSourceAwareParser_ParseSourceMappedPlugins(t *testing.T) {
	parser := NewSourceAwareParser()

	result, err := parser.ParseWithSourceMapping(testSourceAwareContent)
	if err != nil {
		t.Fatalf("ParseWithSourceMapping failed: %v", err)
	}

	plugins := result.SourceMappedProject.SourceMappedPlugins

	// 应该找到至少2个插件。
	if len(plugins) < 2 {
		t.Errorf("Expected at least 2 plugins, got %d", len(plugins))
	}

	// 验证插件的源码位置信息。
	for i, plugin := range plugins {
		if plugin.SourceRange.Start.Line <= 0 {
			t.Errorf("Plugin %d should have valid start line, got %d", i, plugin.SourceRange.Start.Line)
		}

		if plugin.RawText == "" {
			t.Errorf("Plugin %d should have raw text", i)
		}

		// 验证原始文本确实存在于源码中。
		if !strings.Contains(testSourceAwareContent, plugin.RawText) {
			t.Errorf("Plugin %d raw text '%s' not found in source", i, plugin.RawText)
		}
	}

	// 查找Spring Boot插件并验证其信息。
	var springBootPlugin *model.SourceMappedPlugin
	for _, plugin := range plugins {
		if plugin.ID == "org.springframework.boot" {
			springBootPlugin = plugin
			break
		}
	}

	if springBootPlugin == nil {
		t.Error("Should find Spring Boot plugin")
	} else {
		if springBootPlugin.Version != "2.7.0" {
			t.Errorf("Spring Boot plugin version should be '2.7.0', got '%s'", springBootPlugin.Version)
		}

		if !strings.Contains(springBootPlugin.RawText, "org.springframework.boot") {
			t.Errorf("Spring Boot plugin raw text should contain plugin ID, got '%s'", springBootPlugin.RawText)
		}
	}
}

func TestSourceAwareParser_ParseSourceMappedRepositories(t *testing.T) {
	parser := NewSourceAwareParser()

	result, err := parser.ParseWithSourceMapping(testSourceAwareContent)
	if err != nil {
		t.Fatalf("ParseWithSourceMapping failed: %v", err)
	}

	repositories := result.SourceMappedProject.SourceMappedRepositories

	// 应该找到至少2个仓库。
	if len(repositories) < 2 {
		t.Errorf("Expected at least 2 repositories, got %d", len(repositories))
	}

	// 验证仓库的源码位置信息。
	for i, repo := range repositories {
		if repo.SourceRange.Start.Line <= 0 {
			t.Errorf("Repository %d should have valid start line, got %d", i, repo.SourceRange.Start.Line)
		}

		if repo.RawText == "" {
			t.Errorf("Repository %d should have raw text", i)
		}

		// 验证原始文本确实存在于源码中。
		if !strings.Contains(testSourceAwareContent, repo.RawText) {
			t.Errorf("Repository %d raw text '%s' not found in source", i, repo.RawText)
		}
	}

	// 查找mavenCentral仓库。
	var mavenCentralRepo *model.SourceMappedRepository
	for _, repo := range repositories {
		if repo.Name == "mavenCentral" {
			mavenCentralRepo = repo
			break
		}
	}

	if mavenCentralRepo == nil {
		t.Error("Should find mavenCentral repository")
	} else {
		if mavenCentralRepo.Type != "maven" {
			t.Errorf("mavenCentral repository type should be 'maven', got '%s'", mavenCentralRepo.Type)
		}
	}
}

func TestSourceAwareParser_ParseSourceMappedProperties(t *testing.T) {
	parser := NewSourceAwareParser()

	result, err := parser.ParseWithSourceMapping(testSourceAwareContent)
	if err != nil {
		t.Fatalf("ParseWithSourceMapping failed: %v", err)
	}

	properties := result.SourceMappedProject.SourceMappedProperties

	// 应该找到至少3个属性。
	if len(properties) < 3 {
		t.Errorf("Expected at least 3 properties, got %d", len(properties))
	}

	// 验证属性的源码位置信息。
	for i, prop := range properties {
		if prop.SourceRange.Start.Line <= 0 {
			t.Errorf("Property %d should have valid start line, got %d", i, prop.SourceRange.Start.Line)
		}

		if prop.Key == "" {
			t.Errorf("Property %d should have a key", i)
		}

		if prop.Value == "" {
			t.Errorf("Property %d should have a value", i)
		}

		if prop.RawText == "" {
			t.Errorf("Property %d should have raw text", i)
		}
	}

	// 查找特定属性并验证。
	expectedProperties := map[string]string{
		"group":       "com.example",
		"version":     "0.1.0-SNAPSHOT",
		"description": "Test project",
	}

	for expectedKey, expectedValue := range expectedProperties {
		found := false
		for _, prop := range properties {
			if prop.Key == expectedKey {
				found = true
				if prop.Value != expectedValue {
					t.Errorf("Property '%s' should have value '%s', got '%s'", expectedKey, expectedValue, prop.Value)
				}
				break
			}
		}

		if !found {
			t.Errorf("Should find property '%s'", expectedKey)
		}
	}
}

func TestSourceAwareParser_PositionAccuracy(t *testing.T) {
	// 测试位置信息的准确性。
	parser := NewSourceAwareParser()

	result, err := parser.ParseWithSourceMapping(testSourceAwareContent)
	if err != nil {
		t.Fatalf("ParseWithSourceMapping failed: %v", err)
	}

	// 验证依赖位置的准确性。
	for _, dep := range result.SourceMappedProject.SourceMappedDependencies {
		// 从原始文本中提取指定位置的内容。
		startPos := dep.SourceRange.Start.StartPos
		endPos := dep.SourceRange.End.StartPos

		if startPos < 0 || endPos > len(testSourceAwareContent) || startPos >= endPos {
			t.Errorf("Invalid position range for dependency: start=%d, end=%d", startPos, endPos)
			continue
		}

		extractedText := testSourceAwareContent[startPos:endPos]

		// 提取的文本应该包含依赖的原始文本。
		if !strings.Contains(extractedText, dep.RawText) {
			t.Errorf("Extracted text '%s' should contain raw text '%s'", extractedText, dep.RawText)
		}
	}

	// 验证插件位置的准确性。
	for _, plugin := range result.SourceMappedProject.SourceMappedPlugins {
		startPos := plugin.SourceRange.Start.StartPos
		endPos := plugin.SourceRange.End.StartPos

		if startPos < 0 || endPos > len(testSourceAwareContent) || startPos >= endPos {
			t.Errorf("Invalid position range for plugin: start=%d, end=%d", startPos, endPos)
			continue
		}

		extractedText := testSourceAwareContent[startPos:endPos]

		// 提取的文本应该包含插件的原始文本。
		if !strings.Contains(extractedText, plugin.RawText) {
			t.Errorf("Extracted text '%s' should contain raw text '%s'", extractedText, plugin.RawText)
		}
	}
}
