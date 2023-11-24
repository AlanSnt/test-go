package utils

import (
	"testing"
)

func TestExportToCSV(t *testing.T) {
	// Define test columns and data
	columns := []string{"name", "age", "city"}
	records := []interface{}{
		map[string]interface{}{
			"name": "John",
			"age":  30,
			"city": "New York",
		},
		map[string]interface{}{
			"name": "Alice",
			"age":  25,
			"city": "Los Angeles",
		},
		map[string]interface{}{
			"name": "Bob",
			"age":  35,
			"city": "Chicago",
		},
	}

	// Call the test function
	data, err := ExportToCSV(records, columns, ',')
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the result with '\r\n' as the line separator
	expected := "John,30,New York\r\nAlice,25,Los Angeles\r\nBob,35,Chicago\r\n"

	actual := string(data)
	if actual != expected {
		t.Errorf("Incorrect result. Expected: %s Got: %s", expected, actual)
	}
}
