package learning

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
)

type trainingItem struct {
	id   uint8
	data string
}

func Train(trainingDataFilePath, mapOutputFile, modelOutputFile string) error {
	trainingDataFile, err := os.Open(trainingDataFilePath)
	if err != nil {
		return err
	}

	r := csv.NewReader(trainingDataFile)

	intents := make([]string, 0)
	reverseMap := make(map[string]uint8)
	trainingData := make([]trainingItem, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		intent := record[0]
		id, ok := reverseMap[intent]
		if !ok {
			id = uint8(len(intents))
			reverseMap[intent] = id
			intents = append(intents, intent)
		}

		trainingData = append(trainingData, trainingItem{
			id:   id,
			data: record[1],
		})
	}

	intentMapData, err := json.Marshal(intents)
	if err != nil {
		return err
	}

	err = os.WriteFile(mapOutputFile, intentMapData, 0666)
	if err != nil {
		return err
	}

	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, uint8(len(intents)), base.OnlyWordsAndNumbers)
	go model.OnlineLearn(errors)

	for _, row := range trainingData {
		stream <- base.TextDatapoint{
			X: row.data,
			Y: row.id,
		}
	}

	close(stream)

	for {
		err := <-errors
		if err != nil {
			return err
		} else {
			break
		}
	}

	err = model.PersistToFile(modelOutputFile)
	if err != nil {
		return err
	}

	return nil
}
