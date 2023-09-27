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

type LlamaDotCpp struct {
	url string
	log logrus.FieldLogger
}

type llamaDotCppCompletionRequest struct {
	Prompt string `json:"prompt"`
	System string `json:"system"`
}

type llamaDotCppCompletionResponseSegment struct {
	Response string `json:"response,omitempty"`
}

type llamaDotCppEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type llamaDotCppEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func NewLlamaDotCpp(log logrus.FieldLogger, url string) *LlamaDotCpp {
	return &LlamaDotCpp{url, log}
}

func (o *LlamaDotCpp) Completion(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(llamaDotCppCompletionRequest{
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
		var seg llamaDotCppCompletionResponseSegment
		err = json.Unmarshal([]byte(part), &seg)
		if err != nil {
			return "", err
		}
		if seg.Response != "" {
			response.WriteString(seg.Response)
		}
	}

	reponseFull := response.String()

	o.log.Debugf("llamaDotCpp prompt: %s", prompt)
	o.log.Debugf("llamaDotCpp response: %s", reponseFull)

	return reponseFull, nil
}

func (o *LlamaDotCpp) Embedding(ctx context.Context, text string) ([]float32, error) {
	body, err := json.Marshal(llamaDotCppEmbeddingRequest{
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

	var embRes llamaDotCppEmbeddingResponse
	err = json.Unmarshal(resbody, &embRes)
	if err != nil {
		return nil, err
	}

	return embRes.Embedding, nil
}
