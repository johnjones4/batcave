package util

import (
	"testing"
	"time"
)

func TestCleanLLMJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"{sample data}", "{sample data}"},
		{"some text before {sample data}", "{sample data}"},
		{"no curly braces", ""},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := CleanLLMJSON(test.input)
			if result != test.expected {
				t.Errorf("Expected: %s, Got: %s", test.expected, result)
			}
		})
	}
}

func TestParseLLMDate(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
	}{
		{"2023-10-02T10:20:30Z", time.Date(2023, 10, 2, 10, 20, 30, 0, time.UTC)},
		{"2023-10-02T10:20:30.123456789-0700", time.Date(2023, 10, 2, 10, 20, 30, 123456789, time.FixedZone("", -7*60*60))},
		{"invalid", time.Time{}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseLLMDate(test.input)
			if err != nil {
				if !result.IsZero() {
					t.Errorf("Expected zero time for error case, Got: %v", result)
				}
			} else if !result.Equal(test.expected) {
				t.Errorf("Expected: %v, Got: %v", test.expected, result)
			}
		})
	}
}
