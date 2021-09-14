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

func (s DefaultState) ProcessIncomingMessage(caller types.Person, input types.RequestMessage) (types.State, types.ResponseMessage, error) {

	intent, err := GetIntentForIncomingMessage(intentLabel, caller, inputMessage)
	if err != nil {
		return nil, types.ResponseMessage{}, err
	}

	return intent.Execute(s)
}
