package service

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var kasaMQTTConnection mqtt.Client

func InitKasaConnection() error {
	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("KASA_MQTT_URL"))
	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	kasaMQTTConnection = c

	return nil
}

func SetKasaDeviceStatus(id string, on bool) error {
	var message string
	if on {
		message = "on"
	} else {
		message = "off"
	}
	topic := fmt.Sprintf("/%s/switch", id)
	// fmt.Println(topic)
	token := kasaMQTTConnection.Publish(topic, 0, false, message)
	token.Wait()
	return token.Error()
}
