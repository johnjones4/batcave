package hal9000

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type Person struct {
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

func SendMessageToPerson(person Person, m string) error {
	sessions, err := GetActiveSessions(person)
	if err != nil {
		return err
	}
	if len(sessions) == 0 {
		ics := GetInterfacesForPerson(person)
		if len(ics) == 0 {
			return ErrorNoInterfacesAvailable(person)
		}
		for _, ic := range ics {
			ses, err := InitiateNewSession(ic)
			if err != nil {
				return err
			}
			ses.BreakIn(Message{m, "", nil})
		}
		return nil
	}
	for _, ses := range sessions {
		err = ses.BreakIn(Message{m, "", nil})
		if err != nil {
			return err
		}
	}

	return nil
}
