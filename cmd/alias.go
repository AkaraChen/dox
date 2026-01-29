package commands

import (
	"fmt"
	"os"
	"sort"

	composepkg "github.com/AkaraChen/dox/internal/compose"
	"github.com/spf13/cobra"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias [NAME]",
	Short: "Run a custom alias defined in dox.yaml",
	Long: `Run a custom alias defined in dox.yaml.

Aliases are custom command shortcuts defined in your dox.yaml file.
They can chain multiple docker compose commands together.

With no arguments, lists all available aliases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listAliases()
		}
		aliasName := args[0]
		return executeAlias(aliasName)
	},
}

func init() {
	composeGroupCmd.AddCommand(aliasCmd)
}

// listAliases displays all available aliases
func listAliases() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	if cfg == nil || len(cfg.Aliases) == 0 {
		fmt.Println("No aliases defined in dox.yaml")
		return nil
	}

	fmt.Println("Available aliases:")
	// Sort alias names
	names := make([]string, 0, len(cfg.Aliases))
	for name := range cfg.Aliases {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("  %s: %s\n", name, cfg.Aliases[name])
	}

	return nil
}

// executeAlias executes an alias by name
func executeAlias(aliasName string) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	if cfg == nil {
		return fmt.Errorf("no dox.yaml found")
	}

	aliasDef, exists := cfg.Aliases[aliasName]
	if !exists {
		// List available aliases
		var available []string
		for name := range cfg.Aliases {
			available = append(available, name)
		}
		return fmt.Errorf("alias '%s' not found. Available aliases: %v", aliasName, available)
	}

	if IsVerbose() {
		fmt.Printf("Executing alias '%s': %s\n", aliasName, aliasDef)
	}

	commands, err := resolveAlias(aliasDef)
	if err != nil {
		return fmt.Errorf("failed to resolve alias '%s': %w", aliasName, err)
	}

	executor := getComposeExecutor()
	dir, _ := os.Getwd()
	executor.SetDir(dir)

	// Show commands in dry-run mode
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
