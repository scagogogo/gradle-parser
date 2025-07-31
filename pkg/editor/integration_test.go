package editor

import (
	"strings"
	"testing"

	"github.com/scagogogo/gradle-parser/pkg/parser"
)

// 复杂的测试用例，包含各种Gradle语法。
const complexGradleContent = `// Complex Gradle build file。
buildscript {
    repositories {
        gradlePluginPortal()
        mavenCentral()
    }
    dependencies {
        classpath 'org.springframework.boot:spring-boot-gradle-plugin:2.7.0'
    }
}

plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
    id 'org.jetbrains.kotlin.jvm' version '1.6.21'
}

group = 'com.example'
version = '0.1.0-SNAPSHOT'
description = 'A complex Spring Boot application'

java {
    sourceCompatibility = JavaVersion.VERSION_11
    targetCompatibility = JavaVersion.VERSION_11
}

configurations {
    compileOnly {
        extendsFrom annotationProcessor
    }
}

repositories {
    mavenCentral()
    google()
    maven {
        url 'https://repo.spring.io/milestone'
        name 'Spring Milestones'
    }
    maven { url 'https://jitpack.io' }
}

dependencies {
    // Spring Boot starters。
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
    implementation 'org.springframework.boot:spring-boot-starter-security'
    implementation 'org.springframework.boot:spring-boot-starter-actuator'
    
    // Database。
    implementation 'mysql:mysql-connector-java:8.0.29'
    implementation 'org.flywaydb:flyway-core:8.5.13'
    
    // Utilities。
    implementation 'com.google.guava:guava:31.0-jre'
    implementation 'org.apache.commons:commons-lang3:3.12.0'
    implementation 'com.fasterxml.jackson.module:jackson-module-kotlin'
    
    // Kotlin。
    implementation 'org.jetbrains.kotlin:kotlin-reflect'
    implementation 'org.jetbrains.kotlin:kotlin-stdlib-jdk8'
    
    // Development tools。
    developmentOnly 'org.springframework.boot:spring-boot-devtools'
    annotationProcessor 'org.springframework.boot:spring-boot-configuration-processor'
    
    // Testing。
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
    testImplementation 'org.springframework.security:spring-security-test'
    testImplementation 'org.junit.jupiter:junit-jupiter-api:5.8.2'
    testImplementation 'org.junit.jupiter:junit-jupiter-params:5.8.2'
    testRuntimeOnly 'org.junit.jupiter:junit-jupiter-engine:5.8.2'
    testImplementation 'org.testcontainers:junit-jupiter:1.17.3'
    testImplementation 'org.testcontainers:mysql:1.17.3'
}

tasks.named('test') {
    useJUnitPlatform()
    systemProperty 'spring.profiles.active', 'test'
}

task integrationTest(type: Test) {
    group = 'verification'
    description = 'Runs integration tests'
    useJUnitPlatform {
        includeTags 'integration'
    }
}

bootJar {
    archiveFileName = "${project.name}-${project.version}.jar"
    manifest {
        attributes(
            'Implementation-Title': project.name,
            'Implementation-Version': project.version
        )
    }
}
`

