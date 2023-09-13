package processors

import (
	"context"
	"encoding/json"
	"fmt"
	"main/core"
)

type confirmMessageReceiver struct {
	Response string `json:"response"`
}

func (p *Processors) ConfirmMessage(ctx context.Context, req *core.Request, res *core.Response) error {
	if res.Message.Text != "" {
		return nil
	}
	prompt := fmt.Sprintf("If I ask you to execute the following command, \"%s\" what would an affirmative response look like in the JSON format: {\"response\":\"\"}", req.Message.Text)
	text, err := p.LLM.Completion(ctx, prompt)
	if err != nil {
		return err
	}

	var receiver confirmMessageReceiver
	err = json.Unmarshal([]byte(text), &receiver)
	if err != nil {
		return err
	}

	res.Message.Text = receiver.Response

	return nil
}
