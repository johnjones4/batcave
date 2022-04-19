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

type PredictorConfiguration struct {
	IntentMapFilePath string
	ModelFilePath     string
}

type Predictor struct {
	intents []string
	model   *text.NaiveBayes
}

func NewPredictor(configuration PredictorConfiguration) (*Predictor, error) {
	p := Predictor{}

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

func (p *Predictor) PredictIntent(message string) (string, error) {
	class := p.model.Predict(message)
	if class >= uint8(len(p.intents)) {
		return "", ErrorUnknownClassification
	}
	className := p.intents[class]
	return className, nil
}
