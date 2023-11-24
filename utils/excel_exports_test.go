package utils

import (
	"bytes"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestWriteExcelRecord(t *testing.T) {
	// Create a new Excel file.
	f := excelize.NewFile()

	// Define test columns and data
	columns := []string{"name", "age", "city"}
	record := map[string]interface{}{
		"name": "John",
		"age":  30,
		"city": "New York",
	}
	rowIndex := 0

	// Call the test function
	err := writeExcelRecord(f, record, columns, rowIndex)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Get the cell value from the Excel file
	cellValue, err := f.GetCellValue(SHEET_NAME, "A2")
	if err != nil {
		t.Fatalf("Failed to get cell value: %v", err)
	}

	// Check the result
	expected := "John"
	if cellValue != expected {
		t.Errorf("Incorrect result. Expected: %s, Got: %s", expected, cellValue)
	}
}

func TestExportToExcel(t *testing.T) {
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
	result, err := ExportToExcel(records, columns)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Open the Excel file
	xlsx, err := excelize.OpenReader(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("Failed to open Excel file: %v", err)
	}

	// Check the column headers
	for i, column := range columns {
		cellName, _ := excelize.CoordinatesToCellName(i+1, 1)
		cellValue, err := xlsx.GetCellValue(SHEET_NAME, cellName)
		if err != nil {
			t.Fatalf("Failed to get cell value: %v", err)
		}

		if cellValue != column {
			t.Errorf("Incorrect result for header. Expected: %s, Got: %s", column, cellValue)
		}
	}

	// Check the data
	for rowIndex, record := range records {
		recordMap, ok := record.(map[string]interface{})
		if !ok {
			t.Fatalf("Invalid record type at index %d", rowIndex)
		}

		for colIndex, column := range columns {
			cellName, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			cellValue, err := xlsx.GetCellValue(SHEET_NAME, cellName)
			if err != nil {
				t.Fatalf("Failed to get cell value: %v", err)
			}

			stringValue, err := FormatValue(recordMap[column])
			if err != nil {
				t.Fatalf("Failed to format cell value: %v", err)
			}

			if cellValue != stringValue {
				t.Errorf("Incorrect result at row %d, column %s. Expected: %s, Got: %s", rowIndex+2, column, stringValue, cellValue)
			}
		}
	}
}
