package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveRemoteProject(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	// Create test config with project paths
	content := []byte(`
projects:
  webapp:
    path: /home/user/projects/webapp
    description: "Web application"
  api:
    path: /home/user/projects/api
    description: "API service"
`)
	err := os.WriteFile(configFile, content, 0644)
	require.NoError(t, err)

	cfg, err := LoadGlobalConfig(configFile)
	require.NoError(t, err)

	tests := []struct {
		name           string
		input          string
		shouldResolve  bool
		projectName    string
		projectPath    string
		remainingCmd   string
	}{
		{
			name:          "simple project reference",
			input:         "@webapp",
			shouldResolve: true,
			projectName:   "webapp",
			projectPath:   "/home/user/projects/webapp",
			remainingCmd:  "",
		},
		{
			name:          "project with command",
			input:         "@api c up",
			shouldResolve: true,
			projectName:   "api",
			projectPath:   "/home/user/projects/api",
			remainingCmd:  "c up",
		},
		{
			name:          "project with full command",
			input:         "@webapp compose logs -f",
			shouldResolve: true,
			projectName:   "webapp",
			projectPath:   "/home/user/projects/webapp",
			remainingCmd:  "compose logs -f",
		},
		{
			name:           "non-existent project",
			input:          "@nonexistent c up",
			shouldResolve:  false,
			projectName:    "nonexistent",
			projectPath:    "",
			remainingCmd:   "",
		},
		{
			name:           "not a project reference",
			input:          "c up",
			shouldResolve:  false,
			projectName:    "",
			projectPath:    "",
			remainingCmd:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cfg.ResolveRemoteProject(tt.input)
			if tt.shouldResolve {
				require.NoError(t, err)
				assert.Equal(t, tt.projectName, result.ProjectName)
				assert.Equal(t, tt.projectPath, result.ProjectPath)
				assert.Equal(t, tt.remainingCmd, result.RemainingCommand)
			} else {
				if tt.projectName == "" {
					// Not a remote reference
					assert.Nil(t, result)
					assert.NoError(t, err)
				} else {
					// Invalid project reference
					assert.Error(t, err)
				}
			}
		})
	}
}

func TestRemoteResolution_InvalidConfig(t *testing.T) {
	cfg := &GlobalConfig{
		Projects: map[string]ProjectEntry{},
	}

	// Try to resolve when no projects are configured
	result, err := cfg.ResolveRemoteProject("@missing c up")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRemoteResolution_ProjectPathDoesNotExist(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	// Create config with non-existent path
	content := []byte(`
projects:
  ghost:
    path: /nonexistent/path/that/does/not/exist
`)
	err := os.WriteFile(configFile, content, 0644)
	require.NoError(t, err)

	cfg, err := LoadGlobalConfig(configFile)
	require.NoError(t, err)

	// Should still resolve, even if path doesn't exist
	// (validation happens at execution time)
	result, err := cfg.ResolveRemoteProject("@ghost")
	require.NoError(t, err)
	assert.Equal(t, "ghost", result.ProjectName)
	assert.Equal(t, "/nonexistent/path/that/does/not/exist", result.ProjectPath)
}

func TestRemoteResolution_WithWhitespace(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	content := []byte(`
projects:
  myproject:
    path: /home/user/projects/myproject
`)
	err := os.WriteFile(configFile, content, 0644)
	require.NoError(t, err)

	cfg, err := LoadGlobalConfig(configFile)
	require.NoError(t, err)

	// Test with extra whitespace
	result, err := cfg.ResolveRemoteProject("@myproject    c    up")
	require.NoError(t, err)
	assert.Equal(t, "myproject", result.ProjectName)
	assert.Equal(t, "c    up", result.RemainingCommand)
}

func TestRemoteResolution_ProjectNameWithUnderscore(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	content := []byte(`
projects:
  my_project:
    path: /home/user/projects/my_project
`)
	err := os.WriteFile(configFile, content, 0644)
	require.NoError(t, err)

	cfg, err := LoadGlobalConfig(configFile)
	require.NoError(t, err)

	result, err := cfg.ResolveRemoteProject("@my_project c up")
	require.NoError(t, err)
	assert.Equal(t, "my_project", result.ProjectName)
}

func TestRemoteResolution_ProjectNameWithDash(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	content := []byte(`
projects:
  my-project:
    path: /home/user/projects/my-project
`)
	err := os.WriteFile(configFile, content, 0644)
	require.NoError(t, err)

	cfg, err := LoadGlobalConfig(configFile)
	require.NoError(t, err)

	result, err := cfg.ResolveRemoteProject("@my-project c up")
	require.NoError(t, err)
	assert.Equal(t, "my-project", result.ProjectName)
}

func TestListProjects(t *testing.T) {
	cfg := &GlobalConfig{
		Projects: map[string]ProjectEntry{
			"webapp": {Path: "/home/user/webapp", Description: "Web app"},
			"api":    {Path: "/home/user/api", Description: "API"},
		},
	}

	projects := cfg.ListProjects()
	assert.Len(t, projects, 2)

	// Check that we get the project info
	names := make(map[string]bool)
	for _, p := range projects {
		names[p.Name] = true
	}
	assert.True(t, names["webapp"])
	assert.True(t, names["api"])
}
