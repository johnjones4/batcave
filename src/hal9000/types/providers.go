package types

import "time"

type Logger interface {
	LogEvent(event string, info interface{})
	LogError(err error)
}

type PersonProvider interface {
	People() []*Person
	GetPersonByName(name string) (*Person, error)
	SendMessageToPerson(runtime *Runtime, recipient *Person, message ResponseMessage) error
	GetPersonByID(id string) (*Person, error)
}

type DeviceProvider interface {
	Devices() []*Device
	FindDeviceInString(str string) (*Device, error)
}

type AgendaProvider interface {
	GetAgendaForDateRange(start time.Time, end time.Time) ([]Event, error)
}

type DisplayablesProvider interface {
	FindDisplayableInString(str string) (*Displayable, error)
}

type KasaProvider interface {
	SetKasaDeviceStatus(id string, on bool) error
}

type JobProvider interface {
	FindJobById(id string) (*Job, error)
	ReportJobStatus(runtime *Runtime, job *Job, info *JobStatusInfo) error
}

type WeatherProvider interface {
	MakeWeatherAPIAlertCall(lat float64, lon float64) ([]NOAAWeatherAlertFeature, error)
	MakeWeatherAPIForecastCall(lat float64, lon float64, date time.Time) (string, string, error)
	MakeWeatherAPIPointRequest(lat float64, lon float64) (NOAAWeatherPointProperties, error)
	DefaultLatLon() (float64, float64)
}

type GoogleProvider interface {
	RefreshAuthToken(runtime *Runtime) error
	CreateNewEvent(runtime *Runtime, event Event) error
}

type ParserProvider interface {
	ProcessMessage(input RequestMessage) (ParsedRequestMessage, error)
}
