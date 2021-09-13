package hal9000

import (
	"fmt"
	"hal9000/service"
	"hal9000/util"
	"os"
	"strings"
)

type LogStep struct {
	Name string
	Fn   func() error
}

type BackgroundJob struct {
	Name string
	Fn   func(*chan util.ResponseMessage)
}

func BootUp() error {
	fns := [](LogStep){
		LogStep{"logger", util.InitLogger},
		LogStep{"kv store", util.InitFileKVStore},
		LogStep{"message parser", InitializeDefaultIncomingMessageParser},
		LogStep{"people", InitPeople},
		LogStep{"displayables", InitDisplayables},
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

	alertedUserNames := strings.Split(os.Getenv("ALERTED_USER_NAMES"), ",")
	users := make([]Person, len(alertedUserNames))
	for i, alertedUserName := range alertedUserNames {
		person, err := GetPersonByName(alertedUserName)
		if err != nil {
			return err
		}
		users[i] = person
	}

	bgs := [](BackgroundJob){
		BackgroundJob{"google token refresh", service.StartGoogleTokenRefreshCycle},
		BackgroundJob{"weather alert scanner", service.StartWeatherAlertLoop},
	}
	alertChan := make(chan util.ResponseMessage)
	for _, bg := range bgs {
		fmt.Printf("Starting up %s ... ", bg.Name)
		go bg.Fn(&alertChan)
		fmt.Println("done")
	}

	go (func() {
		for {
			alert := <-alertChan
			for _, user := range users {
				err := SendMessageToPerson(user, alert)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	})()

	return nil
}
