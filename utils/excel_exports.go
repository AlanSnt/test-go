package utils

import (
	"fmt"
	"unyc/json-csv-converter/translations"

	"github.com/xuri/excelize/v2"
)

const SHEET_NAME = "Sheet1"

func writeExcelRecord(w *excelize.File, record map[string]interface{}, columns []string, rowIndex int) error {
	for colIndex, column := range columns {
		cellName, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
		value, exists := record[column]
		if exists {
			stringValue, err := FormatValue(value)
			if err != nil {
				return err
			}
			w.SetCellValue(SHEET_NAME, cellName, stringValue)
		}
	}

	return nil
}

// ExportToExcel exports a slice of records to an Excel byte slice.
func ExportToExcel(records []interface{}, columns []string) ([]byte, error) {
	// Create a new Excel file.
	f := excelize.NewFile()

	defer f.Close()

	// Create a new sheet.
	_, err := f.NewSheet(SHEET_NAME)

	if err != nil {
		return nil, err
	}

	// Set the column headers.
	for i, column := range columns {
		cellName, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(SHEET_NAME, cellName, column)
	}

	// Populate the sheet with data.
	for rowIndex, record := range records {
		recordMap, ok := record.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(translations.GetTranslation("errors.excelExport.recordInvalid"))
		}

		err := writeExcelRecord(f, recordMap, columns, rowIndex)

		if err != nil {
			return nil, err
		}
	}

	excelData, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return excelData.Bytes(), nil
}
