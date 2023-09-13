package intent

import (
	"context"
	"encoding/json"
	"main/core"
	"os"

	"github.com/sirupsen/logrus"
)

func ParseIntents(ctx context.Context, log logrus.FieldLogger, llm core.LLM, istore core.IntentEmbeddingStore, intentsFile string) error {
	log.Debugf("Loading intents from %s", intentsFile)
	contents, err := os.ReadFile(intentsFile)
	if err != nil {
		return err
	}

	var intents map[string]string
	err = json.Unmarshal(contents, &intents)
	if err != nil {
		return err
	}

	for intent, desc := range intents {
		log.Debugf("Processing intent %s: \"%s\"", intent, desc)
		embedding, err := llm.Embedding(ctx, desc)
		if err != nil {
			return err
		}
		err = istore.UpdateIntentEmbedding(ctx, intent, embedding)
		if err != nil {
			return err
		}
	}

	return nil
}
