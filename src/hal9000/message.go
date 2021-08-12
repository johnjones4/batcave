package hal9000

import (
	"fmt"

	"github.com/jdkato/prose/v2"
	"github.com/olebedev/when"
	"github.com/sbl/ner"
)

type ResponseMessage struct {
	Text  string      `json:"text"`
	URL   string      `json:"url"`
	Extra interface{} `json:"extra"`
}

type RequestMessage struct {
	Message string `json:"message"`
}

type ParsedRequestMessage struct {
	Original      RequestMessage
	NamedEntities []ner.Entity
	Tokens        []prose.Token
	DateInfo      *when.Result
}

func MessageOk() ResponseMessage {
	return ResponseMessage{"Ok", "", nil}
}

func MessageError(err error) ResponseMessage {
	return ResponseMessage{fmt.Sprintf("Encoutered error: %s", err), "", nil}
}
