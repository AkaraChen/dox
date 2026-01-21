package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetComposeBuilder(t *testing.T) {
	// Change to test fixture directory
	fixtureDir := filepath.Join("..", "test", "fixtures", "simple")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	builder, err := getComposeBuilder()
	assert.NoError(t, err)
	assert.NotNil(t, builder)
}

func TestGetComposeBuilder_WithProfile(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "with-profiles")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	// Set profile flag
	profile = "dev"

	builder, err := getComposeBuilder()
	assert.NoError(t, err)
	assert.NotNil(t, builder)

	// Reset profile
	profile = ""
}

func TestGetComposeBuilder_NoComposeFiles(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	builder, err := getComposeBuilder()
	// Should still succeed, just with no files
	assert.NoError(t, err)
	assert.NotNil(t, builder)
}

func TestGetConfig(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "simple")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	cfg, err := getConfig()
	assert.NoError(t, err)
	// simple fixture has no do.yaml
	assert.Nil(t, cfg)
}

func TestGetConfig_WithDoYaml(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "with-profiles")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	cfg, err := getConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Profiles)
}

func TestGetComposeExecutor(t *testing.T) {
	executor := getComposeExecutor()
	assert.NotNil(t, executor)
}

func TestGetComposeExecutor_WithDryRun(t *testing.T) {
	originalDryRun := dryRun
	defer func() { dryRun = originalDryRun }()

	dryRun = true
	executor := getComposeExecutor()
	assert.NotNil(t, executor)
}

func TestGetComposeExecutor_WithVerbose(t *testing.T) {
	originalVerbose := verbose
	defer func() { verbose = originalVerbose }()

	verbose = true
	executor := getComposeExecutor()
	assert.NotNil(t, executor)
	assert.Contains(t, executor.Env, "DOCKER_COMPOSE_VERBOSE=1")
}

func TestParseHookCommand(t *testing.T) {
	tests := []struct {
		name     string
		hook     string
		expected []string
	}{
		{
			name:     "simple command",
			hook:     "echo hello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "command with flags",
			hook:     "docker network prune -f",
			expected: []string{"docker", "network", "prune", "-f"},
		},
		{
			name:     "empty hook",
			hook:     "",
			expected: []string{},
		},
		{
			name:     "single word",
			hook:     "restart",
			expected: []string{"restart"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHookCommand(tt.hook)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPrintCommand(t *testing.T) {
	// Just verify it doesn't panic
	printCommand("docker compose up")
}

func TestResolveFile_AbsolutePath(t *testing.T) {
	input := "/absolute/path/to/compose.yaml"
	result, err := resolveFile(input)
	assert.NoError(t, err)
	assert.Equal(t, input, result)
}

func TestResolveFile_RelativePath(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	input := "compose.yaml"
	result, err := resolveFile(input)
	assert.NoError(t, err)
	assert.True(t, filepath.IsAbs(result))
	assert.Contains(t, result, "compose.yaml")
}

func TestResolveFile_CurrentDirectoryError(t *testing.T) {
	// We can't easily mock os.Getwd, so we skip this test
	// In real scenarios, os.Getwd rarely fails
	t.Skip("Cannot mock os.Getwd safely")
}

func TestIsKnownCommand(t *testing.T) {
	knownCommands := []string{
		"up", "down", "ps", "logs", "restart", "exec", "build",
		"pull", "push", "start", "stop", "rm", "kill", "run",
		"pause", "unpause", "top", "events", "port", "config",
		"create", "version",
	}

	for _, cmd := range knownCommands {
		t.Run(cmd, func(t *testing.T) {
			assert.True(t, isKnownCommand(cmd))
		})
	}
}

func TestIsKnownCommand_Unknown(t *testing.T) {
	unknownCommands := []string{
		"unknown", "invalid", "fake", "custom",
	}

	for _, cmd := range unknownCommands {
		t.Run(cmd, func(t *testing.T) {
			assert.False(t, isKnownCommand(cmd))
		})
	}
}

func TestIsKnownCommand_Empty(t *testing.T) {
	assert.False(t, isKnownCommand(""))
}

func TestResolveAlias_EmptyError(t *testing.T) {
	_, err := resolveAlias("")
	assert.Error(t, err)
}

func TestResolveAlias_SimpleCommandInDir(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "simple")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	commands, err := resolveAlias("up -d")
	assert.NoError(t, err)
	assert.Len(t, commands, 1)
	assert.Contains(t, commands[0], "up")
}

func TestResolveAlias_ChainedCommandsInDir(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "simple")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	commands, err := resolveAlias("down -v && up -d")
	assert.NoError(t, err)
	assert.Len(t, commands, 2)
}

func TestResolveAlias_UnknownCommand(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "simple")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	commands, err := resolveAlias("custom-command -f")
	assert.NoError(t, err)
	assert.Len(t, commands, 1)
	assert.Equal(t, []string{"custom-command", "-f"}, commands[0])
}

func TestResolveAlias_WithWhitespace(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "simple")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	commands, err := resolveAlias("  down   &&  up  ")
	assert.NoError(t, err)
	assert.Len(t, commands, 2)
}

func TestExecuteHooks_NoConfig(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	// No hooks defined
	err = executeHooks("pre_up")
	assert.NoError(t, err)
}

func TestExecuteHooks_NilConfig(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "simple")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	// simple has no do.yaml
	err = executeHooks("pre_up")
	assert.NoError(t, err)
}

func TestExecuteHooks_WithDryRun(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "with-hooks")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	originalDryRun := dryRun
	defer func() { dryRun = originalDryRun }()

	dryRun = true
	err = executeHooks("pre_up")
	assert.NoError(t, err) // Dry run skips execution
}

func TestExecuteHooks_HookTypeNotFound(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "with-hooks")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	// Hook type that doesn't exist
	err = executeHooks("nonexistent_hook")
	assert.NoError(t, err)
}

func TestResolveAlias_MultipleKnownCommands(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "simple")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	commands, err := resolveAlias("ps && logs && restart")
	assert.NoError(t, err)
	assert.Len(t, commands, 3)
}

func TestResolveAlias_MixedCommands(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "simple")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	commands, err := resolveAlias("echo 'test' && up -d && echo 'done'")
	assert.NoError(t, err)
	assert.Len(t, commands, 3)
}
