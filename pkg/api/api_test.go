// Package api 提供API测试
package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

// 测试用的Gradle文件内容
const testGradleContent = `plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
}

group = 'com.example'
version = '1.0.0'
description = 'Test project'

repositories {
    mavenCentral()
    google()
    maven { url 'https://jitpack.io' }
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'mysql:mysql-connector-java:8.0.29'
    implementation 'com.google.guava:guava:31.1-jre'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
    testImplementation 'org.junit.jupiter:junit-jupiter-api:5.8.2'
}`

// 创建临时测试文件
func createTempGradleFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "build.gradle")

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	return filePath
}

func TestParseFile(t *testing.T) {
	filePath := createTempGradleFile(t, testGradleContent)

	result, err := ParseFile(filePath)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if result == nil {
		t.Fatal("ParseFile() returned nil result")
	}

	if result.Project == nil {
		t.Fatal("ParseFile() returned nil project")
	}

	// 验证项目基本信息
	if result.Project.Group != "com.example" {
		t.Errorf("Expected group 'com.example', got '%s'", result.Project.Group)
	}

	if result.Project.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", result.Project.Version)
	}

	if result.Project.Description != "Test project" {
		t.Errorf("Expected description 'Test project', got '%s'", result.Project.Description)
	}
}

func TestParseString(t *testing.T) {
	result, err := ParseString(testGradleContent)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	if result == nil {
		t.Fatal("ParseString() returned nil result")
	}

	if result.Project == nil {
		t.Fatal("ParseString() returned nil project")
	}

	// 验证解析结果
	if result.Project.Group != "com.example" {
		t.Errorf("Expected group 'com.example', got '%s'", result.Project.Group)
	}
}

