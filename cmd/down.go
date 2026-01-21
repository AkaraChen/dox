package commands

import (
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove Docker Compose services",
	Long: `Stop and remove Docker Compose services, networks, and optionally volumes.

Supports standard docker compose flags like -v (remove volumes) and --remove-orphans.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommand(func(b *Builder, a []string) ([]string, error) {
   return b.BuildDown(a)
	 }, args)
	},
}

func init() {
	composeGroupCmd.AddCommand(downCmd)
}
