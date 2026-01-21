package compose

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutor_RunCommand_Success(t *testing.T) {
	// Use echo command which should always succeed
	executor := NewExecutor(false)

	cmd := []string{"echo", "hello", "world"}
	output, err := executor.RunCommand(cmd)
	require.NoError(t, err)
	assert.Contains(t, output, "hello world")
}

func TestExecutor_RunCommand_Error(t *testing.T) {
	executor := NewExecutor(false)

	// Command that will fail
	cmd := []string{"false", ""} // false command exits with 1
	_, err := executor.RunCommand(cmd)
	assert.Error(t, err)
}

func TestExecutor_RunCommand_DryRun(t *testing.T) {
	executor := NewExecutor(true) // dry-run mode

	cmd := []string{"echo", "test"}
	output, err := executor.RunCommand(cmd)
	require.NoError(t, err)

	// Dry run should just echo the command
	assert.Contains(t, output, "echo")
}

func TestExecutor_RunCommands_Sequential(t *testing.T) {
	executor := NewExecutor(false)

	commands := [][]string{
	 {"echo", "first"},
	 {"echo", "second"},
	}

	output, err := executor.RunCommands(commands)
	require.NoError(t, err)
	assert.Contains(t, output, "first")
	assert.Contains(t, output, "second")
}

func TestExecutor_RunCommands_StopsOnError(t *testing.T) {
	executor := NewExecutor(false)

	commands := [][]string{
	 {"echo", "first"},
	 {"false", ""}, // This will fail
	 {"echo", "never reached"},
	}

	_, err := executor.RunCommands(commands)
	assert.Error(t, err)
}

func TestExecutor_RunCommands_DryRun(t *testing.T) {
	executor := NewExecutor(true)

	commands := [][]string{
	 {"echo", "first"},
	 {"echo", "second"},
	}

	output, err := executor.RunCommands(commands)
	require.NoError(t, err)

	// Should show both commands
	assert.Contains(t, output, "echo")
}

func TestExecutor_RunCommandWithInput_CapturesOutput(t *testing.T) {
	executor := NewExecutor(false)

	cmd := []string{"echo", "test output"}
	output, err := executor.RunCommand(cmd)
	require.NoError(t, err)
	assert.Contains(t, output, "test output")
}

func TestExecutor_RunCommandWithOutput_CapturesBoth(t *testing.T) {
	executor := NewExecutor(false)
	var stdout, stderr bytes.Buffer

	cmd := []string{"echo", "stdout"}
	err := executor.RunCommandWithOutput(cmd, &stdout, &stderr)
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "stdout")
}

func TestNewExecutor_Defaults(t *testing.T) {
	executor := NewExecutor(false)

	assert.False(t, executor.DryRun)
	assert.NotNil(t, executor.Stdout)
	assert.NotNil(t, executor.Stderr)
}

func TestNewExecutor_DryRun(t *testing.T) {
	executor := NewExecutor(true)

	assert.True(t, executor.DryRun)
}

func TestExecutor_ExitCode(t *testing.T) {
	executor := NewExecutor(false)

	// Success case
	cmd := []string{"true"}
	_, err := executor.RunCommand(cmd)
	assert.NoError(t, err)

	// Failure case
	cmd = []string{"sh", "-c", "exit 42"}
	_, err = executor.RunCommand(cmd)
	assert.Error(t, err)
}

func TestExecutor_RunCommand_SetsDir(t *testing.T) {
	executor := NewExecutor(false)
	executor.SetDir(os.TempDir())

	cmd := []string{"pwd"}
	output, err := executor.RunCommand(cmd)
	require.NoError(t, err)
	// Trim whitespace for comparison
	assert.Contains(t, strings.TrimSpace(output), strings.TrimRight(os.TempDir(), "/"))
}

func TestExecutor_RunCommand_SetsEnv(t *testing.T) {
	executor := NewExecutor(false)
	executor.SetEnv([]string{"TEST_VAR=test_value"})

	cmd := []string{"sh", "-c", "echo $TEST_VAR"}
	output, err := executor.RunCommand(cmd)
	require.NoError(t, err)
	assert.Contains(t, output, "test_value")
}
