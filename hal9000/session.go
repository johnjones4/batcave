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
	ID            string               `json:"id"`
	Start         time.Time            `json:"start"`
	StateString   string               `json:"state"`
	History       []HistoricalExchange `json:"history"`
	InterfaceType string               `json:"interface"`
	Interface     Interface
}

var (
	ErrorSessionNotFound = errors.New("session not found")
)

func GetActiveSessions(p Person) ([]Session, error) {
	return []Session{}, nil //TODO
}

func NewSession(ic Interface) (Session, error) {
	ses := Session{
		ID:            uuid.NewString(),
		Start:         time.Now(),
		Interface:     ic,
		InterfaceType: ic.Name(),
		StateString:   util.StateTypeDefault,
		History:       make([]HistoricalExchange, 0),
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

func LoadSession(id string, ic Interface) (Session, error) {
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
	if ses.InterfaceType != ic.Name() {
		return Session{}, fmt.Errorf("Interface mismatch (Received %s, expected %s)", ses.InterfaceType, ic.Name())
	}
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

func (s *Session) BreakIn(p ResponseMessage) error {
	s.Interface.SendMessage(p) //TODO
	return nil
}

func (s *Session) State() State {
	return InitStateByName(s.StateString)
}

func (s *Session) ProcessIncomingMessage(m RequestMessage) (ResponseMessage, error) {
	requestTime := time.Now()

	nextState, response, err := s.State().ProcessIncomingMessage(m)
	if err != nil {
		return ResponseMessage{}, err
	}

	s.History = append(s.History, HistoricalExchange{
		Timestamp:  requestTime,
		Request:    m,
		Response:   response,
		StartState: s.StateString,
		EndState:   nextState.Name(),
	})

	//TODO write log

	s.StateString = nextState.Name()

	return response, nil
}
