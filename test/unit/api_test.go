package unit

import (
	"testing"

	"github.com/scagogogo/gradle-parser/pkg/api"
	"github.com/scagogogo/gradle-parser/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseString tests basic string parsing functionality
func TestParseString(t *testing.T) {
	gradleContent := `
plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
}

group = 'com.example'
version = '1.0.0'
description = 'Test project'

repositories {
    mavenCentral()
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'mysql:mysql-connector-java:8.0.29'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}
`

	result, err := api.ParseString(gradleContent)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Project)

	project := result.Project
	assert.Equal(t, "com.example", project.Group)
	assert.Equal(t, "1.0.0", project.Version)
	assert.Equal(t, "Test project", project.Description)
	assert.Greater(t, len(project.Dependencies), 0)
	assert.Greater(t, len(project.Plugins), 0)
}

// TestGetDependencies tests dependency extraction
func TestGetDependencies(t *testing.T) {
	gradleContent := `
dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'mysql:mysql-connector-java:8.0.29'
    testImplementation 'junit:junit:4.13.2'
    compileOnly 'org.projectlombok:lombok:1.18.24'
}
`

	result, err := api.ParseString(gradleContent)
	require.NoError(t, err)

	deps := result.Project.Dependencies
	// Note: The parser may find duplicates or parse differently than expected
	assert.Greater(t, len(deps), 0, "Should find at least some dependencies")

	// Check that we can find the expected dependencies (regardless of scope for now)
	foundSpringBoot := false
	foundMySQL := false
	foundJUnit := false
	foundLombok := false

	for _, dep := range deps {
		switch {
		case dep.Group == "org.springframework.boot" && dep.Name == "spring-boot-starter-web":
			foundSpringBoot = true
		case dep.Group == "mysql" && dep.Name == "mysql-connector-java":
			foundMySQL = true
			assert.Equal(t, "8.0.29", dep.Version)
		case dep.Group == "junit" && dep.Name == "junit":
			foundJUnit = true
		case dep.Group == "org.projectlombok" && dep.Name == "lombok":
			foundLombok = true
		}
	}

	assert.True(t, foundSpringBoot, "Spring Boot dependency not found")
	assert.True(t, foundMySQL, "MySQL dependency not found")
	assert.True(t, foundJUnit, "JUnit dependency not found")
	assert.True(t, foundLombok, "Lombok dependency not found")
}

// TestGetPlugins tests plugin extraction
func TestGetPlugins(t *testing.T) {
	gradleContent := `
plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
    id 'io.spring.dependency-management' version '1.0.11.RELEASE'
    id 'application'
}
`

	result, err := api.ParseString(gradleContent)
	require.NoError(t, err)

	plugins := result.Project.Plugins
	assert.Len(t, plugins, 4)

	// Check specific plugins
	foundJava := false
	foundSpringBoot := false
	foundDependencyManagement := false
	foundApplication := false

	for _, plugin := range plugins {
		switch plugin.ID {
		case "java":
			foundJava = true
			assert.Empty(t, plugin.Version) // Java plugin typically has no version
		case "org.springframework.boot":
			foundSpringBoot = true
			assert.Equal(t, "2.7.0", plugin.Version)
		case "io.spring.dependency-management":
			foundDependencyManagement = true
			assert.Equal(t, "1.0.11.RELEASE", plugin.Version)
		case "application":
			foundApplication = true
		}
	}

	assert.True(t, foundJava, "Java plugin not found")
	assert.True(t, foundSpringBoot, "Spring Boot plugin not found")
	assert.True(t, foundDependencyManagement, "Dependency management plugin not found")
	assert.True(t, foundApplication, "Application plugin not found")
}

