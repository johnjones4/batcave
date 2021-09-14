package intents

import (
	"errors"
	"hal9000/service"
	"hal9000/types"
	"hal9000/util"
)

type controlIntent struct {
	device types.Device
	on     bool
}

func NewControlIntent(runtime types.Runtime, m types.ParsedRequestMessage, on bool) (controlIntent, error) {
	device, err := runtime.Devices().FindDeviceInString(m.Original.Message)
	if err != nil {
		return controlIntent{}, err
	}

	return controlIntent{device, on}, nil
}

func (i controlIntent) Execute(runtime types.Runtime, lastState types.State) (types.State, types.ResponseMessage, error) {
	var err error
	if i.device.GetType() == service.DeviceTypeGroup {
		for _, device := range i.device.GetDevices(runtime) {
			err = ExecuteOnDevice(runtime, device, i.on)
			if err != nil {
				break
			}
		}
	} else {
		err = ExecuteOnDevice(runtime, i.device, i.on)
	}
	if err != nil {
		return nil, types.ResponseMessage{}, err
	}

	return lastState, util.MessageOk(), nil
}

func ExecuteOnDevice(runtime types.Runtime, device types.Device, on bool) error {
	if device.GetType() == service.DeviceTypeKasa {
		return runtime.Kasa().SetKasaDeviceStatus(device.GetID(), on)
	} else {
		return errors.New("unable to handle device type")
	}
}
