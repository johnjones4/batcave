package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
	"github.com/joho/godotenv"
)

type TrainingDataRow struct {
	Subtype uint8
	String  string
}

func saveSubtypeMap(subtypeMap map[string]uint8) error {
	outMap := make(map[string]string)
	for label, intkey := range subtypeMap {
		outMap[fmt.Sprint(intkey)] = label
	}
	data, err := json.Marshal(outMap)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(os.Getenv("SUBTYPE_MAP_FILE"), data, 0777)
}

func loadTrainingData() ([]TrainingDataRow, map[string]uint8, error) {
	csvfile, err := os.Open(os.Getenv("TRAINING_DATA_FILE"))
	if err != nil {
		return nil, nil, err
	}
	r := csv.NewReader(csvfile)
	data := make([]TrainingDataRow, 0)
	subtypeMap := make(map[string]uint8)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if len(record) != 2 {
			return nil, nil, fmt.Errorf("bad row size: %d", len(record))
		}
		subtype, ok := subtypeMap[record[0]]
		if !ok {
			subtype = uint8(len(subtypeMap))
			subtypeMap[record[0]] = subtype
		}
		data = append(data, TrainingDataRow{subtype, record[1]})
	}
	return data, subtypeMap, nil
}

func train() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	data, subtypeMap, err := loadTrainingData()
	if err != nil {
		return err
	}

	err = saveSubtypeMap(subtypeMap)
	if err != nil {
		return err
	}

	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, uint8(len(subtypeMap)), base.OnlyWordsAndNumbers)
	go model.OnlineLearn(errors)

	for _, row := range data {
		stream <- base.TextDatapoint{
			X: row.String,
			Y: row.Subtype,
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

	err = model.PersistToFile(os.Getenv("MODEL_FILE"))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := train()
	if err != nil {
		fmt.Println(err)
	}
}
