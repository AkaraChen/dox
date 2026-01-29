package commands

import (
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start Docker Compose services",
	Long: `Start Docker Compose services.

Auto-discovers compose.yaml and slice files (compose.*.yaml) in the current directory.
Use -p to select a profile from dox.yaml, or -d for detached mode.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommand(func(b *Builder, a []string) ([]string, error) {
   return b.BuildUp(a)
	 }, args)
	},
}

func init() {
	composeGroupCmd.AddCommand(upCmd)
}