func TestComplexGradleEditing(t *testing.T) {
	// 创建位置感知解析器。
	sourceAwareParser := parser.NewSourceAwareParser()
	result, err := sourceAwareParser.ParseWithSourceMapping(complexGradleContent)
	if err != nil {
		t.Fatalf("Failed to parse complex content: %v", err)
	}

	editor := NewGradleEditor(result.SourceMappedProject)

	// 执行多个编辑操作。
	t.Run("Multiple complex edits", func(t *testing.T) {
		// 1. 更新项目版本。
		err := editor.UpdateProperty("version", "1.0.0")
		if err != nil {
			t.Errorf("Failed to update version: %v", err)
		}

		// 2. 更新Spring Boot插件版本。
		err = editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
		if err != nil {
			t.Errorf("Failed to update Spring Boot plugin: %v", err)
		}

		// 3. 更新Kotlin插件版本。
		err = editor.UpdatePluginVersion("org.jetbrains.kotlin.jvm", "1.7.10")
		if err != nil {
			t.Errorf("Failed to update Kotlin plugin: %v", err)
		}

		// 4. 更新MySQL依赖版本。
		err = editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.30")
		if err != nil {
			t.Errorf("Failed to update MySQL version: %v", err)
		}

		// 5. 更新Guava版本。
		err = editor.UpdateDependencyVersion("com.google.guava", "guava", "31.1-jre")
		if err != nil {
			t.Errorf("Failed to update Guava version: %v", err)
		}

		// 6. 为Spring Boot starter添加版本号。
		err = editor.UpdateDependencyVersion("org.springframework.boot", "spring-boot-starter-web", "2.7.2")
		if err != nil {
			t.Errorf("Failed to add version to Spring Boot starter: %v", err)
		}

		// 7. 添加新依赖。
		err = editor.AddDependency("org.apache.commons", "commons-text", "1.9", "implementation")
		if err != nil {
			t.Errorf("Failed to add new dependency: %v", err)
		}

		// 验证修改数量。
		modifications := editor.GetModifications()
		expectedModifications := 7
		if len(modifications) != expectedModifications {
			t.Errorf("Expected %d modifications, got %d", expectedModifications, len(modifications))
		}

		// 应用修改。
		serializer := NewGradleSerializer(complexGradleContent)
		finalText, err := serializer.ApplyModifications(modifications)
		if err != nil {
			t.Fatalf("Failed to apply modifications: %v", err)
		}

		// 验证结果。
		validateComplexEdits(t, complexGradleContent, finalText)
	})
}

func validateComplexEdits(t *testing.T, original, result string) {
	// 验证所有预期的更改都已应用。
	expectedChanges := map[string]string{
		"version = '0.1.0-SNAPSHOT'":                         "version = '1.0.0'",
		"id 'org.springframework.boot' version '2.7.0'":      "id 'org.springframework.boot' version '2.7.2'",
		"id 'org.jetbrains.kotlin.jvm' version '1.6.21'":     "id 'org.jetbrains.kotlin.jvm' version '1.7.10'",
		"'mysql:mysql-connector-java:8.0.29'":                "'mysql:mysql-connector-java:8.0.30'",
		"'com.google.guava:guava:31.0-jre'":                  "'com.google.guava:guava:31.1-jre'",
		"'org.springframework.boot:spring-boot-starter-web'": "'org.springframework.boot:spring-boot-starter-web:2.7.2'",
	}

	for oldText, newText := range expectedChanges {
		if strings.Contains(result, oldText) {
			t.Errorf("Old text '%s' should have been replaced", oldText)
		}
		if !strings.Contains(result, newText) {
			t.Errorf("New text '%s' should be present in result", newText)
		}
	}

	// 验证新添加的依赖。
	if !strings.Contains(result, "commons-text") {
		t.Errorf("New dependency 'commons-text' should be added")
	}

	// 验证未修改的内容保持不变。
	unchangedContent := []string{
		"// Complex Gradle build file",
		"buildscript {",
		"repositories {",
		"group = 'com.example'",
		"description = 'A complex Spring Boot application'",
		"tasks.named('test') {",
		"bootJar {",
	}

	for _, content := range unchangedContent {
		if !strings.Contains(result, content) {
			t.Errorf("Unchanged content '%s' should be preserved", content)
		}
	}
}

