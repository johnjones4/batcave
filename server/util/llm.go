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
		"2006-01-02T15:04:05.999999999-0700",
		time.DateOnly,
		time.DateTime,
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}
	for _, format := range formats {
		t, err := time.Parse(format, str)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse \"%s\"", str)
}
