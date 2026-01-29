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
	expected := filepath.Join("/test/home", ".config", "dox", "config.yaml")
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

func TestGetAlias(t *testing.T) {
	cfg := &GlobalConfig{
		Aliases: map[string]string{
			"refresh": "down && up --build -d",
			"clean":   "down -v --remove-orphans",
		},
	}

	// Test existing alias
	alias, found := cfg.GetAlias("refresh")
	assert.True(t, found)
	assert.Equal(t, "down && up --build -d", alias)

	// Test non-existent alias
	alias, found = cfg.GetAlias("nonexistent")
	assert.False(t, found)
	assert.Empty(t, alias)
}

func TestGetAlias_EmptyAliases(t *testing.T) {
	cfg := &GlobalConfig{
		Aliases: map[string]string{},
	}

	alias, found := cfg.GetAlias("anything")
	assert.False(t, found)
	assert.Empty(t, alias)
}

func TestGetAlias_NilAliases(t *testing.T) {
	cfg := &GlobalConfig{
		Aliases: nil,
	}

	alias, found := cfg.GetAlias("anything")
	assert.False(t, found)
	assert.Empty(t, alias)
}

func TestHasProject(t *testing.T) {
	cfg := &GlobalConfig{
		Projects: map[string]ProjectEntry{
			"webapp": {Path: "/home/user/webapp"},
			"api":    {Path: "/home/user/api"},
		},
	}

	assert.True(t, cfg.HasProject("webapp"))
	assert.True(t, cfg.HasProject("api"))
	assert.False(t, cfg.HasProject("nonexistent"))
}

func TestHasProject_EmptyProjects(t *testing.T) {
	cfg := &GlobalConfig{
		Projects: map[string]ProjectEntry{},
	}

	assert.False(t, cfg.HasProject("anything"))
}

func TestHasProject_NilProjects(t *testing.T) {
	cfg := &GlobalConfig{
		Projects: nil,
	}

	assert.False(t, cfg.HasProject("anything"))
}

func TestProjectNames(t *testing.T) {
	cfg := &GlobalConfig{
		Projects: map[string]ProjectEntry{
			"webapp": {Path: "/home/user/webapp"},
			"api":    {Path: "/home/user/api"},
			"db":     {Path: "/home/user/db"},
		},
	}

	names := cfg.ProjectNames()
	assert.Len(t, names, 3)

	// Convert to set for easier checking
	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}
	assert.True(t, nameSet["webapp"])
	assert.True(t, nameSet["api"])
	assert.True(t, nameSet["db"])
}

func TestProjectNames_Empty(t *testing.T) {
	cfg := &GlobalConfig{
		Projects: map[string]ProjectEntry{},
	}

	names := cfg.ProjectNames()
	assert.Empty(t, names)
}

func TestProjectNames_Nil(t *testing.T) {
	cfg := &GlobalConfig{
		Projects: nil,
	}

	names := cfg.ProjectNames()
	assert.Empty(t, names)
}

func TestAliasNames(t *testing.T) {
	cfg := &GlobalConfig{
		Aliases: map[string]string{
			"refresh": "down && up",
			"clean":   "down -v",
			"rebuild": "up --build",
		},
	}

	names := cfg.AliasNames()
	assert.Len(t, names, 3)

	// Convert to set for easier checking
	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}
	assert.True(t, nameSet["refresh"])
	assert.True(t, nameSet["clean"])
	assert.True(t, nameSet["rebuild"])
}

func TestAliasNames_Empty(t *testing.T) {
	cfg := &GlobalConfig{
		Aliases: map[string]string{},
	}

	names := cfg.AliasNames()
	assert.Empty(t, names)
}

func TestAliasNames_Nil(t *testing.T) {
	cfg := &GlobalConfig{
		Aliases: nil,
	}

	names := cfg.AliasNames()
	assert.Empty(t, names)
}

func TestIsAtProjectReference(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid project ref", "@myproject", true},
		{"valid with command", "@myproject c up", true},
		{"no at sign", "myproject c up", false},
		{"at only", "@", false},
		{"at in middle", "my@project", false},
		{"empty string", "", false},
		{"at with space", "@ myproject", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAtProjectReference(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
