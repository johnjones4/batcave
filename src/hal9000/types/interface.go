package types

type Interface interface {
	Type() string
	ID() string
	IsStillValid() bool
	SupportsVisuals() bool
	SendMessage(message ResponseMessage) error
}

type InterfaceStore interface {
	Register(person Person, iface Interface)
	GetInterfacesForPerson(p Person, id string) []Interface
	GetVisualInterfacesForPerson(p Person) []Interface
	DetermineInterfaceOwner(runtime Runtime, iface Interface) (Person, error)
}
