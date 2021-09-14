package hal9000

import (
	"fmt"
	"hal9000/types"
	"hal9000/util"
	"time"

	"github.com/google/uuid"
)

type HistoricalExchange struct {
	Timestamp  time.Time             `json:"timestamp"`
	Request    types.RequestMessage  `json:"request"`
	Response   types.ResponseMessage `json:"response"`
	StartState string                `json:"startState"`
	EndState   string                `json:"endState"`
}

func NewSession(caller types.Person, ic types.Interface) Session {
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

func (s *Session) State() types.State {
	return InitStateByName(s.StateString)
}

func (s *Session) ProcessIncomingMessage(m types.RequestMessage) (types.ResponseMessage, error) {
	requestTime := time.Now()

	nextState, response, err := s.State().ProcessIncomingMessage(s.Caller, m)
	if err != nil {
		fmt.Println(err) //todo error logging
		return util.MessageError(err), nil
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

	if !s.Interface.SupportsVisuals() && response.URL != "" {
		sessions := GetUserSessions(s.Caller)
		if len(sessions) == 0 {
			ics := GetVisualInterfacesForPerson(s.Caller)
			for _, ic := range ics {
				sessions = append(sessions, NewSession(s.Caller, ic))
			}
		}
		for _, ses := range sessions {
			if ses.Interface.SupportsVisuals() {
				err := ses.BreakIn(response)
				if err != nil {
					fmt.Println(err) //todo error logging
				}
			}
		}
	}

	return response, nil
}

func (s *Session) BreakIn(m types.ResponseMessage) error {
	util.LogEvent("break_in", map[string]interface{}{
		"session": s.ID,
		"message": m,
	})
	return s.Interface.SendMessage(m)
}
