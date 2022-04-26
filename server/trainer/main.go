package main

import (
	"os"

	"github.com/johnjones4/hal-9000/server/hal9000/learning"
)

func main() {
	trainingFile := os.Getenv("TRAINING_DATA_FILE")
	mapFile := os.Getenv("INTENT_MAP_FILE")
	modelFile := os.Getenv("MODEL_FILE")

	err := learning.TrainIntent(trainingFile, mapFile, modelFile)
	if err != nil {
		panic(err)
	}
}
