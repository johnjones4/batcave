package service

import (
	"encoding/json"
	"fmt"
	"hal9000/types"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
	"github.com/jdkato/prose/v2"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/sbl/ner"
)

type parserProviderConcrete struct {
	subtypeMap    map[uint8]string
	model         *text.NaiveBayes
	nerExtractor  *ner.Extractor
	dateExtractor *when.Parser
}

func InitParserProvider() (types.ParserProvider, error) {
	pp := parserProviderConcrete{}

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
	pp.subtypeMap = outMap

	_model := text.NewNaiveBayes(nil, uint8(len(pp.subtypeMap)), base.OnlyWordsAndNumbers)
	err = _model.RestoreFromFile(os.Getenv("MODEL_FILE"))
	if err != nil {
		return nil, err
	}

	pp.model = _model

	ext, err := ner.NewExtractor(os.Getenv("NER_MODEL_PATH"))
	if err != nil {
		return nil, err
	}
	pp.nerExtractor = ext

	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)
	pp.dateExtractor = w

	return pp, nil
}

func (pp parserProviderConcrete) ProcessMessage(input types.RequestMessage) (types.ParsedRequestMessage, error) {
	class := pp.model.Predict(input.Message)

	intentLabel, ok := pp.subtypeMap[class]
	if !ok {
		return types.ParsedRequestMessage{}, fmt.Errorf("no alias for intent %d", class)
	}

	nerTokens := ner.Tokenize(input.Message)
	es, err := pp.nerExtractor.Extract(nerTokens)
	if err != nil {
		return types.ParsedRequestMessage{}, err
	}

	doc, err := prose.NewDocument(input.Message)
	if err != nil {
		return types.ParsedRequestMessage{}, err
	}
	tokens := doc.Tokens()

	dateInfo, err := pp.dateExtractor.Parse(input.Message, time.Now())
	if err != nil {
		return types.ParsedRequestMessage{}, err
	}

	return types.ParsedRequestMessage{
		Original:      input,
		NamedEntities: es,
		Tokens:        tokens,
		DateInfo:      dateInfo,
		IntentLabel:   intentLabel,
	}, nil
}
