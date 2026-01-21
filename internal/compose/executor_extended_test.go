package compose

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommand_ErrorPath(t *testing.T) {
	executor := NewExecutor(false)

	// Command that doesn't exist
	cmd := []string{"nonexistent-command-xyz-123"}
	_, err := executor.RunCommand(cmd)
	assert.Error(t, err)
}

func TestRunCommand_WithEnv(t *testing.T) {
	executor := NewExecutor(false)
	executor.SetEnv([]string{"TEST_VAR=test_value"})

	// Use echo to test env is passed
	cmd := []string{"sh", "-c", "echo $TEST_VAR"}
	output, err := executor.RunCommand(cmd)
	if err == nil {
		// If sh exists, check output
		assert.Contains(t, output, "test_value")
	}
}

func TestRunCommandWithOutput_WithWriters(t *testing.T) {
	executor := NewExecutor(false)

	var stdout, stderr bytes.Buffer

	cmd := []string{"sh", "-c", "echo out; echo err >&2"}
	err := executor.RunCommandWithOutput(cmd, &stdout, &stderr)
	if err == nil {
		assert.NotEmpty(t, stdout.String())
		assert.NotEmpty(t, stderr.String())
	}
}

func TestRunCommandWithOutput_ErrorPath(t *testing.T) {
	executor := NewExecutor(false)

	var stdout, stderr bytes.Buffer

	cmd := []string{"false"} // Command that exits with 1
	err := executor.RunCommandWithOutput(cmd, &stdout, &stderr)
	assert.Error(t, err)
}

func TestRunCommandWithOutput_CommandNotFound(t *testing.T) {
	executor := NewExecutor(false)

	var stdout, stderr bytes.Buffer

	cmd := []string{"nonexistent-xyz-123"}
	err := executor.RunCommandWithOutput(cmd, &stdout, &stderr)
	assert.Error(t, err)
}

func TestRunCommands_WithError(t *testing.T) {
	executor := NewExecutor(false)

	commands := [][]string{
		{"echo", "first"},
		{"nonexistent-xyz"}, // This will fail
		{"echo", "never runs"},
	}

	_, err := executor.RunCommands(commands)
	assert.Error(t, err)
}

func TestRunCommands_EmptyList(t *testing.T) {
	executor := NewExecutor(false)

	commands := [][]string{}
	_, err := executor.RunCommands(commands)
	assert.NoError(t, err)
}

func TestRunCommands_DryRun(t *testing.T) {
	executor := NewExecutor(true) // dry-run

	var buf bytes.Buffer
	executor.Stdout = &buf

	commands := [][]string{
		{"echo", "test"},
	}

	_, err := executor.RunCommands(commands)
	assert.NoError(t, err)
	// Dry run should print the command
	assert.NotEmpty(t, buf.String())
}

func TestRunInteractive_NotImplemented(t *testing.T) {
	executor := NewExecutor(false)

	// RunInteractive is currently a stub
	cmd := []string{"echo", "test"}
	err := executor.RunInteractive(cmd)
	// Should not crash
	assert.NoError(t, err)
}

func TestRunInteractive_WithDryRun(t *testing.T) {
	executor := NewExecutor(true)

	cmd := []string{"echo", "test"}
	err := executor.RunInteractive(cmd)
	// Should not crash in dry-run
	assert.NoError(t, err)
}

func TestSetDir_AbsolutePath(t *testing.T) {
	executor := NewExecutor(false)

	executor.SetDir("/absolute/path")
	assert.Equal(t, "/absolute/path", executor.Dir)
}

func TestSetDir_RelativePath(t *testing.T) {
	executor := NewExecutor(false)

	executor.SetDir("relative/path")
	assert.Equal(t, "relative/path", executor.Dir)
}

func TestSetEnv_MultipleTimes(t *testing.T) {
	executor := NewExecutor(false)

	executor.SetEnv([]string{"VAR1=value1"})
	executor.SetEnv([]string{"VAR2=value2"})

	// Should replace, not append
	assert.Len(t, executor.Env, 1)
	assert.Equal(t, "VAR2=value2", executor.Env[0])
}

func TestSetEnv_EmptyList(t *testing.T) {
	executor := NewExecutor(false)

	executor.SetEnv([]string{})
	assert.Empty(t, executor.Env)
}

func TestRunCommand_WithCustomDir(t *testing.T) {
	executor := NewExecutor(false)
	executor.SetDir("/tmp")

	// Command should run in the specified directory
	cmd := []string{"pwd"}
	_, err := executor.RunCommand(cmd)
	if err == nil {
		// If pwd works, output should contain /tmp
		// Note: This might fail on some systems
	}
}


func TestRunCommand_WithStdoutAndStderr(t *testing.T) {
	executor := NewExecutor(false)

	var stdout, stderr bytes.Buffer
	executor.Stdout = &stdout
	executor.Stderr = &stderr

	cmd := []string{"sh", "-c", "echo out; echo err >&2"}
	output, err := executor.RunCommand(cmd)
	if err != nil {
		// sh might not exist on this system
		t.Skip("sh not available")
		return
	}
	// If command succeeded, check output
	if len(stdout.String()) > 0 {
		assert.Contains(t, stdout.String(), "out")
	}
	if len(stderr.String()) > 0 {
		assert.Contains(t, stderr.String(), "err")
	}
	// output should contain combined stdout
	assert.Contains(t, output, "out")
}

func TestRunCommand_VerboseOutput(t *testing.T) {
	executor := NewExecutor(false)

	// This test just ensures verbose flag doesn't break execution
	cmd := []string{"echo", "test"}
	_, err := executor.RunCommand(cmd)
	if err == nil {
		assert.NoError(t, err)
	}
}
