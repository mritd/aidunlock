package unlock

import (
	"fmt"
	"net/http/cookiejar"
	"os"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"

	"net/http"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var idCache *cache.Cache

func Boot() {
	var appleIDs []AppleID
	err := viper.UnmarshalKey("AppleIDs", &appleIDs)
	if err != nil {
		logrus.Print("Can't parse server config!")
		os.Exit(1)
	}

	c := cron.New()
	for i := range appleIDs {
		x := i

		logrus.Printf("Apple ID [%s] cron starting\n", appleIDs[x].ID)

		_ = c.AddFunc(appleIDs[x].Cron, func() {

			jar, _ := cookiejar.New(nil)

			appleIDs[x].Client = &http.Client{
				Timeout: 5 * time.Second,
				Jar:     jar,
			}
			if appleIDs[x].Check() {

				// unlock failed count
				var idFailedCount int
				failedCount, ok := idCache.Get(appleIDs[x].ID)
				if ok {
					idFailedCount = failedCount.(int)
				}
				// no more trying to unlock more than 2 times
				if idFailedCount > 1 {
					logrus.Warnf("Apple ID [%s] failed to unlock twice", appleIDs[x].ID)
					return
				}

				err := appleIDs[x].Unlock()
				if err != nil {
					idCache.Set(appleIDs[x].ID, idFailedCount+1, cache.NoExpiration)
					var smtp SMTPConfig
					_ = viper.UnmarshalKey("email", &smtp)
					smtp.Send(fmt.Sprintf("Apple ID [%s] unlock failed: %s\n", appleIDs[x].ID, err.Error()))
				}
			}
		})
	}
	c.Start()
	select {}
}

func init() {
	idCache = cache.New(cache.NoExpiration, cache.NoExpiration)
}
