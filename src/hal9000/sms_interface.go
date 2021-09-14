package hal9000

import (
	"crypto/sha1"
	"fmt"
	"hal9000/types"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type InterfaceTypeSMS struct {
	Number string
}

func (i InterfaceTypeSMS) Type() string {
	return "sms"
}

func (i InterfaceTypeSMS) ID() string {
	h := sha1.New()
	h.Write([]byte(i.Number))
	bs := h.Sum(nil)
	return fmt.Sprintf("sms-%x", bs)
}

func (i InterfaceTypeSMS) IsStillValid() bool {
	return true
}

func (i InterfaceTypeSMS) SupportsVisuals() bool {
	return false
}

func (i InterfaceTypeSMS) SendMessage(m types.ResponseMessage) error {
	accountSid := os.Getenv("TWILIO_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	msgData := url.Values{}
	msgData.Set("To", i.Number)
	msgData.Set("From", os.Getenv("TWILIO_NUMBER_FROM"))
	msgData.Set("Body", m.Text)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err := client.Do(req)
	return err
}
