package hal9000

import (
	"errors"
	"hal9000/service"
	"hal9000/util"
)

type ControlIntent struct {
	Device Device
	On     bool
}

func NewControlIntent(m ParsedRequestMessage, on bool) (ControlIntent, error) {
	device, err := FindDeviceInString(m.Original.Message)
	if err != nil {
		return ControlIntent{}, err
	}

	return ControlIntent{Device: device, On: on}, nil
}

func (i ControlIntent) Execute(lastState State) (State, util.ResponseMessage, error) {
	var err error
	if i.Device.Type == DeviceTypeGroup {
		for _, device := range i.Device.Devices() {
			err = ExecuteOnDevice(device, i.On)
			if err != nil {
				break
			}
		}
	} else {
		err = ExecuteOnDevice(i.Device, i.On)
	}
	if err != nil {
		return nil, util.ResponseMessage{}, err
	}

	return lastState, MessageOk(), nil
}

func ExecuteOnDevice(device Device, on bool) error {
	if device.Type == DeviceTypeKasa {
		return service.SetKasaDeviceStatus(device.ID, on)
	} else {
		return errors.New("unable to handle device type")
	}
}
