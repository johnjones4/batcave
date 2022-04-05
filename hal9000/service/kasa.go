package service

import (
	"encoding/json"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type KasaDeviceGroup struct {
	PreferredName string   `json:"preferredName"`
	Names         []string `json:"names"`
	Devices       []string `json:"devices"`
}

type Kasa struct {
	client       mqtt.Client
	deviceGroups []KasaDeviceGroup
}

func NewKasa() (*Kasa, error) {
	devicesString, err := os.ReadFile(os.Getenv("KASA_DEVICES_FILE"))
	if err != nil {
		return nil, err
	}

	var deviceGroups []KasaDeviceGroup
	err = json.Unmarshal(devicesString, &deviceGroups)
	if err != nil {
		return nil, err
	}

	return &Kasa{
		client:       mqtt.NewClient(mqtt.NewClientOptions().AddBroker(os.Getenv("KASA_MQTT_URL"))),
		deviceGroups: deviceGroups,
	}, nil
}

func (k *Kasa) DeviceGroups() []KasaDeviceGroup {
	return k.deviceGroups
}

func (k *Kasa) SetStatus(id string, on bool) error {
	if !k.client.IsConnected() {
		if token := k.client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	var message string
	if on {
		message = "on"
	} else {
		message = "off"
	}
	topic := fmt.Sprintf("/%s/switch", id)
	token := k.client.Publish(topic, 0, false, message)
	token.Wait()
	return token.Error()
}
