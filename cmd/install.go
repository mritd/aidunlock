package cmd

import (
	"github.com/mritd/aidunlock/unlock"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install AppleID Unlock",
	Long: `
Install AppleID Unlock.`,
	Run: func(cmd *cobra.Command, args []string) {
		unlock.Install()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
