package service

import (
	"fmt"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Kasa struct {
	client  mqtt.Client
	devices []string
}

func NewKasa() (*Kasa, error) {
	deviceListString, err := os.ReadFile(os.Getenv("KASA_DEVICES_FILE"))
	if err != nil {
		return nil, err
	}
	return &Kasa{
		client:  mqtt.NewClient(mqtt.NewClientOptions().AddBroker(os.Getenv("KASA_MQTT_URL"))),
		devices: strings.Split(string(deviceListString), "\n"),
	}, nil
}

func (k *Kasa) DeviceNames() []string {
	return k.devices
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
