package hal9000

import (
	"fmt"
	"os"
)

type Interface interface {
	Name() string
	SendMessage(message ResponseMessage) error
}

type InterfaceTypeSMS struct {
	Number string
}

func (i InterfaceTypeSMS) Name() string {
	return "sms"
}

func (i InterfaceTypeSMS) SendMessage(m ResponseMessage) error {
	fmt.Println(m.Text)
	return nil
}

type InterfaceTypeTerminal struct {
	Output *os.File
}

func (i InterfaceTypeTerminal) Name() string {
	return "terminal"
}

func (i InterfaceTypeTerminal) SendMessage(m ResponseMessage) error {
	i.Output.Write([]byte(m.Text + "\n"))
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
