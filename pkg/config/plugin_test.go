package config

import (
	"testing"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

func TestNewPluginParser(t *testing.T) {
	parser := NewPluginParser()
	if parser == nil {
		t.Error("NewPluginParser() returned nil")
	}
}

func TestParsePluginBlock(t *testing.T) {
	parser := NewPluginParser()

	// Test with nil block。
	_, err := parser.ParsePluginBlock(nil)
	if err == nil {
		t.Error("ParsePluginBlock() should return error for nil block")
	}

	// Test with plugin declarations。
	block := &model.ScriptBlock{
		Name: "plugins",
		Values: map[string]interface{}{
			"id 'java'": "id 'java'",
			"id 'com.android.application' version '7.0.0'":            "id 'com.android.application' version '7.0.0'",
			"id(\"org.jetbrains.kotlin.android\") version \"1.5.30\"": "id(\"org.jetbrains.kotlin.android\") version \"1.5.30\"",
		},
		Closures: make(map[string][]*model.ScriptBlock),
	}

	plugins, err := parser.ParsePluginBlock(block)
	if err != nil {
		t.Fatalf("ParsePluginBlock() error = %v", err)
	}

	if len(plugins) != 3 {
		t.Errorf("ParsePluginBlock() returned %v plugins, want 3", len(plugins))
	}

	// Verify that plugin IDs and versions are correctly parsed。
	var foundJava, foundAndroid, foundKotlin bool
	for _, plugin := range plugins {
		switch plugin.ID {
		case "java":
			foundJava = true
			if plugin.Version != "" {
				t.Errorf("java plugin should have no version, got %s", plugin.Version)
			}
		case "com.android.application":
			foundAndroid = true
			if plugin.Version != "7.0.0" {
				t.Errorf("android plugin should have version 7.0.0, got %s", plugin.Version)
			}
		case "org.jetbrains.kotlin.android":
			foundKotlin = true
			if plugin.Version != "1.5.30" {
				t.Errorf("kotlin plugin should have version 1.5.30, got %s", plugin.Version)
			}
		}
	}

	if !foundJava {
		t.Error("ParsePluginBlock() did not find java plugin")
	}
	if !foundAndroid {
		t.Error("ParsePluginBlock() did not find android plugin")
	}
	if !foundKotlin {
		t.Error("ParsePluginBlock() did not find kotlin plugin")
	}

	// Test with id closures。
	idClosureBlock := &model.ScriptBlock{
		Values: map[string]interface{}{
			"spring-boot": "spring-boot",
		},
	}

	block.Closures["id"] = []*model.ScriptBlock{idClosureBlock}

	plugins, err = parser.ParsePluginBlock(block)
	if err != nil {
		t.Fatalf("ParsePluginBlock() error = %v", err)
	}

	if len(plugins) != 4 {
		t.Errorf("ParsePluginBlock() with id closure returned %v plugins, want 4", len(plugins))
	}
}

func TestExtractPluginsFromText(t *testing.T) {
	parser := NewPluginParser()

	// Test with empty text。
	plugins := parser.ExtractPluginsFromText("")
	if len(plugins) != 0 {
		t.Errorf("ExtractPluginsFromText() with empty text returned %v plugins, want 0", len(plugins))
	}

	// Test with plugin declarations and apply statements。
	text := `plugins {
		id 'java'
		id 'com.android.application' version '7.0.0'
	}
	
	apply plugin: 'kotlin'`

	plugins = parser.ExtractPluginsFromText(text)
	if len(plugins) != 3 {
		t.Errorf("ExtractPluginsFromText() returned %v plugins, want 3", len(plugins))
	}

	// Verify the plugins。
	var foundJava, foundAndroid, foundKotlin bool
	for _, plugin := range plugins {
		switch plugin.ID {
		case "java":
			foundJava = true
		case "com.android.application":
			foundAndroid = true
			if plugin.Version != "7.0.0" {
				t.Errorf("android plugin should have version 7.0.0, got %s", plugin.Version)
			}
		case "kotlin":
			foundKotlin = true
		}
	}

	if !foundJava {
		t.Error("ExtractPluginsFromText() did not find java plugin")
	}
	if !foundAndroid {
		t.Error("ExtractPluginsFromText() did not find android plugin")
	}
	if !foundKotlin {
		t.Error("ExtractPluginsFromText() did not find kotlin plugin")
	}
}

func TestGetPluginConfigurations(t *testing.T) {
	parser := NewPluginParser()

	// Create a root block with some plugin configurations。
	rootBlock := &model.ScriptBlock{
		Name: "root",
		Closures: map[string][]*model.ScriptBlock{
			"android": {
				{
					Name: "android",
					Values: map[string]interface{}{
						"compileSdkVersion": "31",
					},
				},
			},
			"kotlin": {
				{
					Name: "kotlin",
					Values: map[string]interface{}{
						"jvmTarget": "1.8",
					},
				},
			},
		},
	}

	// Create a list of plugins。
	plugins := []*model.Plugin{
		{ID: "com.android.application"},
		{ID: "kotlin"},
		{ID: "unknown-plugin"},
	}

	configs := parser.GetPluginConfigurations(rootBlock, plugins)
	if len(configs) != 2 {
		t.Errorf("GetPluginConfigurations() returned %v configs, want 2", len(configs))
	}

	// Check for android config。
	androidConfig, hasAndroid := configs["com.android.application"]
	if !hasAndroid {
		t.Error("GetPluginConfigurations() did not find android config")
	} else if androidConfig.Name != "android" {
		t.Errorf("Android config has wrong name, got %s", androidConfig.Name)
	}

	// Check for kotlin config。
	kotlinConfig, hasKotlin := configs["kotlin"]
	if !hasKotlin {
		t.Error("GetPluginConfigurations() did not find kotlin config")
	} else if kotlinConfig.Name != "kotlin" {
		t.Errorf("Kotlin config has wrong name, got %s", kotlinConfig.Name)
	}
}

func TestIsAndroidProject(t *testing.T) {
	parser := NewPluginParser()

	tests := []struct {
		name    string
		plugins []*model.Plugin
		want    bool
	}{
		{
			name:    "empty list",
			plugins: []*model.Plugin{},
			want:    false,
		},
		{
			name: "no android plugin",
			plugins: []*model.Plugin{
				{ID: "java"},
				{ID: "kotlin"},
			},
			want: false,
		},
		{
			name: "has android application plugin",
			plugins: []*model.Plugin{
				{ID: "com.android.application"},
			},
			want: true,
		},
		{
			name: "has android library plugin",
			plugins: []*model.Plugin{
				{ID: "com.android.library"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.IsAndroidProject(tt.plugins); got != tt.want {
				t.Errorf("IsAndroidProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSpringBootProject(t *testing.T) {
	parser := NewPluginParser()

	tests := []struct {
		name    string
		plugins []*model.Plugin
		want    bool
	}{
		{
			name:    "empty list",
			plugins: []*model.Plugin{},
			want:    false,
		},
		{
			name: "no spring boot plugin",
			plugins: []*model.Plugin{
				{ID: "java"},
				{ID: "kotlin"},
			},
			want: false,
		},
		{
			name: "has spring boot plugin",
			plugins: []*model.Plugin{
				{ID: "org.springframework.boot"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.IsSpringBootProject(tt.plugins); got != tt.want {
				t.Errorf("IsSpringBootProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsKotlinProject(t *testing.T) {
	parser := NewPluginParser()

	tests := []struct {
		name    string
		plugins []*model.Plugin
		want    bool
	}{
		{
			name:    "empty list",
			plugins: []*model.Plugin{},
			want:    false,
		},
		{
			name: "no kotlin plugin",
			plugins: []*model.Plugin{
				{ID: "java"},
			},
			want: false,
		},
		{
			name: "has kotlin plugin",
			plugins: []*model.Plugin{
				{ID: "kotlin"},
			},
			want: true,
		},
		{
			name: "has kotlin jvm plugin",
			plugins: []*model.Plugin{
				{ID: "org.jetbrains.kotlin.jvm"},
			},
			want: true,
		},
		{
			name: "has kotlin android plugin",
			plugins: []*model.Plugin{
				{ID: "org.jetbrains.kotlin.android"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.IsKotlinProject(tt.plugins); got != tt.want {
				t.Errorf("IsKotlinProject() = %v, want %v", got, tt.want)
			}
		})
	}
}
