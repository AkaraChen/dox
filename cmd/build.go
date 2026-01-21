package commands

import (
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [SERVICE...]",
	Short: "Build service images",
	Long: `Build or rebuild service images.

Can build all services or specific services.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommand(func(b *Builder, a []string) ([]string, error) {
   return b.BuildBuild(a)
	 }, args)
	},
}

func init() {
	composeGroupCmd.AddCommand(buildCmd)
}
