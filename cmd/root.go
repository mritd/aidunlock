package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/mritd/aidunlock/unlock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version string
var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "aidunlock",
	Short: "A simple Apple ID unlock tool.",
	Long: `
A simple Apple ID unlock tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		unlock.Boot()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLog)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is aidunlock.yaml)")
}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cfgFile = "aidunlock.yaml"
		viper.SetConfigFile(cfgFile)

		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			_, err = os.Create(cfgFile)
			unlock.CheckAndExit(err)
			viper.Set("AppleIDs", unlock.ExampleConfig())
			viper.Set("Email", unlock.SMTPExampleConfig())
			err = viper.WriteConfig()
			unlock.CheckAndExit(err)
		} else {
			unlock.CheckAndExit(err)
		}

	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	unlock.CheckAndExit(err)

}

func initLog() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
