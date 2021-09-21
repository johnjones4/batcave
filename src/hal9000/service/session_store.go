package service

import (
	"hal9000/types"
	"hal9000/util"
)

type sessionStoreConcrete struct {
	sessions []*types.Session
}

func InitSessionStore() types.SessionStore {
	return &sessionStoreConcrete{make([]*types.Session, 0)}
}

func (ss *sessionStoreConcrete) SaveSession(ses *types.Session) {
	for i, ses1 := range ss.sessions {
		if ses1.ID == ses.ID {
			ss.sessions[i] = ses
			return
		}
	}
	ss.sessions = append(ss.sessions, ses)
}

func (ss *sessionStoreConcrete) GetVisualUserSessions(p *types.Person) []*types.Session {
	allSessions := ss.GetUserSessions(p)
	visualSessions := make([]*types.Session, 0)
	for _, s := range allSessions {
		if (*s.Interface).SupportsVisuals() {
			visualSessions = append(visualSessions, s)
		}
	}
	return visualSessions
}

func (ss *sessionStoreConcrete) GetUserSessions(p *types.Person) []*types.Session {
	usessions := make([]*types.Session, 0)
	unregs := make([]int, 0)
	for i, ses := range ss.sessions {
		if (*ses.Caller).GetID() == (*p).GetID() {
			if (*ses.Interface).IsStillValid() {
				usessions = append(usessions, ses)
			} else {
				unregs = append(unregs, i)
			}
		}
	}
	// if len(unregs) > 0 {
	// newSessions := ss.sessions
	// for _, i := range unregs {
	// 	if i < len(newSessions) {
	// 		newSessions = append(newSessions[:i], newSessions[i+1:]...)
	// 	}
	// }
	// ss.sessions = newSessions TODO
	// }
	return usessions
}

func (ss *sessionStoreConcrete) GetSessionWithInterfaceID(id string) (*types.Session, error) {
	for _, ses := range ss.sessions {
		if (*ses.Interface).ID() == id {
			return ses, nil
		}
	}
	return nil, util.ErrorSessionNotFound
}

func (ss *sessionStoreConcrete) GetSessionById(id string) (*types.Session, error) {
	for _, ses := range ss.sessions {
		if ses.ID == id {
			return ses, nil
		}
	}
	return nil, util.ErrorSessionNotFound
}
