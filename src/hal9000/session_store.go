package hal9000

import "errors"

var sessions = make([]Session, 0)

var ErrorSessionNotFound = errors.New("session not found")

func SaveSession(ses Session) {
	for i, ses1 := range sessions {
		if ses1.ID == ses.ID {
			sessions[i] = ses
			return
		}
	}
	sessions = append(sessions, ses)
}

func UnregisterSessions(ses []int) {
	newSessions := sessions
	for _, i := range ses {
		newSessions = append(newSessions[:i], newSessions[i+1:]...)
	}
	sessions = newSessions
}

func GetUserSessions(p Person) []Session {
	usessions := make([]Session, 0)
	unregs := make([]int, 0)
	for i, ses := range sessions {
		if ses.Caller.ID == p.ID {
			if ses.Interface.IsStillValid() {
				usessions = append(usessions, ses)
			} else {
				unregs = append(unregs, i)
			}
		}
	}
	if len(unregs) > 0 {
		UnregisterSessions(unregs)
	}
	return usessions
}

func GetSessionWithInterfaceID(id string) (Session, error) {
	for _, ses := range sessions {
		if ses.Interface.ID() == id {
			return ses, nil
		}
	}
	return Session{}, ErrorSessionNotFound
}

func GetSessionById(id string) (Session, error) {
	for _, ses := range sessions {
		if ses.ID == id {
			return ses, nil
		}
	}
	return Session{}, ErrorSessionNotFound
}
