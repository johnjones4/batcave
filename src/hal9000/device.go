package hal9000

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

const (
	DeviceTypeKasa  = "kasa"
	DeviceTypeGroup = "group"
)

type Device struct {
	Names     []string `json:"names"`
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	DeviceIDs []string `json:"devices"`
}

func (d Device) Devices() []Device {
	groupDevices := make([]Device, len(d.DeviceIDs))
	for i, deviceId := range d.DeviceIDs {
		for _, device := range devices {
			if device.ID == deviceId {
				groupDevices[i] = device
				break
			}
		}
	}
	return groupDevices
}

var devices []Device

var ErrorDeviceNotFound = errors.New("device not found")

func InitDevices() error {
	bytes, err := ioutil.ReadFile(os.Getenv("DEVICE_MANIFEST_PATH"))
	if err != nil {
		return err
	}
	devices = nil
	err = json.Unmarshal(bytes, &devices)
	if err != nil {
		return err
	}
	return nil
}

func FindDeviceInString(str string) (Device, error) {
	lcStr := strings.ToLower(str)
	for _, device := range devices {
		for _, name := range device.Names {
			lcName := strings.ToLower(name)
			if strings.Contains(lcStr, lcName) {
				return device, nil
			}
		}
	}
	return Device{}, ErrorDeviceNotFound
}
