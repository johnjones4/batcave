package main

import "fmt"

func HandleInput(session *Session, message string) (Response, error) {
	parsedInputMessage, err := computeParsedInputMessage(message)
	if err != nil {
		return Response{}, err
	}
	intent, err := session.State.InferIntent(parsedInputMessage)
	if err != nil {
		return Response{}, err
	}
	nextState, response, err := intent.Process(parsedInputMessage)
	if err != nil {
		return Response{}, err
	}
	if !session.State.CanTransitionToState(nextState) {
		return Response{}, fmt.Errorf("cannot transition from %s to %s", session.State.Name(), nextState.Name())
	}
	session.State = nextState
	return response, nil
}

func computeParsedInputMessage(message string) (ParsedInputMessage, error) {
	return ParsedInputMessage{}, nil //TODO
}
