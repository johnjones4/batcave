package hal9000

import (
	"encoding/json"
	"errors"
	"hal9000/util"
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

func (d Device) GetNames() []string {
	return d.Names
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

type NameableDeviceSequenceItem struct {
	Name     string
	Nameable Device
}

func FindDeviceInString(str string) (Device, error) {
	nameables := make([]util.Nameable, len(devices))
	for i, d := range devices {
		nameables[i] = d
	}
	sortedNameables := util.GenerateNameableSequence(nameables)
	lcStr := strings.ToLower(str)
	for _, nameable := range sortedNameables {
		lcName := strings.ToLower(nameable.Name)
		if strings.Contains(lcStr, lcName) {
			return nameable.Nameable.(Device), nil
		}
	}
	return Device{}, ErrorDeviceNotFound
}
