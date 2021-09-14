package service

import (
	"encoding/json"
	"hal9000/types"
	"hal9000/util"
	"io/ioutil"
	"os"
	"strings"
)

type DisplayableConcrete struct {
	Names  []string `json:"names"`
	URL    string   `json:"url"`
	Type   string   `json:"type"`
	Source string   `json:"source"`
}

func (d DisplayableConcrete) GetNames() []string {
	return d.Names
}

func (d DisplayableConcrete) GetURL() string {
	return d.URL
}

func (d DisplayableConcrete) GetType() string {
	return d.Type
}

func (d DisplayableConcrete) GetSource() string {
	return d.Source
}

type displayablesProviderConcrete struct {
	displayables []types.Displayable
}

func InitDisplayablesProvider() (types.DisplayablesProvider, error) {
	bytes, err := ioutil.ReadFile(os.Getenv("DISPLAYABLES_MANIFEST_PATH"))
	if err != nil {
		return nil, err
	}
	var displayablesConcrete []DisplayableConcrete
	err = json.Unmarshal(bytes, &displayablesConcrete)
	if err != nil {
		return nil, err
	}
	displayables := make([]types.Displayable, len(displayablesConcrete))
	for i, d := range displayablesConcrete {
		displayables[i] = d
	}
	return displayablesProviderConcrete{displayables}, nil
}

func (dp displayablesProviderConcrete) FindDisplayableInString(str string) (types.Displayable, error) {
	lcStr := strings.ToLower(str)
	for _, displayable := range dp.displayables {
		for _, name := range displayable.GetNames() {
			lcName := strings.ToLower(name)
			if strings.Contains(lcStr, lcName) {
				return displayable, nil
			}
		}
	}
	return nil, util.ErrorDisplayNotFound
}
