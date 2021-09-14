package service

import (
	"fmt"
	"hal9000/types"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type kasaProviderConcrete struct {
	mQTTConnection mqtt.Client
}

func InitKasaProvider() (types.KasaProvider, error) {
	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("KASA_MQTT_URL"))
	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return kasaProviderConcrete{c}, nil
}

func (kp kasaProviderConcrete) SetKasaDeviceStatus(id string, on bool) error {
	var message string
	if on {
		message = "on"
	} else {
		message = "off"
	}
	topic := fmt.Sprintf("/%s/switch", id)
	token := kp.mQTTConnection.Publish(topic, 0, false, message)
	token.Wait()
	return token.Error()
}
