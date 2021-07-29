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
	Timestamp time.Time       `json:"timestamp"`
	Request   RequestMessage  `json:"request"`
	Response  ResponseMessage `json:"response"`
}

type Session struct {
	ID          string               `json:"id"`
	Start       time.Time            `json:"start"`
	Interface   Interface            `json:"interface"`
	StateString string               `json:"state"`
	History     []HistoricalExchange `json:"history"`
}

var (
	SessionNotFoundError = errors.New("session not found")
)

func GetActiveSessions(p Person) ([]Session, error) {
	return []Session{}, nil //TODO
}

func NewSession(ic Interface) (Session, error) {
	ses := Session{
		ID:          uuid.NewString(),
		Start:       time.Now(),
		Interface:   ic,
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
	bytes := util.GetKVValueBytes(key, []byte{})
	if len(bytes) == 0 {
		return Session{}, SessionNotFoundError
	}
	var ses Session
	err := json.Unmarshal(bytes, &ses)
	if err != nil {
		return Session{}, err
	}
	return ses, nil
}

func (s *Session) Save() error {
	sessionData, err := json.Marshal(s)
	if err != nil {
		return err
	}
	key := SessionKeyForID(s.ID)
	return util.SetKVValueBytes(key, sessionData, time.Now().Add(time.Hour))
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

	s.StateString = nextState.Name()

	s.History = append(s.History, HistoricalExchange{
		Timestamp: requestTime,
		Request:   m,
		Response:  response,
	})

	return response, nil
}
