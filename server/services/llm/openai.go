package llm

import (
	"bytes"
	"context"

	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type OpenAI struct {
	log    logrus.FieldLogger
	client *openai.Client
}

func NewOpenAI(log logrus.FieldLogger, key string) *OpenAI {
	var o OpenAI
	o.client = openai.NewClient(key)
	o.log = log
	return &o
}

func (o *OpenAI) Completion(ctx context.Context, prompt string) (string, error) {
	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Return only the JSON data requested in the prompt. Do not provide explanation",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", err
	}
	o.log.Debugf("OpenAI prompt: %s", prompt)
	o.log.Debugf("OpenAI response: %s", resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAI) Embedding(ctx context.Context, text string) ([]float32, error) {
	res, err := o.client.CreateEmbeddings(ctx, openai.EmbeddingRequestStrings{
		Model: openai.AdaEmbeddingV2,
		Input: []string{text},
	})
	if err != nil {
		return nil, err
	}
	return res.Data[0].Embedding, nil
}

func (o *OpenAI) SpeechToText(ctx context.Context, wavBytes []byte) (string, error) {
	trans, err := o.client.CreateTranscription(ctx, openai.AudioRequest{
		Model:    "whisper-1",
		FilePath: "voice.wav",
		Reader:   bytes.NewBuffer(wavBytes),
	})
	if err != nil {
		return "", err
	}

	return trans.Text, nil
}
