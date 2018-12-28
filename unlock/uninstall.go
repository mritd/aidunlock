package unlock

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

func Uninstall() {

	if runtime.GOOS == "linux" {
		CheckRoot()

		log.Println("Stop AppleID Unlock")
		_ = exec.Command("systemctl", "stop", "aidunlock").Run()

		log.Println("Clean files")
		_ = os.Remove(binPath)
		_ = os.Remove(configDir)
		_ = os.Remove(servicePath)

		log.Println("Systemd reload")
		_ = exec.Command("systemctl", "daemon-reload").Run()
	} else {
		log.Println("Uninstall not support this platform!")
	}

}
