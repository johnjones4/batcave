package types

import (
	"github.com/jdkato/prose/v2"
	"github.com/olebedev/when"
	"github.com/sbl/ner"
)

type ResponseMessage struct {
	Text  string      `json:"text"`
	URL   string      `json:"url"`
	Extra interface{} `json:"extra"`
}

type AlertQueue interface {
	Enqueue(m ResponseMessage)
}

type RequestMessage struct {
	Message string `json:"message"`
}

type ParsedRequestMessage struct {
	Original      RequestMessage
	NamedEntities []ner.Entity
	Tokens        []prose.Token
	DateInfo      *when.Result
	IntentLabel   string
}
