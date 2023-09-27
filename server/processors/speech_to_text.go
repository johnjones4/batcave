package processors

import (
	"context"
	"encoding/base64"
	"main/core"
)

func (p *Processors) SpeechToText(ctx context.Context, req *core.Request) error {
	if req.Message.Text == "" && req.Message.Audio.Data != "" {
		b, err := base64.StdEncoding.DecodeString(req.Message.Audio.Data)
		if err != nil {
			return err
		}

		txt, err := p.STT.SpeechToText(ctx, b)
		if err != nil {
			return err
		}

		req.Message.Text = txt
	}
	return nil
}
