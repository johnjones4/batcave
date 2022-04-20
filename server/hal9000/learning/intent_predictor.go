package learning

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
)

var (
	ErrorUnknownClassification = errors.New("unknown classification")
)

type IntentPredictorConfiguration struct {
	IntentMapFilePath string
	ModelFilePath     string
}

type IntentPredictor struct {
	intents []string
	model   *text.NaiveBayes
}

func NewIntentPredictor(configuration IntentPredictorConfiguration) (*IntentPredictor, error) {
	p := IntentPredictor{}

	intentMapBytes, err := os.ReadFile(configuration.IntentMapFilePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(intentMapBytes, &p.intents)
	if err != nil {
		return nil, err
	}

	p.model = text.NewNaiveBayes(nil, uint8(len(p.intents)), base.OnlyWordsAndNumbers)
	err = p.model.RestoreFromFile(configuration.ModelFilePath)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *IntentPredictor) PredictIntent(message string) (string, error) {
	class := p.model.Predict(message)
	if class >= uint8(len(p.intents)) {
		return "", ErrorUnknownClassification
	}
	className := p.intents[class]
	return className, nil
}
