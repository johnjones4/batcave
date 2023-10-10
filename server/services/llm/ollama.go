package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type OllamaCfg struct {
	URL   string `json:"url"`
	Model string `json:"model"`
}

type Ollama struct {
	OllamaCfg
}

func NewOllama(cfg OllamaCfg) *Ollama {
	return &Ollama{
		OllamaCfg: cfg,
	}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaGenerateResponseLine struct {
	Response string `json:"response,omitempty"`
}

type ollamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func (o *Ollama) Completion(ctx context.Context, prompt string) (string, error) {
	rerBytes, err := json.Marshal(ollamaRequest{
		Model:  o.Model,
		Prompt: prompt,
	})
	if err != nil {
		return "", err
	}
	res, err := http.Post(o.URL+"/api/generate", "application/json", bytes.NewBuffer(rerBytes))
	if err != nil {
		return "", err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var completion strings.Builder

	reslines := strings.Split(string(resBytes), "\n")
	for _, line := range reslines {
		if line == "" {
			continue
		}
		var lineStr ollamaGenerateResponseLine
		err = json.Unmarshal([]byte(line), &lineStr)
		if err != nil {
			return "", err
		}
		if lineStr.Response == "" {
			continue
		}
		completion.WriteString(lineStr.Response)
	}

	return completion.String(), nil
}

func (o *Ollama) Embedding(ctx context.Context, text string) ([]float32, error) {
	rerBytes, err := json.Marshal(ollamaRequest{
		Model:  o.Model,
		Prompt: text,
	})
	if err != nil {
		return nil, err
	}
	res, err := http.Post(o.URL+"/api/embeddings", "application/json", bytes.NewBuffer(rerBytes))
	if err != nil {
		return nil, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resSt ollamaEmbeddingResponse
	err = json.Unmarshal(resBytes, &resSt)
	if err != nil {
		return nil, err
	}

	return resSt.Embedding, nil
}
