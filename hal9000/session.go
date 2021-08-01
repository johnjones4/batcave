package hal9000

import (
	"encoding/json"
	"errors"
	"fmt"
	"hal9000/util"
	"time"

	"github.com/google/uuid"
)

type HistoricalExchange struct {
	Timestamp  time.Time       `json:"timestamp"`
	Request    RequestMessage  `json:"request"`
	Response   ResponseMessage `json:"response"`
	StartState string          `json:"startState"`
	EndState   string          `json:"endState"`
}

type Session struct {
	Caller      Person               `json:"caller"`
	ID          string               `json:"id"`
	Start       time.Time            `json:"start"`
	StateString string               `json:"state"`
	History     []HistoricalExchange `json:"history"`
	InterfaceID string               `json:"interface"`
	Interface   Interface
}

var (
	ErrorSessionNotFound = errors.New("session not found")
)

func NewSession(caller Person, ic Interface) (Session, error) {
	ses := Session{
		Caller:      caller,
		ID:          uuid.NewString(),
		Start:       time.Now(),
		Interface:   ic,
		InterfaceID: ic.ID(),
		StateString: util.StateTypeDefault,
		History:     make([]HistoricalExchange, 0),
	}
	err := ses.Save()
	if err != nil {
		return Session{}, err
	}
	return ses, nil
}

func SessionKeyForID(id string) string {
	return fmt.Sprintf("session_%s", id)
}

func LoadSession(id string) (Session, error) {
	key := SessionKeyForID(id)
	bytes := util.KVStoreInstance.GetBytes(key, []byte{})
	if len(bytes) == 0 {
		return Session{}, ErrorSessionNotFound
	}
	var ses Session
	err := json.Unmarshal(bytes, &ses)
	if err != nil {
		return Session{}, err
	}
	interfaces := GetInterfacesForPerson(ses.Caller, ses.InterfaceID)
	if len(interfaces) == 0 {
		return Session{}, errors.New("interface for session not found")
	}
	ses.Interface = interfaces[0]
	return ses, nil
}

func (s *Session) Save() error {
	sessionData, err := json.Marshal(s)
	if err != nil {
		return err
	}
	key := SessionKeyForID(s.ID)
	return util.KVStoreInstance.SetBytes(key, sessionData, time.Now().Add(time.Hour))
}

func (s *Session) State() State {
	return InitStateByName(s.StateString)
}

func (s *Session) ProcessIncomingMessage(m RequestMessage) (ResponseMessage, error) {
	requestTime := time.Now()

	nextState, response, err := s.State().ProcessIncomingMessage(s.Caller, m)
	if err != nil {
		return ResponseMessage{}, err
	}

	exchange := HistoricalExchange{
		Timestamp:  requestTime,
		Request:    m,
		Response:   response,
		StartState: s.StateString,
		EndState:   nextState.Name(),
	}
	s.History = append(s.History, exchange)

	util.LogEvent("exchange", map[string]interface{}{
		"session":  s.ID,
		"exchange": exchange,
	})

	s.StateString = nextState.Name()

	return response, nil
}
