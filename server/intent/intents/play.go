package intents

import (
	"context"
	"fmt"
	"main/core"
	"main/services/tunein"

	"github.com/mitchellh/mapstructure"
)

type Play struct {
	TuneIn *tunein.TuneIn
}

type playIntentParseReceiver struct {
	AudioSource string `json:"audioSource"`
}

func (p *Play) IntentLabel() string {
	return "play"
}

func (p *Play) IntentParsePrompt(req *core.Request) string {
	return fmt.Sprintf("Extract the name of the audio source from the phrase \"%s\" in the JSON format: {\"audioSource\":\"\"}", req.Message.Text)
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

	url, err := p.TuneIn.GetStreamURL(info.AudioSource)
	if err != nil {
		return core.Response{}, err
	}

	return core.Response{
		Media: core.Media{
			URL:  url,
			Type: core.MediaTypeAudioStream,
		},
		Action: core.ActionPlay,
	}, nil
}
