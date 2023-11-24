package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"unyc/json-csv-converter/translations"

	"github.com/xuri/excelize/v2"
)

// CHUNK_SIZE specifies the size of data chunks for processing.
var CHUNK_SIZE int
var MAX_WORKERS int

var wg sync.WaitGroup

func clearFiles(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := file.Name()[len(file.Name())-3:]
		path := fmt.Sprintf("%s/%s", path, file.Name())

		if ext == "xlsx" {
			os.Remove(path)
		} else if ext == "csv" {
			os.Remove(path)
		} else if ext == "json" {
			os.Remove(path)
		}
	}
}

func mergeExcelFiles(path, fileName string, chunkNumber int) error {
	workerSem := make(chan struct{}, MAX_WORKERS)
	errorChan := make(chan error, chunkNumber)

	finalFilename := fmt.Sprintf("%s/%s.xlsx", path, fileName)
	finalRows := 0

	f := excelize.NewFile()
	defer f.Close()

	_, err := f.NewSheet(SHEET_NAME)
	if err != nil {
		return err
	}

	for i := 0; i < chunkNumber; i++ {
		workerSem <- struct{}{} // Acquire a worker
		wg.Add(1)

		go func(index int) {
			defer func() {
				<-workerSem // Release a worker
				wg.Done()
			}()

			tempFilename := fmt.Sprintf("%s/%s_chunk_%d.xlsx", path, fileName, index)

			xlFile, err := excelize.OpenFile(tempFilename)
			if err != nil {
				errorChan <- err
			}

			// Copy the sheets from temp Excel files to the final Excel file
			for _, sheetName := range xlFile.GetSheetList() {
				rows, err := xlFile.GetRows(sheetName)
				if err != nil {
					errorChan <- err
				}

				// Set the data in the new sheet
				for rowIndex, row := range rows {
					for colIndex, cellValue := range row {
						cellName, err := excelize.CoordinatesToCellName(colIndex+1, finalRows+rowIndex+1)
						if err != nil {
							errorChan <- err
						}

						f.SetCellValue(SHEET_NAME, cellName, cellValue)
					}
					finalRows++
				}
				xlFile.Close()
				os.Remove(tempFilename)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		clearFiles(path)
		return err // Return the first error
	}

	// Close the final Excel file
	err = f.SaveAs(finalFilename)
	if err != nil {
		return err
	}

	return nil
}

// mergeFiles merges temporary files into a final file.
func mergeCsvFiles(path string, fileName string, chunkNumber int) error {
	workerSem := make(chan struct{}, MAX_WORKERS)
	errorChan := make(chan error, chunkNumber)

	finalFilename := fmt.Sprintf("%s/%s.csv", path, fileName)
	finalFile, err := os.Create(finalFilename)
	if err != nil {
		return err
	}
	defer finalFile.Close()

	// Read each temporary file, write its content to the final file, and remove it
	for i := 0; i < chunkNumber; i++ {
		workerSem <- struct{}{} // Acquire a worker
		wg.Add(1)

		go func(index int) {
			defer func() {
				<-workerSem // Release a worker
				wg.Done()
			}()

			tempFilename := fmt.Sprintf("%s/%s_chunk_%d.csv", path, fileName, index)
			tempData, err := os.ReadFile(tempFilename)
			if err != nil {
				errorChan <- err
			}

			_, err = finalFile.Write(tempData)
			if err != nil {
				errorChan <- err
			}

			err = os.Remove(tempFilename)
			if err != nil {
				errorChan <- err
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		clearFiles(path)
		return err // Return the first error
	}

	return nil
}

func chunkPayload(path string, fileName string, records []interface{}) int {
	var chunkNumber int

	workerSem := make(chan struct{}, MAX_WORKERS)
	errorChan := make(chan error, len(records)/CHUNK_SIZE+1)

	// Split records into chunks for parallel processing
	for i := 0; i < len(records); i += CHUNK_SIZE {
		workerSem <- struct{}{} // Acquire a worker
		wg.Add(1)

		end := i + CHUNK_SIZE

		if end > len(records) {
			end = len(records)
		}

		chunk := records[i:end]

		go func(index int, chunk []interface{}) {
			defer func() {
				<-workerSem // Release a worker
				wg.Done()
			}()

			chunkData, err := json.Marshal(chunk)

			// Handle errors
			if err != nil {
				errorChan <- err
				return
			}

			if len(chunkData) == 0 {
				return
			}

			// Write the processed data to a temporary file
			tempFilename := fmt.Sprintf("%s/%s_chunk_%d.json", path, fileName, index)

			err = os.WriteFile(tempFilename, chunkData, 0644)
			if err != nil {
				errorChan <- err
				return
			}
		}(i/CHUNK_SIZE, chunk)

		chunkNumber++
	}

	// Wait for all chunks to be processed
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	return chunkNumber
}

func processChunk(path string, fileName string, delimiter rune, exportType string, columns []string, chunkNumber int) error {
	var chunkData []byte
	var err error

	tempPayloadFilename := fmt.Sprintf("%s/%s_chunk_%d.json", path, fileName, chunkNumber)
	tempFileName := fmt.Sprintf("%s/%s_chunk_%d.%s", path, fileName, chunkNumber, exportType)

	chunkByte, err := os.ReadFile(tempPayloadFilename)
	if err != nil {
		return err
	}

	// convert the chunk to a slice of records to procces in export
	var chunk []interface{}
	err = json.Unmarshal(chunkByte, &chunk)
	if err != nil {
		return err
	}

	// Process a chunk based on the exportType
	if exportType == "csv" {
		chunkData, err = ExportToCSV(chunk, columns, delimiter)
	} else if exportType == "xlsx" {
		chunkData, err = ExportToExcel(chunk, columns)
	}

	// Handle errors
	if err != nil {
		return err
	}

	if len(chunkData) == 0 {
		return fmt.Errorf(translations.GetTranslation("errors.export.emptyChunk"))
	}

	err = os.WriteFile(tempFileName, chunkData, 0644)
	if err != nil {
		return err
	}

	defer os.Remove(tempPayloadFilename)

	return nil
}

// ProcessExport exports data to a file based on the given exportType.
func ProcessExport(
	path string,
	exportType string,
	fileName string,
	records []interface{},
	columns []string,
	delimiter rune,
) error {
	chunksNumbers := chunkPayload(path, fileName, records)

	var wg sync.WaitGroup
	workerSem := make(chan struct{}, MAX_WORKERS)
	errorChan := make(chan error, chunksNumbers)

	for i := 0; i < chunksNumbers; i++ {
		workerSem <- struct{}{} // Acquire a worker
		wg.Add(1)

		go func(index int) {
			defer func() {
				<-workerSem // Release a worker
				wg.Done()
			}()

			err := processChunk(path, fileName, delimiter, exportType, columns, index)
			if err != nil {
				errorChan <- err
			}
		}(i)
	}

	// Wait for all chunks to be processed
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Handle errors from goroutines
	for err := range errorChan {
		clearFiles(path)
		return err // Return the first error
	}

	// Merge temporary files into a final file
	mergeFunc := mergeCsvFiles
	if exportType != "csv" {
		mergeFunc = mergeExcelFiles
	}

	return mergeFunc(path, fileName, chunksNumbers)
}

func init() {
	chunkSizeStr := os.Getenv("CHUNK_SIZE")
	maxWorkersStr := os.Getenv("MAX_WORKERS")

	chunkSize, err := strconv.Atoi(chunkSizeStr)
	if err != nil || chunkSize <= 0 {
		CHUNK_SIZE = 10000
	} else {
		CHUNK_SIZE = chunkSize
	}

	maxWorkers, err := strconv.Atoi(maxWorkersStr)
	if err != nil || maxWorkers <= 0 {
		MAX_WORKERS = 3
	} else {
		MAX_WORKERS = maxWorkers
	}
}
