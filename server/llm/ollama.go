package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Ollama struct {
	url string
	log logrus.FieldLogger
}

type ollamaCompletionRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
}

type ollamaCompletionResponseSegment struct {
	Response string `json:"response,omitempty"`
}

type ollamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func NewOllama(log logrus.FieldLogger, url string) *Ollama {
	return &Ollama{url, log}
}

func (o *Ollama) Completion(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(ollamaCompletionRequest{
		Model:  "llama2",
		Prompt: prompt,
		System: "Return only the JSON data requested in the prompt. Do not provide explanation",
	})
	if err != nil {
		return "", err
	}

	res, err := http.Post(o.url+"/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	parts := strings.Split(string(resbody), "\n")
	var response strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		var seg ollamaCompletionResponseSegment
		err = json.Unmarshal([]byte(part), &seg)
		if err != nil {
			return "", err
		}
		if seg.Response != "" {
			response.WriteString(seg.Response)
		}
	}

	reponseFull := response.String()

	o.log.Debugf("Ollama prompt: %s", prompt)
	o.log.Debugf("Ollama response: %s", reponseFull)

	return reponseFull, nil
}

func (o *Ollama) Embedding(ctx context.Context, text string) ([]float32, error) {
	body, err := json.Marshal(ollamaEmbeddingRequest{
		Model:  "llama2",
		Prompt: text,
	})
	if err != nil {
		return nil, err
	}

	res, err := http.Post(o.url+"/api/embeddings", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var embRes ollamaEmbeddingResponse
	err = json.Unmarshal(resbody, &embRes)
	if err != nil {
		return nil, err
	}

	return embRes.Embedding, nil
}
