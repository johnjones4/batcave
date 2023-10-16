package processors

import (
	"context"
	"encoding/base64"
	"main/core"
)

func (p *Processors) TexttToSpeech(ctx context.Context, req *core.Request, res *core.Response) error {
	if res.Message.Text != "" && res.Message.Audio.Data == "" {
		wavBytes, err := p.TTS.TextToSpeech(ctx, res.Message.Text)
		if err != nil {
			return err
		}
		res.Message.Audio.Data = base64.StdEncoding.EncodeToString(wavBytes)
	}
	return nil
}