func TestParseReader(t *testing.T) {
	reader := strings.NewReader(testGradleContent)

	result, err := ParseReader(reader)
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

func TestGetDependencies(t *testing.T) {
	filePath := createTempGradleFile(t, testGradleContent)

	dependencies, err := GetDependencies(filePath)
	if err != nil {
		t.Fatalf("GetDependencies() error = %v", err)
	}

	if len(dependencies) == 0 {
		t.Error("GetDependencies() returned no dependencies")
	}

	// 验证是否包含预期的依赖
	foundSpringBoot := false
	foundMySQL := false

	for _, dep := range dependencies {
		if dep.Group == "org.springframework.boot" && dep.Name == "spring-boot-starter-web" {
			foundSpringBoot = true
		}
		if dep.Group == "mysql" && dep.Name == "mysql-connector-java" {
			foundMySQL = true
		}
	}

	if !foundSpringBoot {
		t.Error("GetDependencies() did not find Spring Boot dependency")
	}

	if !foundMySQL {
		t.Error("GetDependencies() did not find MySQL dependency")
	}
}

func TestGetPlugins(t *testing.T) {
	filePath := createTempGradleFile(t, testGradleContent)

	plugins, err := GetPlugins(filePath)
	if err != nil {
		t.Fatalf("GetPlugins() error = %v", err)
	}

	if len(plugins) == 0 {
		t.Error("GetPlugins() returned no plugins")
	}

	// 验证是否包含预期的插件
	foundJava := false
	foundSpringBoot := false

	for _, plugin := range plugins {
		if plugin.ID == "java" {
			foundJava = true
		}
		if plugin.ID == "org.springframework.boot" {
			foundSpringBoot = true
			if plugin.Version != "2.7.0" {
				t.Errorf("Expected Spring Boot version '2.7.0', got '%s'", plugin.Version)
			}
		}
	}

	if !foundJava {
		t.Error("GetPlugins() did not find Java plugin")
	}

	if !foundSpringBoot {
		t.Error("GetPlugins() did not find Spring Boot plugin")
	}
}

func TestGetRepositories(t *testing.T) {
	filePath := createTempGradleFile(t, testGradleContent)

	repositories, err := GetRepositories(filePath)
	if err != nil {
		t.Fatalf("GetRepositories() error = %v", err)
	}

	if len(repositories) == 0 {
		t.Error("GetRepositories() returned no repositories")
	}

	// 验证是否包含预期的仓库
	foundMavenCentral := false
	foundGoogle := false
	foundJitPack := false

	for _, repo := range repositories {
		if repo.Name == "mavenCentral" {
			foundMavenCentral = true
		}
		if repo.Name == "google" {
			foundGoogle = true
		}
		if strings.Contains(repo.URL, "jitpack.io") {
			foundJitPack = true
		}
	}

	if !foundMavenCentral {
		t.Error("GetRepositories() did not find Maven Central")
	}

	if !foundGoogle {
		t.Error("GetRepositories() did not find Google repository")
	}

	if !foundJitPack {
		t.Error("GetRepositories() did not find JitPack repository")
	}
}

func TestDependenciesByScope(t *testing.T) {
	// 创建测试依赖
	dependencies := []*model.Dependency{
		{Group: "org.springframework.boot", Name: "spring-boot-starter-web", Version: "2.7.0", Scope: "implementation"},
		{Group: "mysql", Name: "mysql-connector-java", Version: "8.0.29", Scope: "implementation"},
		{Group: "org.junit.jupiter", Name: "junit-jupiter-api", Version: "5.8.2", Scope: "testImplementation"},
	}

	dependencySets := DependenciesByScope(dependencies)

	if len(dependencySets) == 0 {
		t.Error("DependenciesByScope() returned no dependency sets")
	}

	// 验证分组结果
	foundImplementation := false
	foundTestImplementation := false

	for _, depSet := range dependencySets {
		if depSet.Scope == "implementation" {
			foundImplementation = true
			if len(depSet.Dependencies) < 2 {
				t.Errorf("Expected at least 2 implementation dependencies, got %d", len(depSet.Dependencies))
			}
		}
		if depSet.Scope == "testImplementation" {
			foundTestImplementation = true
			if len(depSet.Dependencies) < 1 {
				t.Errorf("Expected at least 1 testImplementation dependency, got %d", len(depSet.Dependencies))
			}
		}
	}

	if !foundImplementation {
		t.Error("DependenciesByScope() did not create implementation scope")
	}

	if !foundTestImplementation {
		t.Error("DependenciesByScope() did not create testImplementation scope")
	}
}

func TestProjectTypeDetection(t *testing.T) {
	// 测试Android项目检测
	androidPlugins := []*model.Plugin{
		{ID: "com.android.application", Version: "7.0.0"},
		{ID: "kotlin-android", Version: "1.5.30"},
	}

	if !IsAndroidProject(androidPlugins) {
		t.Error("IsAndroidProject() should return true for Android plugins")
	}

	// 测试Kotlin项目检测
	kotlinPlugins := []*model.Plugin{
		{ID: "org.jetbrains.kotlin.jvm", Version: "1.7.10"},
	}

	if !IsKotlinProject(kotlinPlugins) {
		t.Error("IsKotlinProject() should return true for Kotlin plugins")
	}

	// 测试Spring Boot项目检测
	springBootPlugins := []*model.Plugin{
		{ID: "org.springframework.boot", Version: "2.7.0"},
	}

	if !IsSpringBootProject(springBootPlugins) {
		t.Error("IsSpringBootProject() should return true for Spring Boot plugins")
	}

	// 测试非特定项目类型
	javaPlugins := []*model.Plugin{
		{ID: "java"},
	}

	if IsAndroidProject(javaPlugins) {
		t.Error("IsAndroidProject() should return false for Java-only plugins")
	}

	if IsKotlinProject(javaPlugins) {
		t.Error("IsKotlinProject() should return false for Java-only plugins")
	}

	if IsSpringBootProject(javaPlugins) {
		t.Error("IsSpringBootProject() should return false for Java-only plugins")
	}
}

func TestDefaultOptions(t *testing.T) {
	options := DefaultOptions()

	if options == nil {
		t.Fatal("DefaultOptions() returned nil")
	}

	// 验证默认选项
	if !options.SkipComments {
		t.Error("DefaultOptions() should have SkipComments = true")
	}

	if !options.CollectRawContent {
		t.Error("DefaultOptions() should have CollectRawContent = true")
	}

	if !options.ParsePlugins {
		t.Error("DefaultOptions() should have ParsePlugins = true")
	}

	if !options.ParseDependencies {
		t.Error("DefaultOptions() should have ParseDependencies = true")
	}

	if !options.ParseRepositories {
		t.Error("DefaultOptions() should have ParseRepositories = true")
	}

	if !options.ParseTasks {
		t.Error("DefaultOptions() should have ParseTasks = true")
	}
}

func TestNewParser(t *testing.T) {
	options := &Options{
		SkipComments:      false,
		CollectRawContent: false,
		ParsePlugins:      false,
		ParseDependencies: false,
		ParseRepositories: false,
		ParseTasks:        false,
	}

	parser := NewParser(options)
	if parser == nil {
		t.Fatal("NewParser() returned nil")
	}

	// 测试使用自定义选项的解析器
	result, err := parser.Parse(testGradleContent)
	if err != nil {
		t.Fatalf("Custom parser Parse() error = %v", err)
	}

	if result == nil {
		t.Fatal("Custom parser returned nil result")
	}
}

func TestParseFileWithSourceMapping(t *testing.T) {
	filePath := createTempGradleFile(t, testGradleContent)

	result, err := ParseFileWithSourceMapping(filePath)
	if err != nil {
		t.Fatalf("ParseFileWithSourceMapping() error = %v", err)
	}

	if result == nil {
		t.Fatal("ParseFileWithSourceMapping() returned nil result")
	}

	if result.SourceMappedProject == nil {
		t.Fatal("ParseFileWithSourceMapping() returned nil SourceMappedProject")
	}

	// 验证文件路径被设置
	if result.SourceMappedProject.FilePath != filePath {
		t.Errorf("Expected FilePath '%s', got '%s'", filePath, result.SourceMappedProject.FilePath)
	}

	// 验证原始文本被保存
	if result.SourceMappedProject.OriginalText == "" {
		t.Error("OriginalText should not be empty")
	}

	// 验证源码映射的组件
	if len(result.SourceMappedProject.SourceMappedDependencies) == 0 {
		t.Error("Should have source mapped dependencies")
	}

	if len(result.SourceMappedProject.SourceMappedPlugins) == 0 {
		t.Error("Should have source mapped plugins")
	}

	if len(result.SourceMappedProject.SourceMappedRepositories) == 0 {
		t.Error("Should have source mapped repositories")
	}
}

func TestCreateGradleEditor(t *testing.T) {
	filePath := createTempGradleFile(t, testGradleContent)

	editor, err := CreateGradleEditor(filePath)
	if err != nil {
		t.Fatalf("CreateGradleEditor() error = %v", err)
	}

	if editor == nil {
		t.Fatal("CreateGradleEditor() returned nil editor")
	}

	// 验证编辑器可以获取源码映射项目
	sourceMappedProject := editor.GetSourceMappedProject()
	if sourceMappedProject == nil {
		t.Fatal("Editor should have a SourceMappedProject")
	}

	// 验证初始修改列表为空
	modifications := editor.GetModifications()
	if len(modifications) != 0 {
		t.Errorf("New editor should have 0 modifications, got %d", len(modifications))
	}
}

func TestUpdateDependencyVersion(t *testing.T) {
	filePath := createTempGradleFile(t, testGradleContent)

	newText, err := UpdateDependencyVersion(filePath, "mysql", "mysql-connector-java", "8.0.30")
	if err != nil {
		t.Fatalf("UpdateDependencyVersion() error = %v", err)
	}

	if newText == "" {
		t.Error("UpdateDependencyVersion() returned empty text")
	}

	// 验证版本已更新
	if !strings.Contains(newText, "mysql:mysql-connector-java:8.0.30") {
		t.Error("Updated text should contain new version 8.0.30")
	}

	// 验证旧版本已被替换
	if strings.Contains(newText, "mysql:mysql-connector-java:8.0.29") {
		t.Error("Updated text should not contain old version 8.0.29")
	}
}

func TestUpdatePluginVersion(t *testing.T) {
	filePath := createTempGradleFile(t, testGradleContent)

	newText, err := UpdatePluginVersion(filePath, "org.springframework.boot", "2.7.2")
	if err != nil {
		t.Fatalf("UpdatePluginVersion() error = %v", err)
	}

	if newText == "" {
		t.Error("UpdatePluginVersion() returned empty text")
	}

	// 验证版本已更新
	if !strings.Contains(newText, "id 'org.springframework.boot' version '2.7.2'") {
		t.Error("Updated text should contain new version 2.7.2")
	}

	// 验证旧版本已被替换
	if strings.Contains(newText, "version '2.7.0'") {
		t.Error("Updated text should not contain old version 2.7.0")
	}
}

func TestErrorHandling(t *testing.T) {
	// 测试文件不存在的情况
	t.Run("File not found", func(t *testing.T) {
		_, err := ParseFile("non-existent-file.gradle")
		if err == nil {
			t.Error("ParseFile() should return error for non-existent file")
		}

		_, err = GetDependencies("non-existent-file.gradle")
		if err == nil {
			t.Error("GetDependencies() should return error for non-existent file")
		}

		_, err = GetPlugins("non-existent-file.gradle")
		if err == nil {
			t.Error("GetPlugins() should return error for non-existent file")
		}

		_, err = GetRepositories("non-existent-file.gradle")
		if err == nil {
			t.Error("GetRepositories() should return error for non-existent file")
		}

		_, err = ParseFileWithSourceMapping("non-existent-file.gradle")
		if err == nil {
			t.Error("ParseFileWithSourceMapping() should return error for non-existent file")
		}

		_, err = CreateGradleEditor("non-existent-file.gradle")
		if err == nil {
			t.Error("CreateGradleEditor() should return error for non-existent file")
		}

		_, err = UpdateDependencyVersion("non-existent-file.gradle", "group", "name", "version")
		if err == nil {
			t.Error("UpdateDependencyVersion() should return error for non-existent file")
		}

		_, err = UpdatePluginVersion("non-existent-file.gradle", "plugin", "version")
		if err == nil {
			t.Error("UpdatePluginVersion() should return error for non-existent file")
		}
	})

	// 测试无效内容的情况
	t.Run("Invalid content", func(t *testing.T) {
		invalidContent := "this is not a valid gradle file content"

		result, err := ParseString(invalidContent)
		// 解析器应该能处理无效内容，但可能返回空结果
		if err != nil {
			t.Logf("ParseString() with invalid content returned error: %v", err)
		} else if result != nil && result.Project != nil {
			// 如果没有错误，项目信息应该是空的或默认值
			t.Logf("ParseString() with invalid content returned project: %+v", result.Project)
		}
	})

	// 测试更新不存在的依赖
	t.Run("Update non-existent dependency", func(t *testing.T) {
		filePath := createTempGradleFile(t, testGradleContent)

		_, err := UpdateDependencyVersion(filePath, "non.existent", "dependency", "1.0.0")
		if err == nil {
			t.Error("UpdateDependencyVersion() should return error for non-existent dependency")
		}
	})

	// 测试更新不存在的插件
	t.Run("Update non-existent plugin", func(t *testing.T) {
		filePath := createTempGradleFile(t, testGradleContent)

		_, err := UpdatePluginVersion(filePath, "non.existent.plugin", "1.0.0")
		if err == nil {
			t.Error("UpdatePluginVersion() should return error for non-existent plugin")
		}
	})
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version constant should not be empty")
	}

	// 验证版本格式（应该是语义化版本）
	if !strings.Contains(Version, ".") {
		t.Error("Version should contain dots (semantic versioning)")
	}
}
