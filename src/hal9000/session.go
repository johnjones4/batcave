package hal9000

import (
	"hal9000/types"
	"hal9000/util"
	"time"

	"github.com/google/uuid"
)

type SessionLogInfo struct {
	Timestamp  time.Time             `json:"timestamp"`
	Request    types.RequestMessage  `json:"request"`
	Response   types.ResponseMessage `json:"response"`
	StartState string                `json:"startState"`
	EndState   string                `json:"endState"`
}

func NewSession(runtime *types.Runtime, caller *types.Person, ic *types.Interface) types.Session {
	ses := types.Session{
		Caller:      caller,
		ID:          uuid.NewString(),
		Start:       time.Now(),
		Interface:   ic,
		StateString: util.StateTypeDefault,
	}
	(*(*runtime).SessionStore()).SaveSession(&ses)
	return ses
}

func ProcessIncomingMessage(runtime *types.Runtime, s *types.Session, m types.RequestMessage) (types.ResponseMessage, error) {
	requestTime := time.Now()

	nextState, response, err := initStateByName(s.StateString).ProcessIncomingMessage(runtime, s.Caller, m)
	if err != nil {
		(*(*runtime).Logger()).LogError(err)
		return util.MessageError(err), nil
	}

	(*(*runtime).Logger()).LogEvent("exchange", map[string]interface{}{
		"session": s.ID,
		"exchange": SessionLogInfo{
			Timestamp:  requestTime,
			Request:    m,
			Response:   response,
			StartState: s.StateString,
			EndState:   (*nextState).Name(),
		},
	})

	s.StateString = (*nextState).Name()

	if !(*s.Interface).SupportsVisuals() && response.URL != "" {
		sessions := (*(*runtime).SessionStore()).GetVisualUserSessions(s.Caller)
		if len(sessions) == 0 {
			ics := (*(*runtime).InterfaceStore()).GetVisualInterfacesForPerson(s.Caller)
			for _, ic := range ics {
				ses := NewSession(runtime, s.Caller, ic)
				sessions = append(sessions, &ses)
			}
		}
		for _, ses := range sessions {
			if (*ses.Interface).SupportsVisuals() {
				(*(*runtime).Logger()).LogEvent("visual_override", map[string]interface{}{
					"session": ses.ID,
					"message": response,
				})
				err := (*ses.Interface).SendMessage(response)
				if err != nil {
					(*(*runtime).Logger()).LogError(err)
				}
			}
		}
	}

	return response, nil
}
