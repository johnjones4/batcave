package intents

import (
	"context"
	"fmt"
	"main/core"
	"main/util"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Remind struct {
	Push core.Push
}

type scheduleIntentParseReceiver struct {
	Date                  string `json:"date,omitempty"`
	Description           string `json:"description,omitempty"`
	RecurranceDescription string `json:"recurranceDescription"`
}

func (i *Remind) IntentLabel() string {
	return "remind"
}

func (i *Remind) IntentParsePrompt(req *core.Request) string {
	return fmt.Sprintf("Extract the exact date and time relative to %s or recurrance description in crontab format along with a description from the phrase \"%s\" and return the information in the JSON format {\"date\":\"RFC3339 format\",\"description\":\"\",\"recurranceDescription\":\"\"}", time.Now().String(), req.Message.Text)
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

	if info.RecurranceDescription != "" {
		err = i.Push.SendRecurring(ctx, req.Source, req.ClientID, info.RecurranceDescription, i.IntentLabel(), map[string]any{
			"description": info.Description,
		})
		if err != nil {
			return core.ResponseEmpty, err
		}
		return core.ResponseEmpty, nil
	} else {
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
					Message: core.Message{
						Text: fmt.Sprintf("I'm reminding you about \"%s\".", info.Description),
					},
				},
			},
		)
		if err != nil {
			return core.ResponseEmpty, err
		}

		return core.ResponseEmpty, nil
	}
}

func (i *Remind) ActOnAsyncIntent(ctx context.Context, source, clientId string, md *core.IntentMetadata) (core.PushMessage, error) {
	var info scheduleIntentParseReceiver
	err := mapstructure.Decode(md.IntentParseReceiver, &info)
	if err != nil {
		return core.PushMessage{}, err
	}

	return core.PushMessage{
		OutboundMessage: core.OutboundMessage{
			Message: core.Message{
				Text: fmt.Sprintf("I'm reminding you about \"%s\".", info.Description),
			},
		},
	}, nil
}
