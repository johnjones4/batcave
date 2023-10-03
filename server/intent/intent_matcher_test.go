package intent

import (
	"context"
	"encoding/json"
	"errors"
	"main/core"
	"main/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type intentMatcherTestCase struct {
	request           core.Request
	embedding         []float32
	embeddingError    error
	intentLabel       string
	intentLabelError  error
	intentParseString string
	metadata          map[string]any
	completionError   error
}

var errorTestError = errors.New("test")

func TestIntentMatcher(t *testing.T) {
	ctrl := gomock.NewController(t)

	cases := []intentMatcherTestCase{
		{
			request: core.Request{
				Message: core.Message{
					Text: "hello",
				},
			},
			embedding:         []float32{1},
			intentLabel:       "label",
			intentParseString: "test",
			metadata: map[string]any{
				"k": "v",
			},
		},
		{
			request: core.Request{
				Message: core.Message{
					Text: "hello",
				},
			},
			embeddingError: errorTestError,
		},
		{
			request: core.Request{
				Message: core.Message{
					Text: "hello",
				},
			},
			embedding:        []float32{1},
			intentLabelError: errorTestError,
		},
		{
			request: core.Request{
				Message: core.Message{
					Text: "hello",
				},
			},
			embedding:         []float32{1},
			intentLabel:       "label",
			intentParseString: "test",
			completionError:   errorTestError,
		},
	}

	log := mocks.NewMockFieldLogger(ctrl)
	llm := mocks.NewMockLLM(ctrl)
	iStore := mocks.NewMockIntentEmbeddingStore(ctrl)
	for _, c := range cases {
		mockIntentActor := mocks.NewMockIntentActor(ctrl)

		matcher, err := NewIntentMatcher(log, []core.IntentActor{mockIntentActor}, llm, iStore)
		assert.Nil(t, err)

		llm.EXPECT().Embedding(gomock.Any(), c.request.Message.Text).Return(c.embedding, c.embeddingError)
		if c.embeddingError == nil {
			iStore.EXPECT().ClosestMatchingIntent(gomock.Any(), c.embedding).Return(c.intentLabel, c.intentLabelError)

			if c.intentLabelError == nil {
				mockIntentActor.EXPECT().IntentLabel().Return(c.intentLabel)
				mockIntentActor.EXPECT().IntentParsePrompt(&c.request).Return(c.intentParseString)
				log.EXPECT().Debugf(gomock.Any(), c.intentLabel)

				var mdJson string
				if c.metadata != nil {
					mdJsonB, _ := json.Marshal(c.metadata)
					mdJson = string(mdJsonB)
				}
				llm.EXPECT().Completion(gomock.Any(), c.intentParseString).Return(mdJson, c.completionError)
			}
		}

		intentActor, metadata, err := matcher.Match(context.Background(), &c.request)

		if c.embeddingError != nil {
			assert.Equal(t, c.embeddingError, err)
			assert.Nil(t, intentActor)
		} else if c.intentLabelError != nil {
			assert.Equal(t, c.intentLabelError, err)
			assert.Nil(t, intentActor)
		} else if c.completionError != nil {
			assert.Equal(t, c.completionError, err)
			assert.Nil(t, intentActor)
		} else {
			assert.Equal(t, mockIntentActor, intentActor)
			assert.Equal(t, c.metadata, metadata.IntentParseReceiver)
		}
	}
}
