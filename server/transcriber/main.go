package main

import (
	"log"
	"os"

	"github.com/johnjones4/hal-9000/server/hal9000/learning"
)

func main() {
	t, err := learning.NewVoiceTranscriber(learning.VoiceTranscriberConfiguration{
		ModelPath: os.Getenv("TRANSCRIBER_MODEL_PATH"),
	})
	if err != nil {
		panic(err)
	}
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	text, err := t.Transcribe(file)
	if err != nil {
		panic(err)
	}
	log.Println(text)
}
