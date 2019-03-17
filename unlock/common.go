package unlock

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strings"
)

const (
	UA             = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	Host           = "iforgot.apple.com"
	BaseURL        = "https://iforgot.apple.com"
	Referer        = BaseURL + "/password/verify/appleid"
	AcceptEncoding = ""
	AcceptLanguage = "zh-CN,zh;q=0.9,en;q=0.8"
	AcceptJSON     = "application/json, text/javascript, */*; q=0.01"
	AcceptHTML     = "text/html;format=fragmented"
	ContentType    = "application/json"
)

const (
	JSON = "json"
	HTML = "HTML"
)

const (
	binPath        = "/usr/bin/aidunlock"
	configDir      = "/etc/aidunlock"
	configFilePath = configDir + "/config.yaml"
	servicePath    = "/lib/systemd/system/aidunlock.service"
)

const SystemdConfig = `[Unit]
Description=AppleID Unlock
Documentation=https://github.com/mritd/aidunlock
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
ExecStart=/usr/bin/aidunlock --config /etc/aidunlock/config.yaml
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target`

var letters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func setCommonHeader(req *http.Request, accept, sstt string) {

	req.Header.Set("User-Agent", UA)
	req.Header.Set("Accept-Encoding", AcceptEncoding)
	req.Header.Set("Accept-Language", AcceptLanguage)
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("Host", Host)
	req.Header.Set("Origin", BaseURL)
	req.Header.Set("Referer", Referer)

	switch accept {
	case JSON:
		req.Header.Set("Accept", AcceptJSON)
	case HTML:
		req.Header.Set("Accept", AcceptHTML)
	default:
		req.Header.Set("Accept", AcceptJSON)
	}

	if strings.TrimSpace(sstt) != "" {
		req.Header.Set("sstt", url.QueryEscape(sstt))
	}
}

func CheckErr(err error) bool {
	if err != nil {
		logrus.Print(err)
		return false
	}
	return true
}

func CheckAndExit(err error) {
	if !CheckErr(err) {
		os.Exit(1)
	}
}

func ExampleConfig() []*AppleID {
	return []*AppleID{
		{
			ID:       "apple1@apple.com",
			Birthday: "19990101",
			Password: "password",
			Questions: map[string]string{
				"Questions1？": "Answer1",
				"Questions2？": "Answer2",
				"Questions3？": "Answer3",
			},
		},
		{
			ID:       "apple2@apple.com",
			Birthday: "19990101",
			Password: "password",
			Questions: map[string]string{
				"Questions1？": "Answer1",
				"Questions2？": "Answer2",
				"Questions3？": "Answer3",
			},
		},
	}
}

func CheckRoot() {
	u, err := user.Current()
	CheckAndExit(err)

	if u.Uid != "0" || u.Gid != "0" {
		logrus.Print("This command must be run as root! (sudo)")
		os.Exit(1)
	}
}

func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
