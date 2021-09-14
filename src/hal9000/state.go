package hal9000

import (
	"hal9000/types"
	"hal9000/util"
)

func InitStateByName(name string) types.State {
	if name == util.StateTypeDefault {
		return DefaultState{}
	}
	return DefaultState{}
}

type DefaultState struct{}

func (s DefaultState) Name() string { return util.StateTypeDefault }

func (s DefaultState) ProcessIncomingMessage(r types.Runtime, caller types.Person, input types.RequestMessage) (types.State, types.ResponseMessage, error) {
	inputMessage, err := r.Parser().ProcessMessage(input)
	if err != nil {
		return nil, types.ResponseMessage{}, err
	}

	intent, err := GetIntentForIncomingMessage(r, caller, inputMessage)
	if err != nil {
		return nil, types.ResponseMessage{}, err
	}

	return intent.Execute(r, s)
}
