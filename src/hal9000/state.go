package hal9000

import (
	"encoding/json"
	"fmt"
	"hal9000/util"
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

type State interface {
	Name() string
	ProcessIncomingMessage(c Person, m RequestMessage) (State, ResponseMessage, error)
}

func InitStateByName(name string) State {
	if name == util.StateTypeDefault {
		return DefaultState{}
	}
	return DefaultState{}
}

var subtypeMap map[uint8]string
var model *text.NaiveBayes
var nerExtractor *ner.Extractor
var dateExtractor *when.Parser

func LoadSubtypeMap() error {
	data, err := ioutil.ReadFile(os.Getenv("SUBTYPE_MAP_FILE"))
	if err != nil {
		return err
	}
	var preMap map[string]string
	err = json.Unmarshal(data, &preMap)
	if err != nil {
		return err
	}
	outMap := make(map[uint8]string)
	for key, val := range preMap {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return err
		}
		outMap[uint8(keyInt)] = val
	}
	subtypeMap = outMap

	return nil
}

func InitializeDefaultIncomingMessageParser() error {
	err := LoadSubtypeMap()
	if err != nil {
		return err
	}

	_model := text.NewNaiveBayes(nil, uint8(len(subtypeMap)), base.OnlyWordsAndNumbers)
	err = _model.RestoreFromFile(os.Getenv("MODEL_FILE"))
	if err != nil {
		return err
	}

	model = _model

	ext, err := ner.NewExtractor(os.Getenv("NER_MODEL_PATH"))
	if err != nil {
		return err
	}
	nerExtractor = ext

	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)
	dateExtractor = w

	return nil
}

type DefaultState struct{}

func (s DefaultState) Name() string { return util.StateTypeDefault }

func (s DefaultState) ProcessIncomingMessage(caller Person, input RequestMessage) (State, ResponseMessage, error) {
	class := model.Predict(input.Message)

	intentLabel, ok := subtypeMap[class]
	if !ok {
		return nil, ResponseMessage{}, fmt.Errorf("no alias for intent %d", class)
	}
	fmt.Println(intentLabel)

	nerTokens := ner.Tokenize(input.Message)
	es, err := nerExtractor.Extract(nerTokens)
	if err != nil {
		return nil, ResponseMessage{}, err
	}

	doc, err := prose.NewDocument(input.Message)
	if err != nil {
		return nil, ResponseMessage{}, err
	}
	tokens := doc.Tokens()

	dateInfo, err := dateExtractor.Parse(input.Message, time.Now())
	if err != nil {
		return nil, ResponseMessage{}, err
	}

	inputMessage := ParsedRequestMessage{
		Original:      input,
		NamedEntities: es,
		Tokens:        tokens,
		DateInfo:      dateInfo,
	}

	intent, err := GetIntentForIncomingMessage(intentLabel, caller, inputMessage)
	if err != nil {
		return nil, ResponseMessage{}, err
	}

	return intent.Execute(s)
}
