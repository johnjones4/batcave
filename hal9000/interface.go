package hal9000

import (
	"crypto/sha1"
	"fmt"
)

type Interface interface {
	Type() string
	ID() string
	IsStillValid() bool
	SendMessage(message ResponseMessage) error
}

type InterfaceTypeSMS struct {
	Number string
}

func (i InterfaceTypeSMS) Type() string {
	return "sms"
}

func (i InterfaceTypeSMS) ID() string {
	h := sha1.New()
	h.Write([]byte(i.Number))
	bs := h.Sum(nil)
	return fmt.Sprintf("sms-%x", bs)
}

func (i InterfaceTypeSMS) IsStillValid() bool {
	return true
}

func (i InterfaceTypeSMS) SendMessage(m ResponseMessage) error {
	fmt.Println(m.Text) //TODO
	return nil
}

var transientInterfaceStore map[string][]Interface = make(map[string][]Interface)

func RegisterTransientInterface(person Person, iface Interface) {
	if _, ok := transientInterfaceStore[person.ID]; !ok {
		transientInterfaceStore[person.ID] = make([]Interface, 0)
	}
	transientInterfaceStore[person.ID] = append(transientInterfaceStore[person.ID], iface)
}

func GetInterfacesForPerson(p Person, id string) []Interface {
	interfaces := make([]Interface, 0)
	if transInterfaces, ok := transientInterfaceStore[p.ID]; ok {
		removeSet := make([]int, 0)
		for i, iface := range transInterfaces {
			if iface.IsStillValid() {
				if id == "" || (id != "" && id == iface.ID()) {
					interfaces = append(interfaces, iface)
				}
			} else {
				removeSet = append(removeSet, i)
			}
		}
		if len(removeSet) > 0 {
			for _, i := range removeSet {
				transInterfaces = append(transInterfaces[:i], transInterfaces[i+1:]...)
			}
			transientInterfaceStore[p.ID] = transInterfaces
		}
	}
	return interfaces
}

func ErrorNoInterfacesAvailable(p Person) error {
	return fmt.Errorf("no interfaces ready for %s", p.Names[0])
}
