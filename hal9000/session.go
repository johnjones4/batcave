package hal9000

import (
	"time"
)

type Session interface {
	StartTime() time.Time
	Interface() Interface
	State() State
	BreakIn(p Message) error
	ProcessIncomingMessage(m string) (Message, error)
}

func GetActiveSessions(p Person) ([]Session, error) {
	return []Session{}, nil //TODO
}

func InitiateNewSession(ic Interface) (Session, error) {
	//TODO session store
	return &SessionConcrete{
		IStart:   time.Now(),
		IContext: ic,
		IState:   DefaultState{},
	}, nil
}

type SessionConcrete struct {
	IStart   time.Time
	IContext Interface
	IState   State `json:"state"`
}

func (s *SessionConcrete) StartTime() time.Time {
	return s.IStart
}

func (s *SessionConcrete) Interface() Interface {
	return s.IContext
}

func (s *SessionConcrete) State() State {
	return s.IState
}

func (s *SessionConcrete) BreakIn(p Message) error {
	s.IContext.SendMessage(p) //TODO
	return nil
}

func (s *SessionConcrete) ProcessIncomingMessage(m string) (Message, error) {
	nextState, response, err := s.IState.ProcessIncomingMessage(m)
	if err != nil {
		return Message{}, err
	}

	s.IState = nextState

	return response, nil
}
