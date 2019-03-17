package unlock

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/sirupsen/logrus"
)

func Uninstall() {

	if runtime.GOOS == "linux" {
		CheckRoot()

		logrus.Print("Stop AppleID Unlock")
		_ = exec.Command("systemctl", "stop", "aidunlock").Run()

		logrus.Print("Clean files")
		_ = os.Remove(binPath)
		_ = os.Remove(configDir)
		_ = os.Remove(servicePath)

		logrus.Print("Systemd reload")
		_ = exec.Command("systemctl", "daemon-reload").Run()
	} else {
		logrus.Print("Uninstall not support this platform!")
	}

}
