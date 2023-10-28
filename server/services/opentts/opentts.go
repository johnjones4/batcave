package opentts

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type OpenTTSConfiguration struct {
	URL   string `json:"url"`
	Voice string `json:"voice"`
}

type OpenTTS struct {
	Configuration OpenTTSConfiguration
}

func (tts *OpenTTS) TextToSpeech(ctx context.Context, text string) ([]byte, error) {
	params := url.Values{
		"voice": {tts.Configuration.Voice},
		"text":  {text},
	}
	res, err := http.Get(fmt.Sprintf("%s/api/tts?%s", tts.Configuration.URL, params.Encode()))
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
