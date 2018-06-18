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

		log.Println("Create config dir /etc/aidunlock")
		os.MkdirAll("/etc/aidunlock", 0755)

		log.Println("Copy file to /usr/bin")
		currentPath, err := exec.LookPath(os.Args[0])
		CheckAndExit(err)

		currentFile, err := os.Open(currentPath)
		CheckAndExit(err)
		defer currentFile.Close()

		installFile, err := os.OpenFile("/usr/bin/aidunlock", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
		CheckAndExit(err)
		defer installFile.Close()
		_, err = io.Copy(installFile, currentFile)
		CheckAndExit(err)

		log.Println("Create config file /etc/aidunlock/config.yaml")
		configFile, err := os.Create("/etc/aidunlock/config.yaml")
		CheckAndExit(err)
		defer configFile.Close()
		viper.SetConfigFile("/etc/aidunlock/config.yaml")
		viper.Set("AppleIDs", ExampleConfig())
		CheckAndExit(viper.WriteConfig())

		log.Println("Create systemd config file /lib/systemd/system/aidunlock.service")
		systemdServiceFile, err := os.Create("/lib/systemd/system/aidunlock.service")
		CheckAndExit(err)
		defer systemdServiceFile.Close()
		fmt.Fprint(systemdServiceFile, SystemdConfig)

	} else {
		log.Println("Install not support this platform!")
	}
}
