package unlock

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/viper"
)

func Install() {

	Uninstall()

	if runtime.GOOS == "linux" {

		log.Printf("Create config dir %s\n", configDir)
		_ = os.MkdirAll(configDir, 0755)

		log.Printf("Copy file to %s\n", binPath)
		currentPath, err := exec.LookPath(os.Args[0])
		CheckAndExit(err)

		currentFile, err := os.Open(currentPath)
		CheckAndExit(err)
		defer func() {
			_ = currentFile.Close()
		}()

		installFile, err := os.OpenFile(binPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
		CheckAndExit(err)
		defer func() {
			_ = installFile.Close()
		}()
		_, err = io.Copy(installFile, currentFile)
		CheckAndExit(err)

		log.Printf("Create config file %s\n", configFilePath)
		configFile, err := os.Create(configFilePath)
		CheckAndExit(err)
		defer func() {
			_ = configFile.Close()
		}()
		viper.SetConfigFile(configFilePath)
		viper.Set("AppleIDs", ExampleConfig())
		CheckAndExit(viper.WriteConfig())

		log.Printf("Create systemd config file %s\n", servicePath)
		systemdServiceFile, err := os.Create(servicePath)
		CheckAndExit(err)
		defer func() {
			_ = systemdServiceFile.Close()
		}()
		_, _ = fmt.Fprint(systemdServiceFile, SystemdConfig)

	} else {
		log.Println("Install not support this platform!")
	}
}
