package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
	"github.com/joho/godotenv"
)

func loadSubtypeMap() (map[uint8]string, error) {
	data, err := ioutil.ReadFile(os.Getenv("SUBTYPE_MAP_FILE"))
	if err != nil {
		return nil, err
	}
	var preMap map[string]string
	err = json.Unmarshal(data, &preMap)
	if err != nil {
		return nil, err
	}
	outMap := make(map[uint8]string)
	for key, val := range preMap {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}
		outMap[uint8(keyInt)] = val
	}
	return outMap, nil
}

func predict(input string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	subtypeMap, err := loadSubtypeMap()
	if err != nil {
		return "", err
	}

	model := text.NewNaiveBayes(nil, uint8(len(subtypeMap)), base.OnlyWordsAndNumbers)
	err = model.RestoreFromFile(os.Getenv("MODEL_FILE"))
	if err != nil {
		return "", err
	}

	class := model.Predict(input)

	className, ok := subtypeMap[class]
	if !ok {
		return "", fmt.Errorf("no alias for class %d", class)
	}

	return className, nil
}

func main() {
	class, err := predict(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(class)
}
