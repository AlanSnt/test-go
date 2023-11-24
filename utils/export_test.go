package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestMergeExcelFiles(t *testing.T) {
	tempDir := t.TempDir()

	xlsxFile1 := excelize.NewFile()
	xlsxFile1.SetCellValue(SHEET_NAME, "A1", "Name")
	xlsxFile1.SetCellValue(SHEET_NAME, "B1", "Age")
	xlsxFile1.SetCellValue(SHEET_NAME, "C1", "City")
	xlsxFile1.SetCellValue(SHEET_NAME, "A2", "John")
	xlsxFile1.SetCellValue(SHEET_NAME, "B2", 30)
	xlsxFile1.SetCellValue(SHEET_NAME, "C2", "New York")
	xlsxFile1.SaveAs(tempDir + "/merged_chunk_0.xlsx")

	xlsxFile2 := excelize.NewFile()
	xlsxFile2.SetCellValue(SHEET_NAME, "A1", "Name")
	xlsxFile2.SetCellValue(SHEET_NAME, "B1", "Age")
	xlsxFile2.SetCellValue(SHEET_NAME, "C1", "City")
	xlsxFile2.SetCellValue(SHEET_NAME, "A2", "Alice")
	xlsxFile2.SetCellValue(SHEET_NAME, "B2", 25)
	xlsxFile2.SetCellValue(SHEET_NAME, "C2", "Los Angeles")
	xlsxFile2.SaveAs(tempDir + "/merged_chunk_1.xlsx")

	err := mergeExcelFiles(tempDir, "merged", 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	finalFilePath := fmt.Sprintf("%s/merged.xlsx", tempDir)
	_, err = os.Stat(finalFilePath)
	if err != nil {
		t.Fatalf("The final Excel file does not exist: %v", err)
	}
}

func TestMergeCsvFiles(t *testing.T) {
	tempDir := t.TempDir()

	csvData1 := []byte("Name,Age,City\nJohn,30,New York\n")
	err := os.WriteFile(tempDir+"/merged_chunk_0.csv", csvData1, 0644)
	if err != nil {
		t.Fatalf("csvData1 Unexpected error: %v", err)
	}

	csvData2 := []byte("Name,Age,City\nAlice,25,Los Angeles\n")
	err = os.WriteFile(tempDir+"/merged_chunk_1.csv", csvData2, 0644)
	if err != nil {
		t.Fatalf("csvData2 Unexpected error: %v", err)
	}

	err = mergeCsvFiles(tempDir, "merged", 2)
	if err != nil {
		t.Fatalf("mergeCsvFiles Unexpected error: %v", err)
	}

	finalFilePath := fmt.Sprintf("%s/merged.csv", tempDir)
	_, err = os.Stat(finalFilePath)
	if err != nil {
		t.Fatalf("The final CSV file does not exist: %v", err)
	}
}

func TestProcessExport(t *testing.T) {
	tempDir := t.TempDir()

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
	}
	columns := []string{"name", "age", "city"}

	err := ProcessExport(tempDir, "csv", "exported", records, columns, ',')
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	finalFilePath := fmt.Sprintf("%s/exported.csv", tempDir)
	_, err = os.Stat(finalFilePath)
	if err != nil {
		t.Fatalf("The final CSV file does not exist: %v", err)
	}

	err = ProcessExport(tempDir, "xlsx", "exported", records, columns, ',')
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	finalFilePath = fmt.Sprintf("%s/exported.xlsx", tempDir)
	_, err = os.Stat(finalFilePath)
	if err != nil {
		t.Fatalf("The final Excel file does not exist: %v", err)
	}
}
