package intents

import (
	"context"
	"main/core"
)

type Stop struct {
}

func (p *Stop) IntentLabel() string {
	return "stop"
}

func (p *Stop) IntentParsePrompt(req *core.Request) string {
	return ""
}

func (p *Stop) IntentParseReceiver() any {
	return nil
}

func (p *Stop) ActOnIntent(ctx context.Context, req *core.Request, md *core.IntentMetadata) (core.Response, error) {
	return core.Response{
		OutboundMessage: core.OutboundMessage{
			Action: core.ActionStop,
		},
	}, nil
}
