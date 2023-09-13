package intents

import (
	"context"
	"fmt"
	"main/core"
)

type Unknown struct {
}

type unknownIntentParseReceiver struct {
	Answer string `json:"answer"`
}

func (p *Unknown) IntentLabel() string {
	return "unknown"
}

func (p *Unknown) IntentParsePrompt(req *core.Request) string {
	return fmt.Sprintf("Provide an answer to the statement \"%s\" in the JSON format {\"answer\":\"\"}", req.Message.Text)
}

func (p *Unknown) IntentParseReceiver() any {
	return unknownIntentParseReceiver{}
}

func (p *Unknown) ActOnIntent(ctx context.Context, req *core.Request, md *core.IntentMetadata) (core.Response, error) {
	return core.Response{}, nil //TODO
}
