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
	ProcessIncomingMessage(m string) (State, Message, error)
}

func InitStateByName(name string) (State, error) {
	if name == util.StateTypeDefault {
		return DefaultState{}, nil
	}
	return nil, fmt.Errorf("no state type named %s", name)
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

func (s DefaultState) ProcessIncomingMessage(input string) (State, Message, error) {
	class := model.Predict(input)

	intentLabel, ok := subtypeMap[class]
	if !ok {
		return nil, Message{}, fmt.Errorf("no alias for intent %d", class)
	}

	nerTokens := ner.Tokenize(input)
	es, err := nerExtractor.Extract(nerTokens)
	if err != nil {
		return nil, Message{}, err
	}

	doc, err := prose.NewDocument(input)
	if err != nil {
		return nil, Message{}, err
	}
	tokens := doc.Tokens()

	dateInfo, err := dateExtractor.Parse(input, time.Now())
	if err != nil {
		return nil, Message{}, err
	}

	inputMessage := ParsedMessage{
		Original:      input,
		NamedEntities: es,
		Tokens:        tokens,
		DateInfo:      dateInfo,
	}

	fmt.Println(inputMessage)

	intent, err := GetIntentForIncomingMessage(intentLabel, inputMessage)
	if err != nil {
		return nil, Message{}, err
	}

	return intent.Execute(s)
}
