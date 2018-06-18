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
		os.Remove("/usr/bin/aidunlock")
		os.Remove("/etc/aidunlock")
		os.Remove("/lib/systemd/system/aidunlock.service")

		log.Println("Systemd reload")
		exec.Command("systemctl", "daemon-reload").Run()
	} else {
		log.Println("Uninstall not support this platform!")
	}

}
