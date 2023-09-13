package intents

import (
	"context"
	"fmt"
	"main/core"
	"main/services/push"
	"main/util"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Remind struct {
	Push *push.Push
}

type scheduleIntentParseReceiver struct {
	Date        string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`
}

func (i *Remind) IntentLabel() string {
	return "remind"
}

func (i *Remind) IntentParsePrompt(req *core.Request) string {
	return fmt.Sprintf("Extract the exact date and time relative to %s and a description from the phrase \"%s\" and return the information in the JSON format {\"date\":\"RFC3339 format\",\"description\":\"\"}", time.Now().String(), req.Message.Text)
}

func (i *Remind) IntentParseReceiver() any {
	return scheduleIntentParseReceiver{}
}

func (i *Remind) ActOnIntent(ctx context.Context, req *core.Request, md *core.IntentMetadata) (core.Response, error) {
	var info scheduleIntentParseReceiver
	err := mapstructure.Decode(md.IntentParseReceiver, &info)
	if err != nil {
		return core.ResponseEmpty, err
	}

	parsedDate, err := util.ParseLLMDate(info.Date)
	if err != nil {
		return core.ResponseEmpty, err
	}

	err = i.Push.SendLater(
		ctx,
		parsedDate,
		req.Source,
		req.ClientID,
		core.PushMessage{
			OutboundMessage: core.OutboundMessage{
				EventId: req.EventId,
				Message: core.Message{
					Text: fmt.Sprintf("I'm reminding you about \"%s\".", info.Description),
				},
			},
		},
	)
	if err != nil {
		return core.ResponseEmpty, err
	}

	return core.ResponseEmpty, err

}
