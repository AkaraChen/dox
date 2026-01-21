package commands

import (
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart SERVICE [SERVICE...]",
	Short: "Restart services",
	Long: `Restart specific services.

Requires at least one service name as argument.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommand(func(b *Builder, a []string) ([]string, error) {
   return b.BuildRestart(a)
	 }, args)
	},
}

func init() {
	composeGroupCmd.AddCommand(restartCmd)
}
