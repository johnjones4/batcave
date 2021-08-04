package util

import (
	"sort"

	"github.com/jdkato/prose/v2"
)

type ContigiousUniformTokenSet struct {
	Tokens []prose.Token
	Start  int
	End    int
}

func ContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func ConcatTokensInRange(tokens []prose.Token, start int, end int) string {
	outputStr := ""
	for i := start; i < end; i++ {
		if i != start && tokens[i].Text[0] != '\'' {
			outputStr += " "
		}
		outputStr += tokens[i].Text
	}
	return outputStr
}

func GetContiguousUniformTokens(tokens []prose.Token, tags []string) []ContigiousUniformTokenSet {
	output := make([]ContigiousUniformTokenSet, 0)
	var currentSet *ContigiousUniformTokenSet
	for i, token := range tokens {
		if ContainsString(tags, token.Tag) {
			if currentSet != nil {
				currentSet.Tokens = append(currentSet.Tokens, token)
				currentSet.End = i + 1
			} else {
				currentSet = &ContigiousUniformTokenSet{
					Tokens: []prose.Token{token},
					Start:  i,
					End:    i + 1,
				}
			}
		} else {
			if currentSet != nil {
				output = append(output, *currentSet)
			}
			currentSet = nil
		}
	}
	return output
}

type Nameable interface {
	GetNames() []string
}

type NameableSequenceItem struct {
	Name     string
	Nameable Nameable
}

func GenerateNameableSequence(nameables []Nameable) []NameableSequenceItem {
	nameableSeq := make([]NameableSequenceItem, 0)
	for _, n := range nameables {
		for _, name := range n.GetNames() {
			nameableSeq = append(nameableSeq, NameableSequenceItem{name, n})
		}
	}
	sort.SliceStable(nameableSeq, func(i, j int) bool {
		return len(nameableSeq[i].Name) > len(nameableSeq[j].Name)
	})
	return nameableSeq
}
