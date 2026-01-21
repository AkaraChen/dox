package commands

import (
	"github.com/spf13/cobra"
)

// statusCmd represents the status command (enhanced ps)
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show enhanced status of services",
	Long: `Show enhanced status of running services.

This is an enhanced version of 'ps' with formatted output showing
service names, states, and port mappings in a table format.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommand(func(b *Builder, a []string) ([]string, error) {
   return b.BuildStatus(a)
	 }, args)
	},
}

func init() {
	composeGroupCmd.AddCommand(statusCmd)
}
