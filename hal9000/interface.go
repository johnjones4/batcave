package hal9000

import (
	"fmt"
)

type Interface interface {
	SendMessage(message ResponseMessage) error
}

type InterfaceTypeSMS struct {
	Number string
}

type InterfaceTypeTerminal struct {
}

func (i InterfaceTypeSMS) SendMessage(m ResponseMessage) error {
	fmt.Println(m.Text)
	return nil
}

func GetInterfacesForPerson(p Person) []Interface {
	interfaces := make([]Interface, 0)
	if p.PhoneNumber != "" {
		interfaces = append(interfaces, InterfaceTypeSMS{p.PhoneNumber})
	}
	return interfaces
}

func ErrorNoInterfacesAvailable(p Person) error {
	return fmt.Errorf("no interfaces ready for %s", p.Names[0])
}
