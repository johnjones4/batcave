package processors

import "main/core"

type Processors struct {
	LLM            core.LLM
	STT            core.STT
	TTS            core.TTS
	ClientRegistry core.ClientRegistry
}
