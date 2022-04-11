package service

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const (
	KVKeyGoogleAuthExpiration  = "google_auth_expiration"
	KVKeyGoogleAuthAccessToken = "google_auth_access_token"
)

type Google struct {
	authExpiration  time.Time
	authAccessToken string
	refreshToken    string
	clientID        string
	clientSecret    string
}

func NewGoogle(clientID, clientSecret, refreshToken string) *Google {
	return &Google{
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
	}
}

type GoogleRefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (gp *Google) RefreshAuthToken() error {
	if gp.authAccessToken == "" || gp.authExpiration.Before(time.Now()) {
		params := url.Values{
			"client_id":     {gp.clientID},
			"client_secret": {gp.clientSecret},
			"refresh_token": {gp.refreshToken},
			"grant_type":    {"refresh_token"},
		}
		url := "https://accounts.google.com/o/oauth2/token?" + params.Encode()

		httpResponse, err := http.Post(url, "text/plain", nil)
		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return err
		}

		var response GoogleRefreshTokenResponse
		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return err
		}

		gp.authExpiration = time.Now().Add(time.Duration(response.ExpiresIn * int(time.Second)))
		gp.authAccessToken = response.AccessToken
	}

	return nil
}

type Event struct {
	Name  string
	Start time.Time
	End   time.Time
}

func (gp *Google) CreateNewEvent(e Event) (*calendar.Event, error) {
	err := gp.RefreshAuthToken()
	if err != nil {
		return nil, err
	}

	if gp.authAccessToken == "" {
		return nil, errors.New("no access token available for google")
	}

	config := &oauth2.Config{}
	ctx := context.Background()
	calendarService, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, &oauth2.Token{AccessToken: gp.authAccessToken})))
	if err != nil {
		return nil, err
	}

	event := &calendar.Event{
		Summary: e.Name,
		Start: &calendar.EventDateTime{
			DateTime: e.Start.Format(time.RFC3339),
			TimeZone: e.Start.Local().String(),
		},
		End: &calendar.EventDateTime{
			DateTime: e.End.Format(time.RFC3339),
			TimeZone: e.End.Local().String(),
		},
	}

	calendarId := "primary"

	return calendarService.Events.Insert(calendarId, event).Do()
}
