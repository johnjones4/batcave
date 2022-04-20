package runtime

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
)

var (
	ErrorInputEmpty = errors.New("empty input")
	ErrorNoIntent   = errors.New("no intent found")
)

const (
	ParseMetadataBodyText       = "text"
	ParseMetadataBodyAudio      = "audio"
	ParseMetadataIntentExplicit = "explicit"
	ParseMetadataIntentInferred = "inferred"
)

func (r *Runtime) Parse(in core.InboundBody, client core.Client, state string) (core.Inbound, error) {
	if len(in.Body) == 0 && len(in.Audio) == 0 {
		return core.Inbound{}, ErrorInputEmpty
	}

	request := core.Inbound{
		InboundBody: in,
		Client:      client,
		State:       state,
	}

	if len(request.Body) == 0 && len(request.Audio) != 0 {
		gzipedReader := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(request.Audio)))
		audioReader, err := gzip.NewReader(gzipedReader)
		if err != nil {
			return core.Inbound{}, nil
		}
		audio, err := ioutil.ReadAll(audioReader)
		if err != nil {
			return core.Inbound{}, nil
		}
		body, err := r.VoiceTranscriber.Transcribe(audio)
		if err != nil {
			return core.Inbound{}, nil
		}
		request.Body = body
		request.Audio = "" //TODO save this somewhere?
		request.ParseMetadata.Body = ParseMetadataBodyAudio
	} else {
		request.ParseMetadata.Body = ParseMetadataBodyText
	}

	if request.Body[0] == '/' {
		request.ParseMetadata.Intent = ParseMetadataIntentExplicit
		firstSpace := strings.Index(request.Body, " ")
		if firstSpace < 0 {
			request.Command = request.Body[1:]
			request.Body = ""
		} else {
			request.Command = strings.TrimSpace(request.Body[1:firstSpace])
			request.Body = strings.TrimSpace(request.Body[firstSpace:])
		}
	} else {
		request.ParseMetadata.Intent = ParseMetadataIntentInferred
		var err error
		request.Command, err = r.IntentPredictor.PredictIntent(request.Body)
		if err != nil {
			return core.Inbound{}, err
		}
	}

	if request.Command == "" {
		return core.Inbound{}, ErrorNoIntent
	}

	if len(client.Users) == 1 {
		userRec, err := r.UserStore.GetUser(client.Users[0])
		if err != nil {
			return core.Inbound{}, err
		}
		request.User = userRec.User
	}

	return request, nil
}
