package intent

import (
	"errors"
	"strings"

	"github.com/johnjones4/hal-9000/hal9000/core"
)

func Parse(in core.InboundBody, state core.State) (core.Inbound, error) {
	if len(in.Body) == 0 || in.Body[0] != '/' {
		return core.Inbound{}, errors.New("input not recognized") //TODO
	}

	firstSpace := strings.Index(in.Body, " ")
	var command string
	if firstSpace < 0 {
		command = in.Body[1:]
		in.Body = ""
	} else {
		command = strings.TrimSpace(in.Body[1:firstSpace])
		in.Body = strings.TrimSpace(in.Body[firstSpace:])
	}

	request := core.Inbound{
		InboundBody: in,
		Command:     command,
		State:       state,
	}

	return request, nil
}
