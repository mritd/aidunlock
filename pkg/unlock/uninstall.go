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
		exec.Command("systemctl", "stop", "aidunlock").Run()

		log.Println("Clean files")
		os.Remove(binPath)
		os.Remove(configDir)
		os.Remove(servicePath)

		log.Println("Systemd reload")
		exec.Command("systemctl", "daemon-reload").Run()
	} else {
		log.Println("Uninstall not support this platform!")
	}

}
