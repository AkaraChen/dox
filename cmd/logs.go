package commands

import (
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs [service...]",
	Short: "View container logs",
	Long: `View logs from services.

Can show logs for all services or specific services. Supports -f (follow) and --tail flags.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommand(func(b *Builder, a []string) ([]string, error) {
   return b.BuildLogs(a)
	 }, args)
	},
}

func init() {
	composeGroupCmd.AddCommand(logsCmd)
}
