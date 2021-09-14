package types

type Runtime interface {
	Devices() DeviceProvider
	Agenda() AgendaProvider
	People() PersonProvider
	Displays() DisplayablesProvider
	Kasa() KasaProvider
	Jobs() JobProvider
	Weather() WeatherProvider
	Google() GoogleProvider
	KVStore() KVStore
	Logger() Logger
	AlertQueue() AlertQueue
	SessionStore() SessionStore
	InterfaceStore() InterfaceStore
}
