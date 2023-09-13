package util

import (
	"fmt"
	"strings"
	"time"
)

func CleanLLMJSON(response string) string {
	openBracket := strings.Index(response, "{")
	if openBracket < 0 {
		return ""
	}
	return response[openBracket:]
}

func ParseLLMDate(str string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999-0700",
	}
	for _, format := range formats {
		t, err := time.Parse(format, str)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse \"%s\"", str)
}
