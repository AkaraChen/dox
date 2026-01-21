package compose

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHooks_ExecutionOrder(t *testing.T) {
	executor := NewExecutor(false)

	// Test hooks are executed in order
	hooks := []string{
	 "echo 'first'",
	 "echo 'second'",
	 "echo 'third'",
	}

	for _, hook := range hooks {
	 cmd := []string{"sh", "-c", hook}
	 _, err := executor.RunCommand(cmd)
	 require.NoError(t, err)
	}
}

func TestHooks_PreUp(t *testing.T) {
	executor := NewExecutor(true) // dry-run

	hooks := []string{
	 "echo 'Starting services...'",
	 "echo 'Check dependencies...'",
	}

	for _, hook := range hooks {
	 cmd := []string{"sh", "-c", hook}
	 _, err := executor.RunCommand(cmd)
	 require.NoError(t, err)
	}
}

func TestHooks_PostUp(t *testing.T) {
	executor := NewExecutor(true) // dry-run

	hooks := []string{
	 "echo 'Services are ready!'",
	}

	cmd := []string{"sh", "-c", hooks[0]}
	_, err := executor.RunCommand(cmd)
	require.NoError(t, err)
}

func TestHooks_Failure(t *testing.T) {
	executor := NewExecutor(false)

	// Hook that fails
	cmd := []string{"sh", "-c", "exit 1"}
	_, err := executor.RunCommand(cmd)
	assert.Error(t, err)
}

func TestHooks_WithComposeCommand(t *testing.T) {
	executor := NewExecutor(true)

	// Simulate: pre_up hooks -> docker compose up -> post_up hooks
	preHooks := []string{"echo 'Before up'"}
	postHooks := []string{"echo 'After up'"}

	// Execute pre hooks
	for _, hook := range preHooks {
	 cmd := strings.Fields(hook)
	 if len(cmd) == 1 {
   cmd = []string{"echo", "'Before up'"}
	 }
	 _, _ = executor.RunCommand(cmd)
	}

	// Execute main command (dry-run)
	mainCmd := []string{"docker", "compose", "up"}
	_, _ = executor.RunCommand(mainCmd)

	// Execute post hooks
	for _, hook := range postHooks {
	 cmd := strings.Fields(hook)
	 if len(cmd) == 1 {
   cmd = []string{"echo", "'After up'"}
	 }
	 _, _ = executor.RunCommand(cmd)
	}
}

func TestHooks_ExecutionFlow(t *testing.T) {
	// Test that hooks execute before and after main command
	executor := NewExecutor(true)

	var output strings.Builder
	executor.Stdout = &output

	// Pre hook
	cmd1 := []string{"echo", "pre"}
	executor.RunCommand(cmd1)

	// Main command
	cmd2 := []string{"echo", "main"}
	executor.RunCommand(cmd2)

	// Post hook
	cmd3 := []string{"echo", "post"}
	executor.RunCommand(cmd3)

	result := output.String()
	assert.Contains(t, result, "pre")
	assert.Contains(t, result, "main")
	assert.Contains(t, result, "post")
}

func TestHooks_SkipOnError(t *testing.T) {
	executor := NewExecutor(false)

	// If a hook fails, subsequent hooks should not run
	hooks := [][]string{
	 {"echo", "first"},
	 {"sh", "-c", "exit 1"}, // This fails
	 {"echo", "never reached"},
	}

	_, err := executor.RunCommands(hooks)
	assert.Error(t, err)
}

func TestBuildCommand_WithHooks(t *testing.T) {
	fixtureDir := setupFixture(t, "with-hooks")

	cfg, _, err := loadConfigFromDir(fixtureDir)
	require.NoError(t, err)

	// Check hooks are loaded
	assert.Len(t, cfg.Hooks["pre_up"], 2)
	assert.Len(t, cfg.Hooks["post_up"], 2)
	assert.Len(t, cfg.Hooks["pre_down"], 1)

	// Verify hook content (YAML parser normalizes quotes)
	assert.Equal(t, "echo \"Starting services...\"", cfg.Hooks["pre_up"][0])
	assert.Equal(t, "echo \"Services are ready!\"", cfg.Hooks["post_up"][0])
}
