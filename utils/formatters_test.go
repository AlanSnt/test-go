package utils

import (
	"fmt"
	"testing"
)

func TestRemoveUselessChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello\nWorld", "Hello World"},
		{"    Trim Spaces    ", "Trim Spaces"},
		{"NoChange", "NoChange"},
		{"", ""},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := RemoveUselessChars(test.input)
			if result != test.expected {
				t.Errorf("Expected: %s, Got: %s", test.expected, result)
			}
		})
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{"Hello", "Hello"},
		{"", ""},
		{true, "true"},
		{3.14159265359, "3.14159265359"},
		{map[string]interface{}{"date": "2023-09-27 13:45:00.000000"}, "2023/09/27"},
		{42, "42"},
		{nil, ""},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			result, err := FormatValue(test.input)
			if err != nil {
				if err.Error() != "Type format not supported" {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if result != test.expected {
					t.Errorf("Expected: %s, Got: %s", test.expected, result)
				}
			}
		})
	}
}
