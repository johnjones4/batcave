package intents

import (
	"context"
	"fmt"
	"main/core"

	"github.com/mitchellh/mapstructure"
)

type Play struct {
	TuneIn core.TuneIn
	Push   core.RecurringPush
}

type playIntentParseReceiver struct {
	AudioSource           string `json:"audioSource"`
	RecurranceDescription string `json:"recurranceDescription"`
}

func (p *Play) IntentLabel() string {
	return "play"
}

func (p *Play) IntentParsePrompt(req *core.Request) string {
	return fmt.Sprintf("Extract the name of the audio source and recurrance description in crontab format or blank from the phrase \"%s\" in the JSON format: {\"audioSource\":\"\",\"recurranceDescription\":\"\"}", req.Message.Text)
}

func (p *Play) IntentParseReceiver() any {
	return playIntentParseReceiver{}
}

func (p *Play) ActOnIntent(ctx context.Context, req *core.Request, md *core.IntentMetadata) (core.Response, error) {
	var info playIntentParseReceiver
	err := mapstructure.Decode(md.IntentParseReceiver, &info)
	if err != nil {
		return core.ResponseEmpty, err
	}

	if info.RecurranceDescription != "" {
		err = p.Push.SendRecurring(ctx, req.Source, req.ClientID, info.RecurranceDescription, p.IntentLabel(), map[string]any{
			"audioSource": info.AudioSource,
		})
		if err != nil {
			return core.ResponseEmpty, err
		}
		return core.ResponseEmpty, nil
	}

	url, err := p.TuneIn.GetStreamURL(info.AudioSource)
	if err != nil {
		return core.Response{}, err
	}

	return core.Response{
		OutboundMessage: core.OutboundMessage{
			Media: core.Media{
				URL:  url,
				Type: core.MediaTypeAudioStream,
			},
		},
		Action: core.ActionPlay,
	}, nil
}
