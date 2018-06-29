package unlock

import (
	"bytes"
	"log"
	"net/http"
	"strings"
)

func (appID AppleID) Check() bool {
	log.Printf("Check AppID ==> %s\n", appID.ID)

	// create request
	req, err := http.NewRequest("POST", BaseURL+"/password/verify/appleid", bytes.NewBufferString(`{"id":"`+appID.ID+`"}`))
	if !CheckErr(err) {
		return false
	}

	// set Header
	setCommonHeader(req, JSON, "")

	// request
	resp, err := appID.Client.Do(req)
	if !CheckErr(err) {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	location := resp.Header.Get("Location")

	if strings.HasPrefix(location, "/recovery/options") {
		log.Printf("Apple ID [%s] not lock\n", appID.ID)
		return false
	}

	if strings.HasPrefix(location, "/password/authenticationmethod") {
		log.Printf("Apple ID [%s] locked\n", appID.ID)
		return true
	}

	return false
}
