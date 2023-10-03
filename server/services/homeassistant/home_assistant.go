package homeassistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/core"
	"net/http"

	"github.com/sirupsen/logrus"
)

type HomeAssistantConfiguration struct {
	URLRoot     string                    `json:"urlBase"`
	BearerToken string                    `json:"bearerToken"`
	Groups      []core.HomeAssistantGroup `json:"groups"`
}

type HomeAssistant struct {
	Configuration HomeAssistantConfiguration
	Log           logrus.FieldLogger
}

type deviceControlRequest struct {
	EntityId string `json:"entity_id"`
}

func (ha *HomeAssistant) Groups() []core.HomeAssistantGroup {
	return ha.Groups()
}

func (ha *HomeAssistant) ToggleDeviceState(deviceId string, on bool) error {
	action := "off"
	if on {
		action = "on"
	}
	req := deviceControlRequest{
		EntityId: deviceId,
	}
	path := fmt.Sprintf("/api/services/switch/turn_%s", action)
	return ha.makeRequest("POST", path, req, nil)
}

func (ha *HomeAssistant) makeRequest(method, path string, requestObject any, responseTarget any) error {
	var reqBody io.Reader

	ha.Log.Debugf("Sending %s %s to home assistant", method, path)

	if requestObject != nil {
		reqBodyBytes, err := json.Marshal(requestObject)
		if err != nil {
			return err
		}
		reqBody = bytes.NewBuffer(reqBodyBytes)
	}

	req, err := http.NewRequest(method, ha.Configuration.URLRoot+path, reqBody)
	if err != nil {
		return err
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ha.Configuration.BearerToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%d: %s", res.StatusCode, http.StatusText(res.StatusCode))
	}

	if responseTarget == nil {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseTarget)
	if err != nil {
		return err
	}

	return nil
}
