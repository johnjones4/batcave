package util

import (
	"log"
	"math"
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

func ArrayContains(a []string, v string) bool {
	for _, v1 := range a {
		if v == v1 {
			return true
		}
	}
	return false
}

func FindClosestMatchString(options []string, corpus string) string {
	tokens := strings.Split(corpus, " ")
	closestMatch := ""
	closestDistance := math.MaxInt
	for _, keyword := range options {
		for start := range tokens {
			subTokens := tokens[start:]
			for end := range subTokens {
				text := strings.Join(subTokens[:end+1], " ")
				distance := levenshtein.DistanceForStrings([]rune(keyword), []rune(text), levenshtein.DefaultOptions)
				if distance < len(keyword)*2 && distance < closestDistance {
					if distance == 0 {
						return keyword
					}
					closestDistance = distance
					closestMatch = keyword
				}
			}
		}
	}
	log.Printf("Match: %s (%d)", closestMatch, closestDistance)
	return closestMatch
}
