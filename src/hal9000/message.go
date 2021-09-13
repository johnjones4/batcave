package hal9000

import (
	"fmt"
	"hal9000/util"

	"github.com/jdkato/prose/v2"
	"github.com/olebedev/when"
	"github.com/sbl/ner"
)

type RequestMessage struct {
	Message string `json:"message"`
}

type ParsedRequestMessage struct {
	Original      RequestMessage
	NamedEntities []ner.Entity
	Tokens        []prose.Token
	DateInfo      *when.Result
}

func MessageOk() util.ResponseMessage {
	return util.ResponseMessage{"Ok", "", nil}
}

func MessageError(err error) util.ResponseMessage {
	return util.ResponseMessage{fmt.Sprintf("Encoutered error: %s", err), "", nil}
}
