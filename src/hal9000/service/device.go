package service

import (
	"encoding/json"
	"hal9000/types"
	"hal9000/util"
	"io/ioutil"
	"os"
	"strings"
)

type DeviceConcrete struct {
	Names     []string `json:"names"`
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	DeviceIDs []string `json:"devices"`
}

func (d DeviceConcrete) GetNames() []string {
	return d.Names
}

func (d DeviceConcrete) GetID() string {
	return d.ID
}

func (d DeviceConcrete) GetType() string {
	return d.Type
}

func (d DeviceConcrete) GetDevices(runtime types.Runtime) []types.Device {
	groupDevices := make([]types.Device, len(d.DeviceIDs))
	for i, deviceId := range d.DeviceIDs {
		for _, device := range runtime.Devices().Devices() {
			if device.GetID() == deviceId {
				groupDevices[i] = device
				break
			}
		}
	}
	return groupDevices
}

type deviceProviderConcrete struct {
	devices []types.Device
}

func InitDeviceProvider() (types.DeviceProvider, error) {
	bytes, err := ioutil.ReadFile(os.Getenv("DEVICE_MANIFEST_PATH"))
	if err != nil {
		return nil, err
	}
	var devicesConcrete []DeviceConcrete
	err = json.Unmarshal(bytes, &devicesConcrete)
	if err != nil {
		return nil, err
	}
	devices := make([]types.Device, len(devicesConcrete))
	for i, d := range devicesConcrete {
		devices[i] = d
	}
	return deviceProviderConcrete{devices}, nil
}

func (dp deviceProviderConcrete) Devices() []types.Device {
	return dp.devices
}

// type nameableDeviceSequenceItem struct {
// 	name     string
// 	nameable types.Device
// }

func (dp deviceProviderConcrete) FindDeviceInString(str string) (types.Device, error) {
	nameables := make([]types.Nameable, len(dp.devices))
	for i, d := range dp.devices {
		nameables[i] = d
	}
	sortedNameables := util.GenerateNameableSequence(nameables)
	lcStr := strings.ToLower(str)
	for _, nameable := range sortedNameables {
		lcName := strings.ToLower(nameable.Name)
		if strings.Contains(lcStr, lcName) {
			return nameable.Nameable.(types.Device), nil
		}
	}
	return nil, util.ErrorDeviceNotFound
}
