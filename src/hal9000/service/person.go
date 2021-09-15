package service

import (
	"encoding/json"
	"hal9000/types"
	"hal9000/util"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

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
	return nil, util.ErrorPersonNotFound
}

func (pp personProviderConcrete) SendMessageToPerson(runtime types.Runtime, recipient types.Person, message types.ResponseMessage) error {
	sessions := runtime.SessionStore().GetUserSessions(recipient)
	if len(sessions) == 0 {
		ics := runtime.InterfaceStore().GetInterfacesForPerson(recipient, "")
		if len(ics) == 0 {
			return util.ErrorNoInterfacesAvailable(recipient)
		}
		for _, ic := range ics {
			ses := types.Session{
				Caller:      recipient,
				ID:          uuid.NewString(),
				Start:       time.Now(),
				Interface:   ic,
				StateString: util.StateTypeDefault,
			}
			runtime.SessionStore().SaveSession(ses)
			sessions = append(sessions, ses)
		}
	}
	for _, ses := range sessions {
		runtime.Logger().LogEvent("break_in", map[string]interface{}{
			"session": ses.ID,
			"message": message,
		})
		err := ses.Interface.SendMessage(message)
		if err != nil {
			runtime.Logger().LogError(err)
		}
	}
	return nil
}

func (pp personProviderConcrete) GetPersonByID(id string) (types.Person, error) {
	for _, person := range pp.people {
		if person.GetID() == id {
			return person, nil
		}
	}
	return nil, util.ErrorPersonNotFound
}
