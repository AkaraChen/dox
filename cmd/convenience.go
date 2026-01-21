package commands

import (
	"github.com/spf13/cobra"
)

// dupCmd represents the dup command (down then up)
var dupCmd = &cobra.Command{
	Use:   "dup",
	Short: "Down then up (restart services)",
	Long: `Stop services and start them again.

Equivalent to running 'down' followed by 'up'.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
	 // dup doesn't pass args through - it's a fixed operation
	 return executeCommands(func(b *Builder) ([][]string, error) {
   return b.BuildDup()
	 })
	},
}

// nukeCmd represents the nuke command
var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Remove everything (containers, volumes, orphans)",
	Long: `Completely remove containers, volumes, and orphaned containers.

Equivalent to 'docker compose down -v --remove-orphans'. Use with caution!`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommands(func(b *Builder) ([][]string, error) {
   return b.BuildNuke()
	 })
	},
}

// freshCmd represents the fresh command
var freshCmd = &cobra.Command{
	Use:   "fresh",
	Short: "Clean rebuild (down -v && up --build)",
	Long: `Stop everything, remove volumes, and start fresh with rebuilt images.

Equivalent to 'down -v' followed by 'up --build'.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
	 return executeCommands(func(b *Builder) ([][]string, error) {
   return b.BuildFresh()
	 })
	},
}

func init() {
	composeGroupCmd.AddCommand(dupCmd)
	composeGroupCmd.AddCommand(nukeCmd)
	composeGroupCmd.AddCommand(freshCmd)
}