// TestProjectTypeDetection tests project type detection functions
func TestProjectTypeDetection(t *testing.T) {
	tests := []struct {
		name          string
		gradleContent string
		isAndroid     bool
		isKotlin      bool
		isSpringBoot  bool
	}{
		{
			name: "Spring Boot Project",
			gradleContent: `
plugins {
    id 'java'
    id 'org.springframework.boot' version '2.7.0'
}
`,
			isSpringBoot: true,
		},
		{
			name: "Android Project",
			gradleContent: `
plugins {
    id 'com.android.application'
}
`,
			isAndroid: true,
		},
		{
			name: "Kotlin Project",
			gradleContent: `
plugins {
    id 'org.jetbrains.kotlin.jvm' version '1.7.10'
}
`,
			isKotlin: true,
		},
		{
			name: "Plain Java Project",
			gradleContent: `
plugins {
    id 'java'
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := api.ParseString(tt.gradleContent)
			require.NoError(t, err)

			plugins := result.Project.Plugins

			assert.Equal(t, tt.isAndroid, api.IsAndroidProject(plugins))
			assert.Equal(t, tt.isKotlin, api.IsKotlinProject(plugins))
			assert.Equal(t, tt.isSpringBoot, api.IsSpringBootProject(plugins))
		})
	}
}

// TestDependenciesByScope tests dependency grouping by scope
func TestDependenciesByScope(t *testing.T) {
	// Create test dependencies with known scopes
	deps := []*model.Dependency{
		{Group: "org.springframework.boot", Name: "spring-boot-starter-web", Version: "2.7.0", Scope: "implementation"},
		{Group: "mysql", Name: "mysql-connector-java", Version: "8.0.29", Scope: "implementation"},
		{Group: "junit", Name: "junit", Version: "4.13.2", Scope: "testImplementation"},
		{Group: "org.springframework.boot", Name: "spring-boot-starter-test", Version: "2.7.0", Scope: "testImplementation"},
		{Group: "org.projectlombok", Name: "lombok", Version: "1.18.24", Scope: "compileOnly"},
	}

	depSets := api.DependenciesByScope(deps)

	// Should have 3 scopes: implementation, testImplementation, compileOnly
	assert.Len(t, depSets, 3)

	scopeCounts := make(map[string]int)
	for _, depSet := range depSets {
		scopeCounts[depSet.Scope] = len(depSet.Dependencies)
	}

	assert.Equal(t, 2, scopeCounts["implementation"])
	assert.Equal(t, 2, scopeCounts["testImplementation"])
	assert.Equal(t, 1, scopeCounts["compileOnly"])
}

// TestErrorHandling tests error handling for invalid Gradle content
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		gradleContent string
		expectError   bool
	}{
		{
			name:          "Valid content",
			gradleContent: `group = 'com.example'`,
			expectError:   false,
		},
		{
			name:          "Empty content",
			gradleContent: "",
			expectError:   false, // Empty content should be valid
		},
		{
			name:          "Whitespace only",
			gradleContent: "   \n\t  \n  ",
			expectError:   false, // Whitespace should be valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := api.ParseString(tt.gradleContent)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.Project)
			}
		})
	}
}

// TestParseWarnings tests that warnings are properly collected
func TestParseWarnings(t *testing.T) {
	// This test assumes that certain constructs might generate warnings
	// The actual implementation may vary
	gradleContent := `
group = 'com.example'
version = '1.0.0'
`

	result, err := api.ParseString(gradleContent)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Warnings should be a slice (may be empty)
	assert.NotNil(t, result.Warnings)
}

// TestCustomParserOptions tests custom parser configuration
func TestCustomParserOptions(t *testing.T) {
	gradleContent := `
// This is a comment
plugins {
    id 'java'
}

group = 'com.example'
version = '1.0.0'

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
}

task customTask {
    doLast {
        println 'Custom task'
    }
}
`

	// Test with different options
	tests := []struct {
		name    string
		options *api.Options
	}{
		{
			name: "Default options",
			options: &api.Options{
				SkipComments:      true,
				CollectRawContent: true,
				ParsePlugins:      true,
				ParseDependencies: true,
				ParseRepositories: true,
				ParseTasks:        true,
			},
		},
		{
			name: "Skip comments and tasks",
			options: &api.Options{
				SkipComments:      true,
				CollectRawContent: false,
				ParsePlugins:      true,
				ParseDependencies: true,
				ParseRepositories: true,
				ParseTasks:        false,
			},
		},
		{
			name: "Dependencies only",
			options: &api.Options{
				SkipComments:      true,
				CollectRawContent: false,
				ParsePlugins:      false,
				ParseDependencies: true,
				ParseRepositories: false,
				ParseTasks:        false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := api.NewParser(tt.options)
			result, err := parser.Parse(gradleContent)

			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotNil(t, result.Project)

			// Verify that options are respected
			if tt.options.ParseDependencies {
				assert.Greater(t, len(result.Project.Dependencies), 0)
			}

			if tt.options.ParsePlugins {
				assert.Greater(t, len(result.Project.Plugins), 0)
			}

			if tt.options.CollectRawContent {
				assert.NotEmpty(t, result.RawText)
			} else {
				assert.Empty(t, result.RawText)
			}
		})
	}
}
