package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hal9000/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const (
	KVKeyGoogleAuthExpiration  = "google_auth_expiration"
	KVKeyGoogleAuthAccessToken = "google_auth_access_token"
)

type GoogleRefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func RefreshGoogleTokenIfNeeded() error {
	expiration := util.KVStoreInstance.GetInt(KVKeyGoogleAuthExpiration, 0)
	accessToken := util.KVStoreInstance.GetString(KVKeyGoogleAuthAccessToken, "")
	fmt.Println("token info", expiration, accessToken, time.Unix(int64(expiration), 0), time.Now())
	if accessToken == "" || expiration == 0 || time.Unix(int64(expiration), 0).Before(time.Now()) {
		refreshToken := os.Getenv("GOOGLE_REFRESH_TOKEN")

		params := url.Values{
			"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
			"client_secret": {os.Getenv("GOOGLE_CLIENT_SECRET")},
			"refresh_token": {refreshToken},
			"grant_type":    {"refresh_token"},
		}
		url := "https://www.googleapis.com/oauth2/v4/token?" + params.Encode()

		httpResponse, err := http.Post(url, "text/plain", nil)
		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))

		var response GoogleRefreshTokenResponse
		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return err
		}

		util.KVStoreInstance.Set(KVKeyGoogleAuthExpiration, int(time.Now().Unix())+(response.ExpiresIn/2), time.Time{})
		util.KVStoreInstance.Set(KVKeyGoogleAuthAccessToken, response.AccessToken, time.Time{})
	}

	return nil
}

func StartGoogleTokenRefreshCycle(_ *chan util.ResponseMessage) {
	for {
		err := RefreshGoogleTokenIfNeeded()
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Minute)
	}
}

func CreateNewEvent(title string, start, end time.Time) error {
	err := RefreshGoogleTokenIfNeeded()
	if err != nil {
		return err
	}

	accessToken := util.KVStoreInstance.GetString(KVKeyGoogleAuthAccessToken, "")
	if accessToken == "" {
		return errors.New("no access token available for google")
	}

	config := &oauth2.Config{}
	ctx := context.Background()
	calendarService, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, &oauth2.Token{AccessToken: accessToken})))
	if err != nil {
		return err
	}

	event := &calendar.Event{
		Summary: title,
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: start.Local().String(),
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: start.Local().String(),
		},
	}

	calendarId := "primary"
	_, err = calendarService.Events.Insert(calendarId, event).Do()

	return err
}
