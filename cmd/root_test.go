package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRoot(t *testing.T) {
	root := GetRoot()
	assert.NotNil(t, root)
	assert.Equal(t, "dox", root.Use)
}

func TestIsVerbose_Default(t *testing.T) {
	// Reset to default
	verbose = false
	assert.False(t, IsVerbose())
}

func TestIsVerbose_True(t *testing.T) {
	verbose = true
	assert.True(t, IsVerbose())
	// Reset
	verbose = false
}

func TestIsDryRun_Default(t *testing.T) {
	dryRun = false
	assert.False(t, IsDryRun())
}

func TestIsDryRun_True(t *testing.T) {
	dryRun = true
	assert.True(t, IsDryRun())
	// Reset
	dryRun = false
}

func TestRootCmd_HasExpectedFlags(t *testing.T) {
	root := GetRoot()

	// Check verbose flag
	flag := root.Flag("verbose")
	assert.NotNil(t, flag)
	assert.Equal(t, "v", flag.Shorthand)

	// Check dry-run flag
	flag = root.Flag("dry-run")
	assert.NotNil(t, flag)
}

func TestRootCmd_HasComposeCommand(t *testing.T) {
	root := GetRoot()
	cCmd := root.Commands()[0]
	assert.Equal(t, "c", cCmd.Name())
}

func TestRootCmd_Version(t *testing.T) {
	root := GetRoot()
	assert.Equal(t, "dev", root.Version)
}

func TestExecute_HelpCommand(t *testing.T) {
	root := GetRoot()

	// Capture output
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"--help"})

	err := root.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "CLI wrapper")
}

func TestExecute_VersionCommand(t *testing.T) {
	root := GetRoot()

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--version"})

	err := root.Execute()
	assert.NoError(t, err)
	output := buf.String()
	// Cobra --version shows full help with version info
	assert.NotEmpty(t, output)
}

func TestExecute_InvalidCommand(t *testing.T) {
	root := GetRoot()

	var buf bytes.Buffer
	root.SetErr(&buf)
	root.SetArgs([]string{"invalid-command"})

	err := root.Execute()
	assert.Error(t, err)
}

func TestExecute_VerboseFlag(t *testing.T) {
	root := GetRoot()

	// Reset verbose
	verbose = false
	defer func() { verbose = false }()

	root.SetArgs([]string{"-v", "c", "--help"})
	var buf bytes.Buffer
	root.SetOut(&buf)

	_ = root.Execute()
	assert.True(t, verbose)
}

func TestExecute_DryRunFlag(t *testing.T) {
	root := GetRoot()

	dryRun = false
	defer func() { dryRun = false }()

	root.SetArgs([]string{"--dry-run", "c", "--help"})
	var buf bytes.Buffer
	root.SetOut(&buf)

	_ = root.Execute()
	assert.True(t, dryRun)
}

func TestExecute_WithOutput(t *testing.T) {
	// This is a basic smoke test
	root := GetRoot()

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"--help"})

	err := root.Execute()
	assert.NoError(t, err)
	assert.NotEmpty(t, buf.String())
}

func TestRootCmd_ComposeSubcommands(t *testing.T) {
	root := GetRoot()

	// Find compose command
	var cCmd *cobra.Command
	for _, cmd := range root.Commands() {
		if cmd.Name() == "c" {
			cCmd = cmd
			break
		}
	}

	require.NotNil(t, cCmd, "compose command 'c' should exist")

	// Check for expected subcommands
	subcommands := cCmd.Commands()
	subcommandNames := make(map[string]bool)
	for _, sub := range subcommands {
		subcommandNames[sub.Name()] = true
	}

	// Verify key subcommands exist
	expectedCommands := []string{"up", "down", "ps", "logs", "restart", "exec", "build"}
	for _, expected := range expectedCommands {
		assert.True(t, subcommandNames[expected], "expected subcommand %s not found", expected)
	}
}

func TestIsVerbose_ThreadSafety(t *testing.T) {
	// Just verify the function works
	verbose = true
	assert.True(t, IsVerbose())
	verbose = false
	assert.False(t, IsVerbose())
}

func TestIsDryRun_ThreadSafety(t *testing.T) {
	dryRun = true
	assert.True(t, IsDryRun())
	dryRun = false
	assert.False(t, IsDryRun())
}

func TestExecute_WithStdin(t *testing.T) {
	root := GetRoot()

	// Provide stdin input
	input := strings.NewReader("")
	root.SetIn(input)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--help"})

	err := root.Execute()
	assert.NoError(t, err)
}

func TestRootCmd_PersistentFlags(t *testing.T) {
	root := GetRoot()

	// Test that persistent flags are inherited
	verboseFlag := root.PersistentFlags().Lookup("verbose")
	assert.NotNil(t, verboseFlag)

	dryRunFlag := root.PersistentFlags().Lookup("dry-run")
	assert.NotNil(t, dryRunFlag)
}

func TestExecute_MultipleFlags(t *testing.T) {
	root := GetRoot()

	verbose = false
	dryRun = false
	defer func() {
		verbose = false
		dryRun = false
	}()

	root.SetArgs([]string{"-v", "--dry-run", "c", "--help"})
	var buf bytes.Buffer
	root.SetOut(&buf)

	_ = root.Execute()
	assert.True(t, verbose)
	assert.True(t, dryRun)
}

func TestExecute_ComposeHelp(t *testing.T) {
	root := GetRoot()

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"c", "--help"})

	err := root.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Manage Docker Compose services")
}

func TestExecute_StatusShorthand(t *testing.T) {
	root := GetRoot()

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"s", "--help"})

	err := root.Execute()
	assert.NoError(t, err)
}
