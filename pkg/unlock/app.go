package unlock

import (
	"log"
	"os"

	"net/http/cookiejar"

	"net/http"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

func Boot() {
	var appleIDs []AppleID
	err := viper.UnmarshalKey("AppleIDs", &appleIDs)
	if err != nil {
		log.Println("Can't parse server config!")
		os.Exit(1)
	}

	c := cron.New()
	for i := range appleIDs {
		x := i

		log.Printf("Apple ID [%s] cron starting\n", appleIDs[x].ID)

		c.AddFunc(appleIDs[x].Cron, func() {

			jar, _ := cookiejar.New(nil)

			appleIDs[x].Client = &http.Client{
				Timeout: 5 * time.Second,
				Jar:     jar,
			}
			if appleIDs[x].Check() {
				appleIDs[x].Unlock()
			}
		})
	}
	c.Start()
	select {}
}
