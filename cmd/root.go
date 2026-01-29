package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	verbose bool
	dryRun  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dox",
	Short: "Docker Compose wrapper for simplified multi-file management",
	Long: `dox is a CLI wrapper for Docker Compose that simplifies working with
multiple compose files (slices) and environment-specific configurations.

It auto-discovers compose.yaml and slice files (compose.*.yaml), supports
profile-based configuration, and provides shorthand commands for common
operations.`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
	 fmt.Fprintln(os.Stderr, err)
	 os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show commands without executing")
}

// GetRoot returns the root command
func GetRoot() *cobra.Command {
	return rootCmd
}

// IsVerbose returns true if verbose mode is enabled
func IsVerbose() bool {
	return verbose
}

// IsDryRun returns true if dry-run mode is enabled
func IsDryRun() bool {
	return dryRun
}
