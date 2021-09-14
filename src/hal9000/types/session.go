package types

import "time"

type Session struct {
	Caller      Person
	ID          string
	Start       time.Time
	StateString string
	Interface   Interface
}

type SessionStore interface {
	SaveSession(ses Session)
	GetUserSessions(p Person) []Session
	GetVisualUserSessions(p Person) []Session
	GetSessionWithInterfaceID(id string) (Session, error)
	GetSessionById(id string) (Session, error)
}
