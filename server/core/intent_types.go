package core

import "context"

type IntentMetadata struct {
	IntentParseCompletion string
	IntentParseReceiver   map[string]any
}

type IntentMatcher interface {
	Match(ctx context.Context, req *Request) (IntentActor, IntentMetadata, error)
}

type IntentActor interface {
	IntentLabel() string
	IntentParsePrompt(req *Request) string
	ActOnIntent(ctx context.Context, req *Request, md *IntentMetadata) (Response, error)
}

type PushIntentActor interface {
	IntentActor
	ActOnAsyncIntent(ctx context.Context, source, clientId string, md *IntentMetadata) (PushMessage, error)
}

type PushIntentFactory interface {
	PushIntent(named string) PushIntentActor
}
