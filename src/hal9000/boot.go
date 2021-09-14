package hal9000

import (
	"fmt"
	"hal9000/service"
	"hal9000/types"
)

type LogStep struct {
	Name string
	Fn   func() error
}

type runtimeConcrete struct {
	logger         types.Logger
	kvStore        types.KVStore
	parser         types.ParserProvider
	people         types.PersonProvider
	displayables   types.DisplayablesProvider
	jobs           types.JobProvider
	devices        types.DeviceProvider
	agenda         types.AgendaProvider
	kasa           types.KasaProvider
	google         types.GoogleProvider
	weather        types.WeatherProvider
	alertQueue     types.AlertQueue
	sessionStore   types.SessionStore
	interfaceStore types.InterfaceStore
}

func (rt *runtimeConcrete) Devices() types.DeviceProvider {
	return rt.devices
}

func (rt *runtimeConcrete) Agenda() types.AgendaProvider {
	return rt.agenda
}

func (rt *runtimeConcrete) People() types.PersonProvider {
	return rt.people
}

func (rt *runtimeConcrete) Displays() types.DisplayablesProvider {
	return rt.displayables
}

func (rt *runtimeConcrete) Kasa() types.KasaProvider {
	return rt.kasa
}

func (rt *runtimeConcrete) Jobs() types.JobProvider {
	return rt.jobs
}

func (rt *runtimeConcrete) Weather() types.WeatherProvider {
	return rt.weather
}

func (rt *runtimeConcrete) Google() types.GoogleProvider {
	return rt.google
}

func (rt *runtimeConcrete) KVStore() types.KVStore {
	return rt.kvStore
}

func (rt *runtimeConcrete) Logger() types.Logger {
	return rt.logger
}

func (rt *runtimeConcrete) AlertQueue() types.AlertQueue {
	return rt.alertQueue
}

func (rt *runtimeConcrete) SessionStore() types.SessionStore {
	return rt.sessionStore
}

func (rt *runtimeConcrete) InterfaceStore() types.InterfaceStore {
	return rt.interfaceStore
}

func (rt *runtimeConcrete) Parser() types.ParserProvider {
	return rt.parser
}

func BootUp() (types.Runtime, error) {
	rt := runtimeConcrete{}

	fns := [](LogStep){
		LogStep{"logger", func() error {
			logger, err := service.InitLogger()
			if err != nil {
				return err
			}
			rt.logger = logger
			return nil
		}},
		LogStep{"kv store", func() error {
			kvStore, err := service.InitFileKVStore()
			if err != nil {
				return err
			}
			rt.kvStore = kvStore
			return nil
		}},
		LogStep{"message parser", func() error {
			parser, err := service.InitParserProvider()
			if err != nil {
				return err
			}
			rt.parser = parser
			return nil
		}},
		LogStep{"people", func() error {
			people, err := service.InitPersonProvider()
			if err != nil {
				return err
			}
			rt.people = people
			return nil
		}},
		LogStep{"displayables", func() error {
			displayables, err := service.InitDisplayablesProvider()
			if err != nil {
				return err
			}
			rt.displayables = displayables
			return nil
		}},
		LogStep{"jobs", func() error {
			// jobs, err := service.InitJobProvider()
			// if err != nil {
			// 	return err
			// }
			// rt.jobs = jobs
			return nil
		}},
		LogStep{"devices", func() error {
			devices, err := service.InitDeviceProvider()
			if err != nil {
				return err
			}
			rt.devices = devices
			return nil
		}},
		LogStep{"agenda", func() error {
			agenda, err := service.InitAgendaProvider()
			if err != nil {
				return err
			}
			rt.agenda = agenda
			return nil
		}},
		LogStep{"kasa", func() error {
			kasa, err := service.InitKasaProvider()
			if err != nil {
				return err
			}
			rt.kasa = kasa
			return nil
		}},
		LogStep{"google", func() error {
			google, err := service.InitGoogleProvider(&rt)
			if err != nil {
				return err
			}
			rt.google = google
			return nil
		}},
		LogStep{"weather", func() error {
			weather, err := service.InitWeatherProvider(&rt)
			if err != nil {
				return err
			}
			rt.weather = weather
			return nil
		}},
		LogStep{"alerts", func() error {
			alertQueue, err := service.InitAlertQueue(&rt)
			if err != nil {
				return err
			}
			rt.alertQueue = alertQueue
			return nil
		}},
		LogStep{"session store", func() error {
			rt.sessionStore = service.InitSessionStore()
			return nil
		}},
		LogStep{"interface store", func() error {
			rt.interfaceStore = service.InitInterfaceStore()
			for _, person := range rt.People().People() {
				sms := InterfaceTypeSMS{person.GetPhoneNumber()}
				rt.interfaceStore.Register(person, sms)
			}
			return nil
		}},
	}
	for _, fn := range fns {
		fmt.Printf("Initializing %s ... ", fn.Name)
		err := fn.Fn()
		if err != nil {
			return nil, err
		}
		fmt.Println("done")
	}
	return &rt, nil
}
