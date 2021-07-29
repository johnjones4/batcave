package hal9000

import (
	"github.com/jdkato/prose/v2"
	"github.com/olebedev/when"
	"github.com/sbl/ner"
)

type Message struct {
	Text  string      `json:"text"`
	URL   string      `json:"url"`
	Extra interface{} `json:"extra"`
}

type ParsedMessage struct {
	Original      string
	NamedEntities []ner.Entity
	Tokens        []prose.Token
	DateInfo      *when.Result
}

func MessageOk() Message {
	return Message{"Ok", "", nil}
}
