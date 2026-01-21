package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// getConfig returns the config for current directory
func getConfig() (*config.Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfg, _, err := config.LoadConfigFromDirectory(dir)
	return cfg, err
}

// getComposeExecutor creates an executor with current settings
func getComposeExecutor() *composepkg.Executor {
	executor := composepkg.NewExecutor(IsDryRun())
	if IsVerbose() {
		executor.SetEnv([]string{"DOCKER_COMPOSE_VERBOSE=1"})
	}
	return executor
}

// executeHooks executes hooks for a given hook type
func executeHooks(hookType string) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	if cfg == nil || cfg.Hooks == nil {
		return nil
	}

	hooks, exists := cfg.Hooks[hookType]
	if !exists || len(hooks) == 0 {
		return nil
	}

	executor := getComposeExecutor()
	dir, _ := os.Getwd()
	executor.SetDir(dir)

	if IsVerbose() {
		fmt.Printf("Executing %s hooks...\n", hookType)
	}

	for _, hook := range hooks {
		if IsDryRun() || IsVerbose() {
			fmt.Printf("  hook: %s\n", hook)
		}

		if IsDryRun() {
			continue
		}

		// Parse hook command
		cmd := parseHookCommand(hook)
		if _, err := executor.RunCommand(cmd); err != nil {
			return fmt.Errorf("hook failed: %s\nError: %w", hook, err)
		}
	}

	return nil
}

// parseHookCommand parses a hook string into command arguments
func parseHookCommand(hook string) []string {
	// Simple parsing - split by spaces
	// TODO: Handle quoted arguments properly
	return strings.Fields(hook)
}

// executeCommand builds and executes a command with hooks
func executeCommand(buildFunc func(*Builder, []string) ([]string, error), args []string) error {
	// Execute pre hooks if this is an 'up' command
	// Determine command type from the build function
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

	// Check if this is an 'up' command and run pre_up hooks
	if len(cmd) > 0 && cmd[len(cmd)-1] == "up" {
		if err := executeHooks("pre_up"); err != nil {
			return err
		}
	}

	// Check if this is a 'down' command and run pre_down hooks
	if len(cmd) > 0 && cmd[len(cmd)-1] == "down" {
		if err := executeHooks("pre_down"); err != nil {
			return err
		}
	}

	if IsDryRun() || IsVerbose() {
		output := composepkg.FormatCommand(cmd)
		printCommand(output)
	}

	if IsDryRun() {
		// Show post hooks in dry-run mode
		if len(cmd) > 0 && cmd[len(cmd)-1] == "up" {
			executeHooks("post_up")
		}
		if len(cmd) > 0 && cmd[len(cmd)-1] == "down" {
			executeHooks("post_down")
		}
		return nil
	}

	_, err = executor.RunCommand(cmd)
	if err != nil {
		return err
	}

	// Execute post hooks after successful command
	if len(cmd) > 0 && cmd[len(cmd)-1] == "up" {
		if err := executeHooks("post_up"); err != nil {
			return err
		}
	}

	if len(cmd) > 0 && cmd[len(cmd)-1] == "down" {
		if err := executeHooks("post_down"); err != nil {
			return err
		}
	}

	return nil
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

// resolveAlias resolves an alias definition into commands
func resolveAlias(aliasDef string) ([][]string, error) {
	if aliasDef == "" {
		return nil, fmt.Errorf("empty alias definition")
	}

	// Split by && to find chained commands
	parts := strings.Split(aliasDef, "&&")
	commands := make([][]string, 0, len(parts))

	builder, err := getComposeBuilder()
	if err != nil {
		return nil, err
	}

	for _, part := range parts {
		cmd := strings.Fields(strings.TrimSpace(part))
		if len(cmd) == 0 {
			continue
		}

		// Check if this is a known command
		if isKnownCommand(cmd[0]) {
			// Build the appropriate docker compose command
			var fullCmd []string
			switch cmd[0] {
			case "up":
				fullCmd, _ = builder.BuildUp(cmd[1:])
			case "down":
				fullCmd, _ = builder.BuildDown(cmd[1:])
			case "ps":
				fullCmd, _ = builder.BuildPs(cmd[1:])
			case "logs":
				fullCmd, _ = builder.BuildLogs(cmd[1:])
			case "restart":
				fullCmd, _ = builder.BuildRestart(cmd[1:])
			case "exec":
				fullCmd, _ = builder.BuildExec(cmd[1:])
			case "build":
				fullCmd, _ = builder.BuildBuild(cmd[1:])
			default:
				// Unknown command, pass through as-is
				fullCmd = []string{"docker", "compose"}
				fullCmd = append(fullCmd, cmd...)
			}
			commands = append(commands, fullCmd)
		} else {
			// Pass through command (could be a shell command)
			commands = append(commands, cmd)
		}
	}

	return commands, nil
}

// isKnownCommand checks if a command word is a known docker compose subcommand
func isKnownCommand(cmd string) bool {
	knownCommands := []string{
		"up", "down", "ps", "logs", "restart", "exec", "build",
		"pull", "push", "start", "stop", "rm", "kill", "run",
		"pause", "unpause", "top", "events", "port", "config",
		"create", "version",
	}
	for _, known := range knownCommands {
		if cmd == known {
			return true
		}
	}
	return false
}
