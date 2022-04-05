package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	AbodeModeAway    = "away"
	AbodeModeHome    = "home"
	AbodeModeStandby = "standby"
)

type Abode struct {
	username    string
	password    string
	accessToken string
	tokenType   string
	expiration  time.Time
}

func NewAbode() *Abode {
	return &Abode{
		username: os.Getenv("ABODE_USERNAME"),
		password: os.Getenv("ABODE_PASSWORD"),
	}
}

type abodeLoginResponse struct {
	Token string `json:"token"`
}

type abodeClaimResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (a *Abode) authorize() error {
	if a.accessToken == "" || a.expiration.Before(time.Now()) {
		form := make(url.Values)
		form.Add("id", a.username)
		form.Add("password", a.password)
		form.Add("uuid", uuid.New().String())
		form.Add("locale_code", "en-US")

		loginHttpResp, err := http.PostForm("https://my.goabode.com/api/auth2/login", form)
		if err != nil {
			return err
		}

		loginRespBody, err := io.ReadAll(loginHttpResp.Body)
		if err != nil {
			return err
		}

		loginResp := abodeLoginResponse{}

		err = json.Unmarshal(loginRespBody, &loginResp)
		if err != nil {
			return err
		}

		claimReq, err := http.NewRequest("GET", "https://my.goabode.com/api/auth2/claims", nil)
		if err != nil {
			return err
		}

		claimReq.Header.Add("ABODE-API-KEY", loginResp.Token)

		claimHttpResp, err := http.DefaultClient.Do(claimReq)
		if err != nil {
			return err
		}

		claimRespBody, err := io.ReadAll(claimHttpResp.Body)
		if err != nil {
			return err
		}

		claimResp := abodeClaimResponse{}
		err = json.Unmarshal(claimRespBody, &claimResp)
		if err != nil {
			return err
		}

		a.accessToken = claimResp.AccessToken
		a.tokenType = claimResp.TokenType
		a.expiration = time.Now().UTC().Add(time.Second * time.Duration(claimResp.ExpiresIn/2))
	}

	return nil
}

func (a *Abode) request(method, path string, res interface{}) error {
	err := a.authorize()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, "https://my.goabode.com"+path, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", a.tokenType, a.accessToken))

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("abode error: %d", httpResp.StatusCode)
	}

	if res == nil {
		return nil
	}

	httpBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(httpBytes, &res)
	if err != nil {
		return err
	}

	return nil
}

func (a *Abode) SetMode(mode string) error {
	return a.request("PUT", "/api/v1/panel/mode/1/"+mode, nil)
}

type AbodeDeviceStatus struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status_display"`
}

func (a *Abode) GetDeviceStatuses() ([]AbodeDeviceStatus, error) {
	var statuses []AbodeDeviceStatus
	err := a.request("GET", "/api/v1/devices", &statuses)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

type AbodePanelMode struct {
	Label string `json:"area_1_label"`
}

type AbodePanel struct {
	Mode AbodePanelMode `json:"mode"`
}

func (a *Abode) GetPanel() (AbodePanel, error) {
	var panel AbodePanel
	err := a.request("GET", "/api/v1/panel", &panel)
	if err != nil {
		return AbodePanel{}, err
	}

	return panel, nil
}
