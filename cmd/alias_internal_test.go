package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListAliases_NoConfig tests listing aliases with no do.yaml
func TestListAliases_NoConfig(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	// Capture output
	var output strings.Builder
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = listAliases()

	w.Close()
	os.Stdout = original

	// Read output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output.Write(buf[:n])

	assert.NoError(t, err)
	assert.Contains(t, output.String(), "No aliases defined")
}

// TestListAliases_WithAliases tests listing aliases with do.yaml
func TestListAliases_WithAliases(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "with-aliases")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	// Capture output
	var output strings.Builder
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = listAliases()

	w.Close()
	os.Stdout = original

	// Read output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output.Write(buf[:n])

	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Available aliases")
}

// TestListAliases_EmptyAliases tests listing aliases with empty aliases map
func TestListAliases_EmptyAliases(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "empty-profiles")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	// Capture output
	var output strings.Builder
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = listAliases()

	w.Close()
	os.Stdout = original

	// Read output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output.Write(buf[:n])

	assert.NoError(t, err)
	assert.Contains(t, output.String(), "No aliases defined")
}

// TestExecuteAlias_Valid tests executing a valid alias
func TestExecuteAlias_Valid(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "with-aliases")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	originalDryRun := dryRun
	defer func() { dryRun = originalDryRun }()

	dryRun = true

	err = executeAlias("fresh")
	assert.NoError(t, err)
}

// TestExecuteAlias_NotFound tests executing a non-existent alias
func TestExecuteAlias_NotFound(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "with-aliases")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	err = executeAlias("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestExecuteAlias_NoConfig tests executing an alias with no do.yaml
func TestExecuteAlias_NoConfig(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	err = executeAlias("any")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no do.yaml found")
}

// TestExecuteAlias_SimpleAlias tests executing a simple alias
func TestExecuteAlias_SimpleAlias(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "with-aliases")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	// This should work fine with our test fixtures
	originalDryRun := dryRun
	defer func() { dryRun = originalDryRun }()

	dryRun = true

	err = executeAlias("fresh")
	assert.NoError(t, err)
}

// TestExecuteAlias_ChainedCommands tests alias with multiple commands
func TestExecuteAlias_ChainedCommands(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "complex-aliases")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	originalDryRun := dryRun
	defer func() { dryRun = originalDryRun }()

	dryRun = true

	err = executeAlias("chained")
	assert.NoError(t, err)
}

// TestExecuteAlias_MultiChain tests alias with many chained commands
func TestExecuteAlias_MultiChain(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "complex-aliases")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	originalDryRun := dryRun
	defer func() { dryRun = originalDryRun }()

	dryRun = true

	err = executeAlias("multi")
	assert.NoError(t, err)
}

// TestExecuteAlias_WithDryRun tests that dry-run doesn't execute
func TestExecuteAlias_WithDryRun(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "with-aliases")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	originalDryRun := dryRun
	defer func() { dryRun = originalDryRun }()

	dryRun = true

	err = executeAlias("fresh")
	assert.NoError(t, err) // Should not execute in dry-run
}

// TestExecuteAlias_VerboseOutput tests verbose output for alias execution
func TestExecuteAlias_VerboseOutput(t *testing.T) {
	fixtureDir := filepath.Join("..", "test", "fixtures", "with-aliases")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(fixtureDir)
	require.NoError(t, err)

	originalVerbose := verbose
	originalDryRun := dryRun
	defer func() {
		verbose = originalVerbose
		dryRun = originalDryRun
	}()

	verbose = true
	dryRun = true

	// Capture output to check for verbose messages
	var output strings.Builder
	original := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = executeAlias("fresh")

	w.Close()
	os.Stdout = original

	buf := make([]byte, 2048)
	n, _ := r.Read(buf)
	output.Write(buf[:n])

	assert.NoError(t, err)
	// Verbose should show which alias is being executed
	assert.Contains(t, output.String(), "Executing alias")
}
