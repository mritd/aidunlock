package cmd

import (
	"github.com/mritd/aidunlock/unlock"
	"github.com/spf13/cobra"
)

var reinstallCmd = &cobra.Command{
	Use:   "reinstall",
	Short: "Reinstall Apple ID unlock tool",
	Long: `
Reinstall Apple ID unlock tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		unlock.Reinstall()
	},
}

func init() {
	rootCmd.AddCommand(reinstallCmd)
}
