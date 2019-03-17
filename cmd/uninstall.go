package cmd

import (
	"github.com/mritd/aidunlock/unlock"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall AppleID Unlock",
	Long: `
Uninstall AppleID Unlock.`,
	Run: func(cmd *cobra.Command, args []string) {
		unlock.Uninstall()
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
