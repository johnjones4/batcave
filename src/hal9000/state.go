package hal9000

import (
	"hal9000/types"
	"hal9000/util"
)

func initStateByName(name string) types.State {
	if name == util.StateTypeDefault {
		return defaultState{}
	}
	return defaultState{}
}

type defaultState struct{}

func (s defaultState) Name() string { return util.StateTypeDefault }

func (s defaultState) ProcessIncomingMessage(r types.Runtime, caller types.Person, input types.RequestMessage) (types.State, types.ResponseMessage, error) {
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