func TestMinimalDiffValidation(t *testing.T) {
	// 测试最小diff的详细验证。
	sourceAwareParser := parser.NewSourceAwareParser()
	result, err := sourceAwareParser.ParseWithSourceMapping(complexGradleContent)
	if err != nil {
		t.Fatalf("Failed to parse content: %v", err)
	}

	editor := NewGradleEditor(result.SourceMappedProject)

	// 执行单个修改。
	err = editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.30")
	if err != nil {
		t.Fatalf("Failed to update dependency: %v", err)
	}

	// 应用修改。
	serializer := NewGradleSerializer(complexGradleContent)
	modifications := editor.GetModifications()
	finalText, err := serializer.ApplyModifications(modifications)
	if err != nil {
		t.Fatalf("Failed to apply modifications: %v", err)
	}

	// 详细的diff分析。
	originalLines := strings.Split(complexGradleContent, "\n")
	resultLines := strings.Split(finalText, "\n")

	if len(originalLines) != len(resultLines) {
		t.Errorf("Line count changed: original %d, result %d", len(originalLines), len(resultLines))
	}

	// 计算实际变更的行数。
	changedLines := 0
	changedLineNumbers := []int{}

	for i := 0; i < len(originalLines) && i < len(resultLines); i++ {
		if originalLines[i] != resultLines[i] {
			changedLines++
			changedLineNumbers = append(changedLineNumbers, i+1)
		}
	}

	// 应该只有一行发生变更。
	if changedLines != 1 {
		t.Errorf("Expected exactly 1 changed line, got %d (lines: %v)", changedLines, changedLineNumbers)

		// 输出变更的行以便调试。
		for _, lineNum := range changedLineNumbers {
			if lineNum <= len(originalLines) && lineNum <= len(resultLines) {
				t.Logf("Line %d changed:", lineNum)
				t.Logf("  Original: %s", originalLines[lineNum-1])
				t.Logf("  Result:   %s", resultLines[lineNum-1])
			}
		}
	}

	// 验证变更的行包含预期的内容。
	if changedLines == 1 {
		changedLineIndex := changedLineNumbers[0] - 1
		changedLine := resultLines[changedLineIndex]

		if !strings.Contains(changedLine, "8.0.30") {
			t.Errorf("Changed line should contain new version '8.0.30', got: %s", changedLine)
		}

		if strings.Contains(changedLine, "8.0.29") {
			t.Errorf("Changed line should not contain old version '8.0.29', got: %s", changedLine)
		}
	}
}

func TestBatchEditingMinimalDiff(t *testing.T) {
	// 测试批量编辑时的最小diff。
	sourceAwareParser := parser.NewSourceAwareParser()
	result, err := sourceAwareParser.ParseWithSourceMapping(complexGradleContent)
	if err != nil {
		t.Fatalf("Failed to parse content: %v", err)
	}

	editor := NewGradleEditor(result.SourceMappedProject)

	// 执行3个修改操作。
	modifications := []struct {
		action func() error
		desc   string
	}{
		{
			action: func() error { return editor.UpdateProperty("version", "1.0.0") },
			desc:   "version update",
		},
		{
			action: func() error { return editor.UpdatePluginVersion("org.springframework.boot", "2.7.2") },
			desc:   "Spring Boot plugin update",
		},
		{
			action: func() error { return editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.30") },
			desc:   "MySQL dependency update",
		},
	}

	for _, mod := range modifications {
		if err := mod.action(); err != nil {
			t.Fatalf("Failed to execute %s: %v", mod.desc, err)
		}
	}

	// 应用修改。
	serializer := NewGradleSerializer(complexGradleContent)
	allModifications := editor.GetModifications()
	finalText, err := serializer.ApplyModifications(allModifications)
	if err != nil {
		t.Fatalf("Failed to apply modifications: %v", err)
	}

	// 验证只有3行发生变更。
	originalLines := strings.Split(complexGradleContent, "\n")
	resultLines := strings.Split(finalText, "\n")

	changedLines := 0
	for i := 0; i < len(originalLines) && i < len(resultLines); i++ {
		if originalLines[i] != resultLines[i] {
			changedLines++
		}
	}

	expectedChangedLines := 3
	if changedLines != expectedChangedLines {
		t.Errorf("Expected exactly %d changed lines, got %d", expectedChangedLines, changedLines)
	}

	// 验证修改的准确性。
	if !strings.Contains(finalText, "version = '1.0.0'") {
		t.Errorf("Version should be updated to 1.0.0")
	}

	if !strings.Contains(finalText, "id 'org.springframework.boot' version '2.7.2'") {
		t.Errorf("Spring Boot plugin should be updated to 2.7.2")
	}

	if !strings.Contains(finalText, "'mysql:mysql-connector-java:8.0.30'") {
		t.Errorf("MySQL dependency should be updated to 8.0.30")
	}
}
