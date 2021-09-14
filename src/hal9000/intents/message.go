package intents

import (
	"fmt"
	"hal9000/service"
	"hal9000/types"
	"hal9000/util"
)

type messageIntent struct {
	caller  messageSender
	person  types.Person
	message types.ResponseMessage
}

type messageSender interface {
	GetOriginName() string
}

func NewMessageIntent(runtime types.Runtime, c messageSender, m types.ParsedRequestMessage) (messageIntent, error) {
	person, messageStart, err := getPersonInParsedRequestMessage(runtime, m)
	if err != nil {
		return messageIntent{}, err
	}

	sendMessage := types.ResponseMessage{
		Text:  fmt.Sprintf("Message from %s: \"%s\"", c.GetOriginName(), util.ConcatTokensInRange(m.Tokens, messageStart, len(m.Tokens))),
		URL:   "",
		Extra: nil,
	}

	return messageIntent{c, person, sendMessage}, nil
}

func (i messageIntent) Execute(runtime types.Runtime, lastState types.State) (types.State, types.ResponseMessage, error) {
	err := runtime.People().SendMessageToPerson(i.person, i.message)
	if err != nil {
		return nil, types.ResponseMessage{}, err
	}

	return lastState, util.MessageOk(), nil
}

func getPersonInParsedRequestMessage(runtime types.Runtime, m types.ParsedRequestMessage) (types.Person, int, error) {
	for _, entity := range m.NamedEntities {
		person, err := runtime.People().GetPersonByName(entity.Name)
		if err != nil && err != service.ErrorPersonNotFound {
			return nil, 0, err
		} else if err != service.ErrorPersonNotFound {
			return person, entity.Range.End, nil
		}
	}

	nouns := util.GetContiguousUniformTokens(m.Tokens, []string{"NN", "NNP", "NNPS", "NNS"})
	for _, nounSet := range nouns {
		for i := 0; i < len(nounSet.Tokens); i++ {
			for j := len(nounSet.Tokens); j >= i+1; j-- {
				nounStr := util.ConcatTokensInRange(nounSet.Tokens, i, j)
				person, err := runtime.People().GetPersonByName(nounStr)
				if err != nil && err != service.ErrorPersonNotFound {
					return nil, 0, err
				} else if err != service.ErrorPersonNotFound {
					return person, nounSet.End, nil
				}
			}
		}
	}

	return nil, 0, service.ErrorPersonNotFound
}
