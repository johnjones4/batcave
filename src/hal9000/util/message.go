package util

import (
	"fmt"
	"hal9000/types"
)

func MessageOk() types.ResponseMessage {
	return types.ResponseMessage{"Ok", "", nil}
}

func MessageError(err error) types.ResponseMessage {
	return types.ResponseMessage{fmt.Sprintf("Encoutered error: %s", err), "", nil}
}
