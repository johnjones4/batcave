package hal9000

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"hal9000/util"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Interface interface {
	Type() string
	ID() string
	IsStillValid() bool
	SupportsVisuals() bool
	SendMessage(message util.ResponseMessage) error
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

func (i InterfaceTypeSMS) SupportsVisuals() bool {
	return false
}

func (i InterfaceTypeSMS) SendMessage(m util.ResponseMessage) error {
	accountSid := os.Getenv("TWILIO_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	msgData := url.Values{}
	msgData.Set("To", i.Number)
	msgData.Set("From", os.Getenv("TWILIO_NUMBER_FROM"))
	msgData.Set("Body", m.Text)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err := client.Do(req)
	return err
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

func GetVisualInterfacesForPerson(p Person) []Interface {
	interfaces := make([]Interface, 0)
	for _, iface := range GetInterfacesForPerson(p, "") {
		if iface.SupportsVisuals() {
			interfaces = append(interfaces, iface)
		}
	}
	return interfaces
}

func DetermineOwnerOfInterface(iface Interface) (Person, error) {
	for owner, ifaces := range transientInterfaceStore {
		for _, _iface := range ifaces {
			if iface.ID() == _iface.ID() {
				return GetPersonByID(owner)
			}
		}
	}
	return Person{}, errors.New("no owner for interface")
}

func ErrorNoInterfacesAvailable(p Person) error {
	return fmt.Errorf("no interfaces ready for %s", p.Names[0])
}
