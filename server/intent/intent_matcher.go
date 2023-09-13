package intent

import (
	"context"
	"encoding/json"
	"errors"
	"main/core"
	"main/util"

	"github.com/sirupsen/logrus"
)

var ErrorNoActorForIntent = errors.New("intent not found")
var ErrorUnexpectedIntentInput = errors.New("unexpected input")

type IntentMatcher struct {
	intents []core.IntentActor
	llm     core.LLM
	istore  core.IntentEmbeddingStore
	log     logrus.FieldLogger
}

func NewIntentMatcher(log logrus.FieldLogger, intents []core.IntentActor, llm core.LLM, istore core.IntentEmbeddingStore) (*IntentMatcher, error) {
	p := IntentMatcher{
		intents: intents,
		llm:     llm,
		istore:  istore,
		log:     log,
	}

	return &p, nil
}

func (id *IntentMatcher) Match(ctx context.Context, req *core.Request) (core.IntentActor, core.IntentMetadata, error) {
	emedding, err := id.llm.Embedding(ctx, req.Message.Text)
	if err != nil {
		return nil, core.IntentMetadata{}, err
	}

	intentLabel, err := id.istore.GetClosestMatchingIntent(ctx, emedding)
	if err != nil {
		return nil, core.IntentMetadata{}, err
	}

	if intentLabel == "" {
		intentLabel = "unknown"
	}

	id.log.Debugf("Incoming intent: %s", intentLabel)

	for _, intentActor := range id.intents {
		if intentLabel == intentActor.IntentLabel() {
			var md core.IntentMetadata
			prompt := intentActor.IntentParsePrompt(req)
			if prompt != "" {
				intentParseResponse, err := id.llm.Completion(ctx, prompt)
				if err != nil {
					return nil, core.IntentMetadata{}, err
				}

				md.IntentParseCompletion = util.CleanLLMJSON(intentParseResponse)

				var receiver map[string]any
				err = json.Unmarshal([]byte(md.IntentParseCompletion), &receiver)
				if err != nil {
					return nil, core.IntentMetadata{}, err
				}
				md.IntentParseReceiver = receiver
			}
			return intentActor, md, nil
		}
	}

	return nil, core.IntentMetadata{}, ErrorNoActorForIntent
}
