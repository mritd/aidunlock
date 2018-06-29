package unlock

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"log"

	"errors"

	"regexp"

	"fmt"

	"strings"

	"github.com/json-iterator/go"
	"github.com/spf13/viper"
)

type AppleID struct {
	ID        string            `yml:"id"`
	Birthday  string            `yml:"birthday"`
	Password  string            `yml:"password"`
	Questions map[string]string `yml:"questions"`
	Client    *http.Client      `yml:"client,omitempty"`
	Cron      string            `yml:"cron"`
}

type Question struct {
	ID       int    `json:"id"`
	Number   int    `json:"number"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Questions struct {
	Questions []Question `json:"questions"`
}

func (appID AppleID) Unlock() error {
	var sstt, location string
	var err error

	sstt, err = appID.ValidateAppleID()
	if err != nil {
		return err
	}

	sstt, location, err = appID.SelectAuthenticationMethod(sstt)
	if err != nil {
		return err
	}

	sstt, location, err = appID.ValidateBirthday(sstt, location)
	if err != nil {
		return err
	}

	sstt, location, err = appID.AnswerQuestion(sstt, location)
	if err != nil {
		return err
	}

	if strings.HasPrefix(location, "/password/reset") {
		return appID.RestPassword(sstt, location)
	} else {
		sstt, location, err = appID.UnlockAppleID(sstt, location)
		if err != nil {
			return err
		}

		return appID.ValidatePassword(sstt, location)
	}

	return nil
}

func (appID AppleID) ValidateAppleID() (string, error) {

	log.Printf("Validate AppID ==> %s\n", appID.ID)

	// create request
	req, err := http.NewRequest("POST", BaseURL+"/password/verify/appleid", bytes.NewBufferString(`{"id":"`+appID.ID+`"}`))
	if !CheckErr(err) {
		return "", err
	}

	// set Header
	setCommonHeader(req, JSON, "")

	// request
	resp, err := appID.Client.Do(req)
	if !CheckErr(err) {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Validate AppID Request Failed ")
	}

	location := resp.Header.Get("Location")

	if strings.HasPrefix(location, "/recovery/options") {
		log.Println("Apple ID []")
	}

	// get sstt and question
	req, err = http.NewRequest("GET", BaseURL+location, nil)
	if !CheckErr(err) {
		return "", err
	}

	// set Header
	setCommonHeader(req, JSON, "")

	// request
	resp, err = appID.Client.Do(req)
	if !CheckErr(err) {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if !CheckErr(err) {
		return "", err
	}

	any := jsoniter.Get(b, "sstt")
	sstt := any.ToString()
	return sstt, nil
}

func (appID AppleID) SelectAuthenticationMethod(sstt string) (string, string, error) {

	log.Println("Select Authentication Method")

	// create request
	req, err := http.NewRequest("GET", BaseURL+"/password/authenticationmethod", nil)
	if !CheckErr(err) {
		return "", "", err
	}

	// set header
	setCommonHeader(req, HTML, sstt)

	// request
	resp, err := appID.Client.Do(req)
	if !CheckErr(err) {
		return "", "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if !CheckErr(err) {
		return "", "", err
	}

	reg := regexp.MustCompile(`appConfig\.bootData\.sstt = encodeURIComponent\("(.*)"\)`)
	strs := reg.FindStringSubmatch(string(b))
	if len(strs) != 2 {
		return "", "", errors.New("Find sstt Failed ")
	}
	sstt = strs[1]

	// create request
	req, err = http.NewRequest("POST", BaseURL+"/password/authenticationmethod", bytes.NewBufferString(`{"type":"questions"}`))
	if !CheckErr(err) {
		return "", "", err
	}

	// set header
	setCommonHeader(req, JSON, sstt)

	// request
	resp, err = appID.Client.Do(req)
	if !CheckErr(err) {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("Select Authentication Method Failed ")
	}

	sstt = resp.Header.Get("sstt")
	location := resp.Header.Get("Location")
	return sstt, location, nil
}

func (appID AppleID) ValidateBirthday(sstt, location string) (string, string, error) {

	log.Println("Validate Birthday")

	// get sstt
	req, err := http.NewRequest("GET", BaseURL+location, nil)
	if !CheckErr(err) {
		return "", "", err
	}

	setCommonHeader(req, JSON, sstt)

	// request
	resp, err := appID.Client.Do(req)
	if !CheckErr(err) {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("Validate Birthday Failed ")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if !CheckErr(err) {
		return "", "", err
	}

	any := jsoniter.Get(b, "sstt")
	sstt = any.ToString()

	// get sstt
	req, err = http.NewRequest("POST", BaseURL+"/password/verify/birthday", bytes.NewBufferString(fmt.Sprintf(`{"monthOfYear":"%s","dayOfMonth":"%s","year":"%s"}`, appID.Birthday[4:6], appID.Birthday[6:], appID.Birthday[:4])))
	if !CheckErr(err) {
		return "", "", err
	}

	setCommonHeader(req, JSON, sstt)

	// request
	resp, err = appID.Client.Do(req)
	if !CheckErr(err) {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("Validate Birthday Failed ")
	}

	sstt = resp.Header.Get("sstt")
	location = resp.Header.Get("Location")
	return sstt, location, nil
}

func (appID AppleID) AnswerQuestion(sstt, location string) (string, string, error) {
	log.Println("Answer Question")

	// get sstt
	req, err := http.NewRequest("GET", BaseURL+location, nil)
	if !CheckErr(err) {
		return "", "", err
	}

	setCommonHeader(req, JSON, sstt)

	// request
	resp, err := appID.Client.Do(req)
	if !CheckErr(err) {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("Answer Question Failed ")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if !CheckErr(err) {
		return "", "", err
	}

	any := jsoniter.Get(b, "sstt")
	sstt = any.ToString()
	any = jsoniter.Get(b, "questions")

	questions := Questions{}

	for i := 0; i < any.Size(); i++ {
		log.Printf("Question [%d]: %s\n", i+1, any.Get(i, "question").ToString())
		q := Question{
			ID:       any.Get(i, "id").ToInt(),
			Number:   any.Get(i, "number").ToInt(),
			Question: any.Get(i, "question").ToString(),
			Answer:   appID.Questions[any.Get(i, "question").ToString()],
		}
		questions.Questions = append(questions.Questions, q)
	}
	json, err := jsoniter.Marshal(questions)
	if !CheckErr(err) {
		return "", "", err
	}

	log.Printf("Answer Question JSON: %s\n", string(json))

	// get sstt
	req, err = http.NewRequest("POST", BaseURL+"/password/verify/questions", bytes.NewBuffer(json))
	if !CheckErr(err) {
		return "", "", err
	}

	setCommonHeader(req, JSON, sstt)

	// request
	resp, err = appID.Client.Do(req)
	if !CheckErr(err) {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err = ioutil.ReadAll(resp.Body)
		if !CheckErr(err) {
			return "", "", err
		}
		log.Println(string(b))
		return "", "", errors.New("Answer Question Failed ")
	}

	sstt = resp.Header.Get("sstt")
	location = resp.Header.Get("Location")
	return sstt, location, nil
}

func (appID AppleID) UnlockAppleID(sstt, location string) (string, string, error) {
	log.Println("Unlock AppleID")

	// get sstt
	req, err := http.NewRequest("GET", BaseURL+location, nil)
	if !CheckErr(err) {
		return "", "", err
	}

	setCommonHeader(req, JSON, sstt)

	// request
	resp, err := appID.Client.Do(req)
	if !CheckErr(err) {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("Unlock AppleID Failed ")
	}

	// get sstt
	sstt = resp.Header.Get("sstt")

	// unlock
	req, err = http.NewRequest("POST", BaseURL+"/password/reset/options", bytes.NewBufferString(`{"type":"unlock_account"}`))
	if !CheckErr(err) {
		return "", "", err
	}

	setCommonHeader(req, JSON, "")
	req.Header.Set("sstt", sstt)

	// request
	resp, err = appID.Client.Do(req)
	if !CheckErr(err) {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("Unlock AppleID Failed ")
	}

	sstt = resp.Header.Get("sstt")
	location = resp.Header.Get("Location")
	return sstt, location, nil
}

func (appID AppleID) RestPassword(sstt, location string) error {
	log.Println("Reset AppleID Password")

	// get sstt
	req, err := http.NewRequest("GET", BaseURL+location, nil)
	if !CheckErr(err) {
		return err
	}

	setCommonHeader(req, JSON, sstt)

	// request
	resp, err := appID.Client.Do(req)
	if !CheckErr(err) {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Reset AppleID Failed ")
	}

	// get sstt
	sstt = resp.Header.Get("sstt")
	password := RandStr(16)

	log.Printf("Reset Password: %s\n", password)

	// unlock
	req, err = http.NewRequest("POST", BaseURL+"/password/reset", bytes.NewBufferString(`{"password":"`+password+`"}`))
	if !CheckErr(err) {
		return err
	}

	setCommonHeader(req, JSON, "")
	req.Header.Set("sstt", sstt)

	// request
	resp, err = appID.Client.Do(req)
	if !CheckErr(err) {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 260 {
		return errors.New("Reset Password Failed ")
	} else {
		var smtp SMTPConfig
		viper.UnmarshalKey("email", &smtp)
		smtp.Send(fmt.Sprintf("Apple ID [%s] Password Reset Success!\nNew Password: %s\n", appID.ID, password))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if !CheckErr(err) {
		return err
	}
	if jsoniter.Get(b, "unlockCompleted").ToBool() {
		log.Printf("Unlock Apple ID [%s] success\n", appID.ID)
	}
	return nil
}

func (appID AppleID) ValidatePassword(sstt, location string) error {
	log.Println("Validate Password")

	// get sstt
	req, err := http.NewRequest("GET", BaseURL+location, nil)
	if !CheckErr(err) {
		return err
	}

	setCommonHeader(req, JSON, sstt)

	// request
	resp, err := appID.Client.Do(req)
	if !CheckErr(err) {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Validate Password Failed ")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if !CheckErr(err) {
		return err
	}

	any := jsoniter.Get(b, "sstt")
	sstt = any.ToString()

	// send password
	req, err = http.NewRequest("POST", BaseURL+"/password/unlock", bytes.NewBufferString("{\"password\":\""+appID.Password+"\"}"))
	if !CheckErr(err) {
		return err
	}

	setCommonHeader(req, JSON, sstt)

	// request
	resp, err = appID.Client.Do(req)
	if !CheckErr(err) {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 260 {
		return errors.New("Validate Password Failed ")
	}

	b, err = ioutil.ReadAll(resp.Body)
	if !CheckErr(err) {
		return err
	}
	if jsoniter.Get(b, "unlockCompleted").ToBool() {
		log.Printf("Unlock Apple ID [%s] success\n", appID.ID)
	}
	return nil
}
