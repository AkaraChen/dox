package commands

import (
	"github.com/spf13/cobra"
)

// sCmd is shorthand for status
var sCmd = &cobra.Command{
	Use:   "s",
	Short: "Shorthand for status",
	Long:  `Shorthand for 'do c status'. Shows enhanced status of services.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommand(func(b *Builder, a []string) ([]string, error) {
   return b.BuildStatus(a)
	 }, args)
	},
}

func init() {
	rootCmd.AddCommand(sCmd)
}
