package hal9000

import (
	"fmt"
	"hal9000/service"
	"hal9000/util"
)

type LogStep struct {
	Name string
	Fn   func() error
}

func BootUp() error {
	fns := [](LogStep){
		LogStep{"logger", util.InitLogger},
		LogStep{"kv store", util.InitFileKVStore},
		LogStep{"message parser", InitializeDefaultIncomingMessageParser},
		LogStep{"people", InitPeople},
		LogStep{"displayables", InitDisplays},
		LogStep{"devices", InitDevices},
		LogStep{"calendars", InitCalendarSchedules},
		LogStep{"kasa", service.InitKasaConnection},
	}
	for _, fn := range fns {
		fmt.Printf("Initializing %s ... ", fn.Name)
		err := fn.Fn()
		if err != nil {
			return err
		}
		fmt.Println("done")
	}

	return nil
}
