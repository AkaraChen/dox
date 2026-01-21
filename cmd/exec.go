package commands

import (
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec SERVICE COMMAND [ARGS...]",
	Short: "Execute command in service container",
	Long: `Execute a command in a running service container.

Requires at least a service name. Common usage: do c exec api bash`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommand(func(b *Builder, a []string) ([]string, error) {
   return b.BuildExec(a)
	 }, args)
	},
}

func init() {
	composeGroupCmd.AddCommand(execCmd)
}
