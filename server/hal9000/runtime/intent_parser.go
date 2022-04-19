package runtime

import (
	"errors"
	"strings"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
)

var (
	ErrorInputEmpty = errors.New("empty input")
	ErrorNoIntent   = errors.New("no intent found")
)

func (r *Runtime) Parse(in core.InboundBody, client core.Client, state string) (core.Inbound, error) {
	if len(in.Body) == 0 {
		return core.Inbound{}, ErrorInputEmpty
	}

	var command string
	var parseType string
	if in.Body[0] == '/' {
		parseType = "explicit"
		firstSpace := strings.Index(in.Body, " ")
		if firstSpace < 0 {
			command = in.Body[1:]
			in.Body = ""
		} else {
			command = strings.TrimSpace(in.Body[1:firstSpace])
			in.Body = strings.TrimSpace(in.Body[firstSpace:])
		}
	} else {
		parseType = "inferred"
		var err error
		command, err = r.Predictor.PredictIntent(in.Body)
		if err != nil {
			return core.Inbound{}, err
		}
	}

	if command == "" {
		return core.Inbound{}, ErrorNoIntent
	}

	var user core.User
	if len(client.Users) == 1 {
		userRec, err := r.UserStore.GetUser(client.Users[0])
		if err != nil {
			return core.Inbound{}, err
		}
		user = userRec.User
	}

	request := core.Inbound{
		InboundBody: in,
		Command:     command,
		State:       state,
		Client:      client,
		User:        user,
		ParseType:   parseType,
	}

	return request, nil
}
