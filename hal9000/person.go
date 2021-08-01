package hal9000

import (
	"encoding/json"
	"errors"
	"fmt"
	"hal9000/util"
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
	ics := GetInterfacesForPerson(recipient, "")
	if len(ics) == 0 {
		return ErrorNoInterfacesAvailable(recipient)
	}
	message := fmt.Sprintf("Message from %s: \"%s\"", sender.Names[0], m)
	for _, ic := range ics {
		util.LogEvent("break_in", map[string]interface{}{
			"from": sender.ID,
			"to":   recipient.ID,
		})
		err := ic.SendMessage(ResponseMessage{message, "", nil})
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
