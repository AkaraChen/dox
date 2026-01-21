package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadGlobalConfig(t *testing.T) {
	// Create temp directory for global config
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	// Create test config
	content := []byte(`
projects:
  myproject:
    path: /home/user/projects/myproject
    description: "My main project"
  another:
    path: /home/user/projects/another
    description: "Another project"
`)
	err := os.WriteFile(configFile, content, 0644)
	require.NoError(t, err)

	// Load config
	cfg, err := LoadGlobalConfig(configFile)
	require.NoError(t, err)

	// Verify projects
	assert.Len(t, cfg.Projects, 2)
	assert.Equal(t, "/home/user/projects/myproject", cfg.Projects["myproject"].Path)
	assert.Equal(t, "My main project", cfg.Projects["myproject"].Description)
	assert.Equal(t, "/home/user/projects/another", cfg.Projects["another"].Path)
}

func TestLoadGlobalConfig_NotFound(t *testing.T) {
	// Try to load non-existent config
	cfg, err := LoadGlobalConfig("/nonexistent/config.yaml")
	require.NoError(t, err)
	assert.Nil(t, cfg)
}

func TestLoadGlobalConfig_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	// Create invalid YAML
	content := []byte(`
projects:
  myproject:
    path: [invalid yaml
`)
	err := os.WriteFile(configFile, content, 0644)
	require.NoError(t, err)

	// Load should fail gracefully
	cfg, err := LoadGlobalConfig(configFile)
	require.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadGlobalConfig_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	// Create empty file
	err := os.WriteFile(configFile, []byte{}, 0644)
	require.NoError(t, err)

	// Load should return empty config
	cfg, err := LoadGlobalConfig(configFile)
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Empty(t, cfg.Projects)
}

func TestLoadGlobalConfig_OnlyAliases(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	// Create config with only aliases, no projects
	content := []byte(`
aliases:
  refresh: "down && up --build -d"
  clean: "down -v --remove-orphans"
`)
	err := os.WriteFile(configFile, content, 0644)
	require.NoError(t, err)

	// Load config
	cfg, err := LoadGlobalConfig(configFile)
	require.NoError(t, err)

	// Verify aliases
	assert.Len(t, cfg.Aliases, 2)
	assert.Equal(t, "down && up --build -d", cfg.Aliases["refresh"])
	assert.Equal(t, "down -v --remove-orphans", cfg.Aliases["clean"])
}

func TestGetGlobalConfigPath(t *testing.T) {
	// Test default path
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	os.Setenv("HOME", "/test/home")
	path := GetGlobalConfigPath()
	expected := filepath.Join("/test/home", ".config", "do", "config.yaml")
	assert.Equal(t, expected, path)
}

func TestResolveProjectPath(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	// Create test config
	content := []byte(`
projects:
  myproject:
    path: /home/user/projects/myproject
  another:
    path: /home/user/projects/another
`)
	err := os.WriteFile(configFile, content, 0644)
	require.NoError(t, err)

	// Load config
	cfg, err := LoadGlobalConfig(configFile)
	require.NoError(t, err)

	// Resolve existing project
	path, found := cfg.ResolveProjectPath("myproject")
	assert.True(t, found)
	assert.Equal(t, "/home/user/projects/myproject", path)

	// Resolve non-existent project
	path, found = cfg.ResolveProjectPath("nonexistent")
	assert.False(t, found)
	assert.Empty(t, path)
}

func TestGlobalConfig_LoadOrDefault(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Test with non-existent config
	cfg, err := LoadGlobalConfigOrDefault(filepath.Join(tempDir, "nonexistent.yaml"))
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Empty(t, cfg.Projects)

	// Test with valid config
	configFile := filepath.Join(tempDir, "config.yaml")
	content := []byte(`
projects:
  test:
    path: /test/path
`)
	err = os.WriteFile(configFile, content, 0644)
	require.NoError(t, err)

	cfg, err = LoadGlobalConfigOrDefault(configFile)
	require.NoError(t, err)
	assert.Len(t, cfg.Projects, 1)
}

func TestParseAtProjectReference(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		isAtRef  bool
		project  string
		rest     string
	}{
		{
			name:    "simple project reference",
			input:   "@myproject",
			isAtRef: true,
			project: "myproject",
			rest:    "",
		},
		{
			name:    "project reference with command",
			input:   "@myproject c up",
			isAtRef: true,
			project: "myproject",
			rest:    "c up",
		},
		{
			name:    "project reference with full command",
			input:   "@myproject compose up -d",
			isAtRef: true,
			project: "myproject",
			rest:    "compose up -d",
		},
		{
			name:    "not an at reference",
			input:   "myproject c up",
			isAtRef: false,
			project: "",
			rest:    "",
		},
		{
			name:    "at sign but no project name",
			input:   "@ c up",
			isAtRef: false,
			project: "",
			rest:    "",
		},
		{
			name:    "at sign in middle",
			input:   "myproject@c up",
			isAtRef: false,
			project: "",
			rest:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isAtRef, project, rest := ParseAtProjectReference(tt.input)
			assert.Equal(t, tt.isAtRef, isAtRef)
			assert.Equal(t, tt.project, project)
			assert.Equal(t, tt.rest, rest)
		})
	}
}
