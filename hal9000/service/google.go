package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hal9000/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	KVKeyGoogleAuthExpiration  = "google_auth_expiration"
	KVKeyGoogleAuthAccessToken = "google_auth_access_token"
)

type GoogleRefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type GoogleDeviceCommandRequest struct {
	Command string      `json:"command"`
	Params  interface{} `json:"params"`
}

func RefreshGoogleTokenIfNeeded() error {
	expiration := util.KVStoreInstance.GetInt(KVKeyGoogleAuthExpiration, 0)
	accessToken := util.KVStoreInstance.GetString(KVKeyGoogleAuthAccessToken, "")
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

		var response GoogleRefreshTokenResponse
		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return err
		}

		util.SetKVValue(KVKeyGoogleAuthExpiration, int(time.Now().Unix())+(response.ExpiresIn/2), time.Time{})
		util.SetKVValue(KVKeyGoogleAuthAccessToken, response.AccessToken, time.Time{})
	}

	return nil
}

type GoogleStreamURLResponseStreamUrls struct {
	RTSPUrl string `json:"rtspUrl"`
}

type GoogleStreamURLResponseResults struct {
	StreamUrls           GoogleStreamURLResponseStreamUrls `json:"streamUrls"`
	StreamToken          string                            `json:"streamToken"`
	StreamExtensionToken string                            `json:"streamExtensionToken"`
	ExpiresAt            time.Time                         `json:"expiresAt"`
}

type GoogleStreamURLResponse struct {
	Results GoogleStreamURLResponseResults `json:"results"`
}

type GoogleStreamRefreshRequest struct {
	StreamExtensionToken string    `json:"streamExtensionToken"`
	ExpiresAt            time.Time `json:"expiresAt"`
}

func GetGoogleVideoStreamURL(deviceId string) (string, GoogleStreamRefreshRequest, error) {
	err := RefreshGoogleTokenIfNeeded()
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	accessToken := util.GetKVValueString(KVKeyGoogleAuthAccessToken, "")
	if accessToken == "" {
		return "", GoogleStreamRefreshRequest{}, errors.New("no access token available for google")
	}

	url := fmt.Sprintf("https://smartdevicemanagement.googleapis.com/v1/enterprises/%s/devices/%s:executeCommand", os.Getenv("GOOGLE_ACCOUNT_ID"), deviceId)

	reqObject := GoogleDeviceCommandRequest{"sdm.devices.commands.CameraLiveStream.GenerateRtspStream", map[string]string{}}
	body, err := json.Marshal(reqObject)
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	req, err := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewBuffer(body)))
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	var response GoogleStreamURLResponse
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	return response.Results.StreamUrls.RTSPUrl, GoogleStreamRefreshRequest{response.Results.StreamExtensionToken, response.Results.ExpiresAt}, nil
}

func RefreshGoogleVideoStreamURL(oldUrl, deviceId, extensionToken string) (string, GoogleStreamRefreshRequest, error) {
	err := RefreshGoogleTokenIfNeeded()
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	accessToken := util.GetKVValueString(KVKeyGoogleAuthAccessToken, "")
	if accessToken == "" {
		return "", GoogleStreamRefreshRequest{}, errors.New("no access token available for google")
	}

	requrl := fmt.Sprintf("https://smartdevicemanagement.googleapis.com/v1/enterprises/%s/devices/%s:executeCommand", os.Getenv("GOOGLE_ACCOUNT_ID"), deviceId)

	reqObject := GoogleDeviceCommandRequest{"sdm.devices.commands.CameraLiveStream.ExtendRtspStream", map[string]string{
		"streamExtensionToken": extensionToken,
	}}
	body, err := json.Marshal(reqObject)
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	req, err := http.NewRequest("POST", requrl, ioutil.NopCloser(bytes.NewBuffer(body)))
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	var response GoogleStreamURLResponse
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	parsedUrl, err := url.Parse(oldUrl)
	if err != nil {
		return "", GoogleStreamRefreshRequest{}, err
	}

	q := parsedUrl.Query()
	q.Del("auth")
	q.Add("auth", response.Results.StreamToken)
	parsedUrl.RawQuery = q.Encode()

	return parsedUrl.String(), GoogleStreamRefreshRequest{response.Results.StreamExtensionToken, response.Results.ExpiresAt}, nil
}

//TODO subscribe to stream
