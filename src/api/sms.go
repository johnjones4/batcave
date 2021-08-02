package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"hal9000"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
)

func generateTwilioSignature(authToken string, URL string, postForm url.Values) string {
	keys := make([]string, 0, len(postForm))
	for key := range postForm {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	str := URL
	for _, key := range keys {
		str += key + postForm[key][0]
	}
	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(str))
	expectedMac := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(expectedMac)
}

func validateTwilioRequest(authToken string, URL string, request *http.Request, formValues url.Values) error {
	expectedTwilioSignature := generateTwilioSignature(authToken, URL, formValues)
	if request.Header.Get("X-Twilio-Signature") != expectedTwilioSignature {
		return errors.New("Bad X-Twilio-Signature")
	}
	return nil
}

func handleSMS(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}
	formValues, err := url.ParseQuery(string(body))
	if err != nil {
		errorResponse(w, err)
		return
	}
	err = validateTwilioRequest(os.Getenv("TWILIO_AUTH_TOKEN"), os.Getenv("SMS_ENDPOINT_URL"), req, formValues)
	if err != nil {
		errorResponse(w, err)
		return
	}

	iface := hal9000.InterfaceTypeSMS{Number: formValues.Get("From")}
	ses, err := hal9000.GetSessionWithInterfaceID(iface.ID())
	if err == hal9000.ErrorSessionNotFound {
		owner, err := hal9000.DetermineOwnerOfInterface(iface)
		if err != nil {
			errorResponse(w, err)
			return
		}
		ses = hal9000.NewSession(owner, iface)
	} else if err != nil {
		errorResponse(w, err)
		return
	}

	response, err := ses.ProcessIncomingMessage(hal9000.RequestMessage{Message: formValues.Get("Body")})
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = ses.Interface.SendMessage(response)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(""))
}
