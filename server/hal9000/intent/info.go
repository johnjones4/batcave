package intent

import (
	"fmt"
	"strings"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
)

type Info struct {
	Intents []core.Intent
}

func (c *Info) SupportedComandsForState(s core.State) map[string]string {
	return map[string]string{
		"commands": "Get a list of currently available commands.",
	}
}

func (c *Info) Execute(req core.Inbound) (core.Outbound, error) {
	lines := make([]string, 0)
	for _, intent := range c.Intents {
		for command, description := range intent.SupportedComandsForState(req.State) {
			lines = append(lines, fmt.Sprintf("/%s: %s", command, description))
		}
	}
	return core.Outbound{
		OutboundBody: core.OutboundBody{
			Body: "Available commands:\n" + strings.Join(lines, "\n"),
		},
		State: req.State,
	}, nil
}
