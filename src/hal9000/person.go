package hal9000

import (
	"encoding/json"
	"errors"
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

func (p Person) GetNames() []string {
	return p.Names
}

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
	nameables := make([]util.Nameable, len(people))
	for i, p := range people {
		nameables[i] = p
	}
	sortedNameables := util.GenerateNameableSequence(nameables)
	lcName := strings.ToLower(name)
	for _, nameable := range sortedNameables {
		if strings.ToLower(nameable.Name) == lcName {
			return nameable.Nameable.(Person), nil
		}
	}
	return Person{}, ErrorPersonNotFound
}

func SendMessageToPerson(recipient Person, message util.ResponseMessage) error {
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
		err := ses.BreakIn(message)
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
