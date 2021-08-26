package hal9000

import (
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
	Interface   Interface            `json:"interface"`
}

func NewSession(caller Person, ic Interface) Session {
	ses := Session{
		Caller:      caller,
		ID:          uuid.NewString(),
		Start:       time.Now(),
		Interface:   ic,
		StateString: util.StateTypeDefault,
		History:     make([]HistoricalExchange, 0),
	}
	SaveSession(ses)
	return ses
}

func (s *Session) State() State {
	return InitStateByName(s.StateString)
}

func (s *Session) ProcessIncomingMessage(m RequestMessage) (ResponseMessage, error) {
	requestTime := time.Now()

	nextState, response, err := s.State().ProcessIncomingMessage(s.Caller, m)
	if err != nil {
		fmt.Println(err) //todo error logging
		return MessageError(err), nil
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

func (s *Session) BreakIn(m ResponseMessage) error {
	util.LogEvent("break_in", map[string]interface{}{
		"session": s.ID,
		"message": m,
	})
	return s.Interface.SendMessage(m)
}
