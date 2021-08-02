package hal9000

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Person struct {
	ID          string   `json:"id"`
	Names       []string `json:"names"`
	PhoneNumber string   `json:"phoneNumber"`
}

var people []Person

var ErrorPersonNotFound = errors.New("person not found")

func InitPeople() error {
	bytes, err := ioutil.ReadFile(os.Getenv("PEOPLE_MANIFEST_PATH"))
	if err != nil {
		return err
	}
	people = nil
	err = json.Unmarshal(bytes, &people)
	if err != nil {
		return err
	}

	for _, p := range people {
		RegisterTransientInterface(p, InterfaceTypeSMS{p.PhoneNumber})
	}

	return nil
}

func GetPersonByName(name string) (Person, error) {
	lcName := strings.ToLower(name)
	for _, person := range people {
		for _, pName := range person.Names {
			if strings.ToLower(pName) == lcName {
				return person, nil
			}
		}
	}
	return Person{}, ErrorPersonNotFound
}

func SendMessageToPerson(sender Person, recipient Person, m string) error {
	message := fmt.Sprintf("Message from %s: \"%s\"", sender.Names[0], m)
	sessions := GetUserSessions(recipient)
	if len(sessions) == 0 {
		ics := GetInterfacesForPerson(recipient, "")
		if len(ics) == 0 {
			return ErrorNoInterfacesAvailable(recipient)
		}
		for _, ic := range ics {
			sessions = append(sessions, NewSession(recipient, ic))
		}
	}
	for _, ses := range sessions {
		err := ses.BreakIn(ResponseMessage{message, "", nil})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetPersonByID(id string) (Person, error) {
	for _, person := range people {
		if person.ID == id {
			return person, nil
		}
	}
	return Person{}, ErrorPersonNotFound
}
