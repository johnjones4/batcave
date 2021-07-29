package hal9000

import (
	"hal9000/service"
	"hal9000/util"
)

func BootUp() error {
	fns := [](func() error){
		util.InitKVStore,
		InitializeDefaultIncomingMessageParser,
		InitPeople,
		InitDisplays,
		InitDevices,
		service.InitKasaConnection,
	}
	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}

	return nil
}
