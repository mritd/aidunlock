package unlock

import (
	"io"
	"log"
	"os"
	"runtime"
)

func Reinstall() {
	if runtime.GOOS == "linux" {
		CheckRoot()

		configExist := false

		log.Println("Backup config")
		// check config
		if _, err := os.Stat(configFilePath); err == nil {
			f, err := os.Open(configFilePath)
			if err != nil {
				log.Panicln("Backup config failed")
			}
			defer func() {
				_ = f.Close()
			}()

			bak, err := os.OpenFile("/tmp/aidunlock.yaml", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
			if err != nil {
				log.Panicln("Creat backup config failed")
			}
			defer func() {
				_ = bak.Close()
			}()

			_, err = io.Copy(bak, f)
			if err != nil {
				log.Panicln("Backup config failed")
			}
			configExist = true
		}

		// Reinstall
		Install()

		if configExist {

			log.Println("Restore config")

			f, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
			if err != nil {
				log.Panicln("Restore config failed")
			}
			defer func() {
				_ = f.Close()
			}()

			bak, err := os.Open("/tmp/aidunlock.yaml")
			if err != nil {
				log.Panicln("Open backup config failed")
			}
			defer func() {
				_ = bak.Close()
			}()

			_, err = io.Copy(f, bak)
			if err != nil {
				log.Panicln("Restore config failed")
			}
		}

	} else {
		log.Println("Install not support this platform!")
	}

}
