package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_ValidYAML(t *testing.T) {
	content := `
version: 1
defaults:
  profile: dev
profiles:
  dev:
    slices: [dev]
  prod:
    slices: [prod]
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dox.yaml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, 1, config.Version)
	assert.Equal(t, "dev", config.Defaults.Profile)
	assert.Len(t, config.Profiles, 2)
	assert.Contains(t, config.Profiles, "dev")
	assert.Contains(t, config.Profiles, "prod")
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	content := `
version: 1
  invalid: yaml: syntax
    badly: indented
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dox.yaml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	_, err = LoadConfig(configPath)
	assert.Error(t, err)
}

func TestLoadConfig_MinimalConfig(t *testing.T) {
	content := `version: 1`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dox.yaml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, 1, config.Version)
	assert.Empty(t, config.Profiles)
	assert.Empty(t, config.Aliases)
}

func TestLoadConfig_WithAliases(t *testing.T) {
	content := `
version: 1
aliases:
  fresh: "down -v && up --build -d"
  restart-all: "down && up -d"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dox.yaml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Len(t, config.Aliases, 2)
	assert.Equal(t, "down -v && up --build -d", config.Aliases["fresh"])
	assert.Equal(t, "down && up -d", config.Aliases["restart-all"])
}

func TestLoadConfig_WithHooks(t *testing.T) {
	content := `
version: 1
hooks:
  pre_up:
    - echo "Starting..."
    - echo "Check..."
  post_up:
    - echo "Ready!"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dox.yaml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Len(t, config.Hooks, 2)
	assert.Len(t, config.Hooks["pre_up"], 2)
	assert.Equal(t, "echo \"Starting...\"", config.Hooks["pre_up"][0])
	assert.Len(t, config.Hooks["post_up"], 1)
}

func TestLoadConfig_WithEnvFiles(t *testing.T) {
	content := `
version: 1
env_files:
  dev: .env.dev
  prod: .env.prod
profiles:
  dev:
    slices: [dev]
    env: dev
  prod:
    slices: [prod]
    env: prod
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dox.yaml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Len(t, config.EnvFiles, 2)
	assert.Equal(t, ".env.dev", config.EnvFiles["dev"])
	assert.Equal(t, ".env.prod", config.EnvFiles["prod"])
}

func TestLoadConfig_WithDiscoveryConfig(t *testing.T) {
	content := `
version: 1
discovery:
  enabled: true
  pattern: "docker-compose.*.yml"
  base: "docker-compose.yml"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dox.yaml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.True(t, config.Discovery.Enabled)
	assert.Equal(t, "docker-compose.*.yml", config.Discovery.Pattern)
	assert.Equal(t, "docker-compose.yml", config.Discovery.Base)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent.yaml")

	_, err := LoadConfig(configPath)
	assert.Error(t, err)
}

func TestFindConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dox.yaml")
	os.WriteFile(configPath, []byte("version: 1"), 0644)

	found, err := FindConfigFile(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, configPath, found)
}

func TestFindConfigFile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	found, err := FindConfigFile(tmpDir)
	assert.Error(t, err)
	assert.Empty(t, found)
}

func TestLoadConfigFromDirectory_WithValidConfig(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "test", "fixtures", "with-profiles")

	config, configPath, err := LoadConfigFromDirectory(fixtureDir)
	require.NoError(t, err)
	assert.NotEmpty(t, configPath)
	assert.NotNil(t, config)
	assert.Equal(t, 1, config.Version)
	assert.Equal(t, "dev", config.Defaults.Profile)
}

func TestLoadConfigFromDirectory_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()

	config, configPath, err := LoadConfigFromDirectory(tmpDir)
	assert.NoError(t, err) // No error, just no config
	assert.Empty(t, configPath)
	assert.Nil(t, config)
}
