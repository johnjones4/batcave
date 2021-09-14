package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hal9000/types"
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

type googleProviderConcrete struct {
}

func InitGoogleProvider(runtime types.Runtime) (types.GoogleProvider, error) {
	gp := googleProviderConcrete{}
	go (func() {
		for {
			err := gp.RefreshAuthToken(runtime)
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Minute)
		}
	})()
	return gp, nil
}

type GoogleRefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (gp googleProviderConcrete) RefreshAuthToken(runtime types.Runtime) error {
	expiration := runtime.KVStore().GetInt(KVKeyGoogleAuthExpiration, 0)
	accessToken := runtime.KVStore().GetString(KVKeyGoogleAuthAccessToken, "")
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

		runtime.KVStore().Set(KVKeyGoogleAuthExpiration, int(time.Now().Unix())+(response.ExpiresIn/2), time.Time{})
		runtime.KVStore().Set(KVKeyGoogleAuthAccessToken, response.AccessToken, time.Time{})
	}

	return nil
}

func (gp googleProviderConcrete) CreateNewEvent(runtime types.Runtime, e types.Event) error {
	err := gp.RefreshAuthToken(runtime)
	if err != nil {
		return err
	}

	accessToken := runtime.KVStore().GetString(KVKeyGoogleAuthAccessToken, "")
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
		Summary: e.GetName(),
		Start: &calendar.EventDateTime{
			DateTime: e.GetStartTime().Format(time.RFC3339),
			TimeZone: e.GetStartTime().Local().String(),
		},
		End: &calendar.EventDateTime{
			DateTime: e.GetEndTime().Format(time.RFC3339),
			TimeZone: e.GetEndTime().Local().String(),
		},
	}

	calendarId := "primary"
	_, err = calendarService.Events.Insert(calendarId, event).Do()

	return err
}
