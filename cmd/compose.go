package commands

import (
	"fmt"
	"os"
	"path/filepath"

	composepkg "github.com/akrc/do/internal/compose"
	"github.com/akrc/do/internal/config"
	"github.com/spf13/cobra"
)

var (
	profile string
)

// Type aliases for use in other files
type Builder = composepkg.Builder

// composeGroupCmd represents the compose command group
var composeGroupCmd = &cobra.Command{
	Use:   "c",
	Short: "Docker Compose commands (shorthand for compose)",
	Long: `Manage Docker Compose services with auto-discovery of compose files and profiles.

The compose command group (aliased as 'c') provides shorthand access to common
docker compose operations. It automatically discovers compose.yaml and slice files
(compose.*.yaml) in the current directory.`,
}

func init() {
	rootCmd.AddCommand(composeGroupCmd)

	// Global flags for compose commands
	composeGroupCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "profile to use from do.yaml")
}

// getComposeBuilder creates a builder for the current directory
func getComposeBuilder() (*Builder, error) {
	dir, err := os.Getwd()
	if err != nil {
	 return nil, err
	}

	// Load config if exists
	cfg, _, err := config.LoadConfigFromDirectory(dir)
	if err != nil {
	 return nil, err
	}

	// Use profile from flag or config
	profileToUse := profile
	if profileToUse == "" && cfg != nil {
	 profileToUse = cfg.GetDefaultProfile()
	}

	return composepkg.NewBuilder(dir, cfg, profileToUse), nil
}

// getComposeExecutor creates an executor with current settings
func getComposeExecutor() *composepkg.Executor {
	executor := composepkg.NewExecutor(IsDryRun())
	if IsVerbose() {
	 executor.SetEnv([]string{"DOCKER_COMPOSE_VERBOSE=1"})
	}
	return executor
}

// executeCommand builds and executes a command
func executeCommand(buildFunc func(*Builder, []string) ([]string, error), args []string) error {
	builder, err := getComposeBuilder()
	if err != nil {
	 return err
	}

	cmd, err := buildFunc(builder, args)
	if err != nil {
	 return err
	}

	executor := getComposeExecutor()

	// Set working directory
	dir, _ := os.Getwd()
	executor.SetDir(dir)

	if IsDryRun() || IsVerbose() {
	 output := composepkg.FormatCommand(cmd)
	 printCommand(output)
	}

	if IsDryRun() {
	 return nil
	}

	_, err = executor.RunCommand(cmd)
	return err
}

// executeCommands builds and executes multiple commands
func executeCommands(buildFunc func(*Builder) ([][]string, error)) error {
	builder, err := getComposeBuilder()
	if err != nil {
	 return err
	}

	commands, err := buildFunc(builder)
	if err != nil {
	 return err
	}

	executor := getComposeExecutor()

	// Set working directory
	dir, _ := os.Getwd()
	executor.SetDir(dir)

	if IsDryRun() || IsVerbose() {
	 for _, cmd := range commands {
   output := composepkg.FormatCommand(cmd)
   printCommand(output)
	 }
	}

	if IsDryRun() {
	 return nil
	}

	_, err = executor.RunCommands(commands)
	return err
}

// printCommand prints a command in a formatted way
func printCommand(cmd string) {
	fmt.Println(cmd)
}

// resolveFile resolves a -f flag to an absolute path
func resolveFile(file string) (string, error) {
	if filepath.IsAbs(file) {
	 return file, nil
	}
	dir, err := os.Getwd()
	if err != nil {
	 return "", err
	}
	return filepath.Join(dir, file), nil
}
