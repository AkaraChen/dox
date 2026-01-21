package commands

import (
	"github.com/spf13/cobra"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running containers",
	Long: `List running containers for the current compose project.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommand(func(b *Builder, a []string) ([]string, error) {
   return b.BuildPs(a)
	 }, args)
	},
}

func init() {
	composeGroupCmd.AddCommand(psCmd)
}
