package service

import (
	"encoding/json"
	"errors"
	"hal9000/types"
	"hal9000/util"
	"io/ioutil"
	"os"
	"strings"
)

var ErrorPersonNotFound = errors.New("person not found")

type PersonConcrete struct {
	ID          string   `json:"id"`
	Names       []string `json:"names"`
	PhoneNumber string   `json:"phoneNumber"`
}

func (p PersonConcrete) GetID() string {
	return p.ID
}

func (p PersonConcrete) GetNames() []string {
	return p.Names
}

func (p PersonConcrete) GetOriginName() string {
	return p.Names[0]
}

func (p PersonConcrete) GetPhoneNumber() string {
	return p.PhoneNumber
}

type personProviderConcrete struct {
	people []types.Person
}

func InitPersonProvider() (types.PersonProvider, error) {
	bytes, err := ioutil.ReadFile(os.Getenv("PEOPLE_MANIFEST_PATH"))
	if err != nil {
		return nil, err
	}
	var personsConcrete []PersonConcrete
	err = json.Unmarshal(bytes, &personsConcrete)
	if err != nil {
		return nil, err
	}
	people := make([]types.Person, len(personsConcrete))
	for i, p := range personsConcrete {
		people[i] = p
	}

	// for _, p := range people {
	// RegisterTransientInterface(p, InterfaceTypeSMS{p.PhoneNumber}) TODO
	// }

	return personProviderConcrete{people}, nil
}

func (pp personProviderConcrete) People() []types.Person {
	return pp.people
}

func (pp personProviderConcrete) GetPersonByName(name string) (types.Person, error) {
	nameables := make([]types.Nameable, len(pp.people))
	for i, p := range pp.people {
		nameables[i] = p
	}
	sortedNameables := util.GenerateNameableSequence(nameables)
	lcName := strings.ToLower(name)
	for _, nameable := range sortedNameables {
		if strings.ToLower(nameable.Name) == lcName {
			return nameable.Nameable.(types.Person), nil
		}
	}
	return nil, ErrorPersonNotFound
}

func (pp personProviderConcrete) SendMessageToPerson(recipient types.Person, message types.ResponseMessage) error {
	// sessions := GetUserSessions(recipient)
	// if len(sessions) == 0 {
	// 	ics := GetInterfacesForPerson(recipient, "")
	// 	if len(ics) == 0 {
	// 		return ErrorNoInterfacesAvailable(recipient)
	// 	}
	// 	for _, ic := range ics {
	// 		sessions = append(sessions, NewSession(recipient, ic))
	// 	}
	// }
	// for _, ses := range sessions {
	// 	err := ses.BreakIn(message)
	// 	if err != nil {
	// 		return err
	// 	}
	// } TODO
	return nil
}

func (pp personProviderConcrete) GetPersonByID(id string) (types.Person, error) {
	for _, person := range pp.people {
		if person.GetID() == id {
			return person, nil
		}
	}
	return nil, ErrorPersonNotFound
}
