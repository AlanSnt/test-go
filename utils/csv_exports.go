package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"unyc/json-csv-converter/translations"
)

type Writer struct {
	Comma   rune // Field delimiter (set to ',' by NewWriter)
	UseCRLF bool // True to use \r\n as the line terminator
}

// writeCsvRecord writes a CSV record to a writer.
func writeCsvRecord(w *csv.Writer, record map[string]interface{}, columns []string) error {
	entry := make([]string, len(columns))

	for i, column := range columns {
		value, ok := record[column]
		if !ok {
			return fmt.Errorf(translations.GetTranslationWithArgs("errors.csvExport.columnNotFound", map[string]interface{}{
				"Column": column,
			}))
		}

		stringValue, err := FormatValue(value)

		if err != nil {
			return err
		}

		if i < len(entry) {
			entry[i] = stringValue
		}
	}

	return w.Write(entry)
}

// ExportToCSV exports a slice of records to a CSV byte slice.
func ExportToCSV(records []interface{}, columns []string, delimiter rune) ([]byte, error) {
	var buffer bytes.Buffer
	csvWriter := csv.NewWriter(&buffer)

	// Set the appropriate delimiter and line terminator
	csvWriter.Comma = delimiter
	csvWriter.UseCRLF = true

	// Loop through records and write them to the buffer
	for _, record := range records {
		recordMap, ok := record.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(translations.GetTranslation("errors.csvExport.recordInvalid"))
		}

		// Process and write the CSV record
		err := writeCsvRecord(csvWriter, recordMap, columns)
		if err != nil {
			return nil, err
		}
	}

	// Flush the CSV writer to ensure all data is written
	csvWriter.Flush()

	// Check for errors during flushing
	if err := csvWriter.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
