package config

import (
	"testing"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

func TestNewRepositoryParser(t *testing.T) {
	parser := NewRepositoryParser()
	if parser == nil {
		t.Error("NewRepositoryParser() returned nil")
	}
}

func TestParseRepositoryBlock(t *testing.T) {
	parser := NewRepositoryParser()

	// Test with nil block
	_, err := parser.ParseRepositoryBlock(nil)
	if err == nil {
		t.Error("ParseRepositoryBlock() should return error for nil block")
	}

	// Test with predefined repository names
	block := &model.ScriptBlock{
		Name: "repositories",
		Values: map[string]interface{}{
			"mavenCentral()": "mavenCentral()",
			"jcenter()":      "jcenter()",
		},
		Closures: make(map[string][]*model.ScriptBlock),
	}

	repos, err := parser.ParseRepositoryBlock(block)
	if err != nil {
		t.Fatalf("ParseRepositoryBlock() error = %v", err)
	}

	if len(repos) != 2 {
		t.Errorf("ParseRepositoryBlock() returned %v repositories, want 2", len(repos))
	}

	// Test with maven closures
	mavenBlock := &model.ScriptBlock{
		Values: map[string]interface{}{
			"url 'https://jitpack.io'": "url 'https://jitpack.io'",
		},
	}

	mavenWithCredentialsBlock := &model.ScriptBlock{
		Values: map[string]interface{}{
			"url 'https://example.com'": "url 'https://example.com'",
		},
		Closures: map[string][]*model.ScriptBlock{
			"credentials": {
				{
					Values: map[string]interface{}{
						"username": "'user'",
						"password": "'pass'",
					},
				},
			},
		},
	}

	block.Closures["maven"] = []*model.ScriptBlock{mavenBlock, mavenWithCredentialsBlock}

	repos, err = parser.ParseRepositoryBlock(block)
	if err != nil {
		t.Fatalf("ParseRepositoryBlock() error = %v", err)
	}

	if len(repos) != 4 {
		t.Errorf("ParseRepositoryBlock() returned %v repositories, want 4", len(repos))
	}

	// Check that we correctly parsed the URL
	var foundJitpack bool
	for _, repo := range repos {
		if repo.URL == "https://jitpack.io" {
			foundJitpack = true
			break
		}
	}

	if !foundJitpack {
		t.Error("ParseRepositoryBlock() did not correctly parse maven URL")
	}

	// Check that we correctly parsed credentials
	var foundCredentials bool
	for _, repo := range repos {
		if repo.URL == "https://example.com" && repo.Username == "user" && repo.Password == "pass" {
			foundCredentials = true
			break
		}
	}

	if !foundCredentials {
		t.Error("ParseRepositoryBlock() did not correctly parse maven credentials")
	}
}

func TestExtractRepositoriesFromText(t *testing.T) {
	parser := NewRepositoryParser()

	// Test with empty text
	repos := parser.ExtractRepositoriesFromText("")
	if len(repos) != 0 {
		t.Errorf("ExtractRepositoriesFromText() with empty text returned %v repositories, want 0", len(repos))
	}

	// Test with repositories block
	text := `repositories {
		mavenCentral()
		maven { url 'https://jitpack.io' }
	}`

	repos = parser.ExtractRepositoriesFromText(text)
	if len(repos) != 2 {
		t.Errorf("ExtractRepositoriesFromText() returned %v repositories, want 2", len(repos))
	}

	// Verify the repositories
	var foundMavenCentral, foundJitpack bool
	for _, repo := range repos {
		if repo.Name == "mavenCentral" {
			foundMavenCentral = true
		}
		if repo.URL == "https://jitpack.io" {
			foundJitpack = true
		}
	}

	if !foundMavenCentral {
		t.Error("ExtractRepositoriesFromText() did not find mavenCentral()")
	}
	if !foundJitpack {
		t.Error("ExtractRepositoriesFromText() did not find jitpack URL")
	}
}

func TestGetDefaultRepositories(t *testing.T) {
	parser := NewRepositoryParser()
	repos := parser.GetDefaultRepositories()

	if len(repos) != 4 {
		t.Errorf("GetDefaultRepositories() returned %v repositories, want 4", len(repos))
	}

	expectedRepos := map[string]bool{
		"mavenCentral": false,
		"mavenLocal":   false,
		"google":       false,
		"jcenter":      false,
	}

	for _, repo := range repos {
		expectedRepos[repo.Name] = true
	}

	for name, found := range expectedRepos {
		if !found {
			t.Errorf("GetDefaultRepositories() did not include %s", name)
		}
	}
}

func TestHasJitPackRepository(t *testing.T) {
	parser := NewRepositoryParser()

	tests := []struct {
		name  string
		repos []*model.Repository
		want  bool
	}{
		{
			name:  "empty list",
			repos: []*model.Repository{},
			want:  false,
		},
		{
			name: "no jitpack",
			repos: []*model.Repository{
				{Name: "mavenCentral", URL: "https://repo.maven.apache.org/maven2/"},
			},
			want: false,
		},
		{
			name: "has jitpack",
			repos: []*model.Repository{
				{Name: "jitpack", URL: "https://jitpack.io"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.HasJitPackRepository(tt.repos); got != tt.want {
				t.Errorf("HasJitPackRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasCustomRepository(t *testing.T) {
	parser := NewRepositoryParser()

	tests := []struct {
		name  string
		repos []*model.Repository
		want  bool
	}{
		{
			name:  "empty list",
			repos: []*model.Repository{},
			want:  false,
		},
		{
			name: "no custom repo",
			repos: []*model.Repository{
				{Name: "mavenCentral"},
				{Name: "google"},
			},
			want: false,
		},
		{
			name: "has custom repo",
			repos: []*model.Repository{
				{Name: "customRepo", URL: "https://custom.example.com"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.HasCustomRepository(tt.repos); got != tt.want {
				t.Errorf("HasCustomRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
