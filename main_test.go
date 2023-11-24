package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestExportValidation_ValidPayload(t *testing.T) {
	// Create a GIN router in test mode
	gin.SetMode(gin.TestMode)
	gin.New()

	// Create a simulated HTTP response recorder
	recorder := httptest.NewRecorder()

	// Create a fake HTTP request with a valid payload
	req := httptest.NewRequest("POST", "/export?type=csv", nil)
	req.Header.Add("Content-Type", "application/json")

	// Set the request body with a valid payload
	requestBody := `{
			"fileName": "test.csv",
			"records": [{"field1": "value1", "field2": "value2"}],
			"columns": ["field1", "field2"],
			"delimiter": ","
	}`
	req.Body = io.NopCloser(strings.NewReader(requestBody))

	// Create a GIN context from the request
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	// Call the ExportValidation function
	format, delimiter, _, err := ExportValidation(c)

	// Check that the function returns the expected values
	if err != nil {
		t.Errorf("ExportValidation returned an unexpected error: %v", err)
	}

	if format != "csv" {
		t.Errorf("Expected format to be 'csv', but it was %s", format)
	}

	if delimiter != ',' {
		t.Errorf("Expected delimiter to be ',', but it was %c", delimiter)
	}

	// You can also check other values in 'data' here

	// Check the HTTP status code of the simulated response
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected HTTP status code to be %d, but it was %d", http.StatusOK, recorder.Code)
	}
}

func TestExportValidation_InvalidPayload(t *testing.T) {
	// Create a GIN router in test mode
	gin.SetMode(gin.TestMode)
	gin.New()

	// Create a simulated HTTP response recorder
	w := httptest.NewRecorder()

	// Create a fake HTTP request with an invalid payload (e.g., missing fileName)
	req := httptest.NewRequest("POST", "/export?type=csv", nil)
	req.Header.Add("Content-Type", "application/json")

	// Create a GIN context from the request
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	t.Run("No fileName", func(t *testing.T) {
		requestBody := `{
			"records": [{"field1": "value1", "field2": "value2"}],
			"columns": ["field1", "field2"]
		}`
		req.Body = io.NopCloser(strings.NewReader(requestBody))

		// Call ExportValidation and capture the error
		_, _, _, err := ExportValidation(c)
		if err == nil {
			t.Errorf("ExportValidation did not return an error")
		}
	})

	t.Run("No Records", func(t *testing.T) {
		requestBody := `{
			"fileName": "test.csv",
			"records": [],
			"columns": ["field1", "field2"]
		}`
		req.Body = io.NopCloser(strings.NewReader(requestBody))

		// Call ExportValidation and capture the error
		_, _, _, err := ExportValidation(c)
		if err == nil {
			t.Errorf("ExportValidation did not return an error")
		}
	})

	t.Run("No Columns", func(t *testing.T) {
		requestBody := `{
			"fileName": "test.csv",
			"records": [{"field1": "value1", "field2": "value2"}],
			"columns": []
		}`
		req.Body = io.NopCloser(strings.NewReader(requestBody))

		// Call ExportValidation and capture the error
		_, _, _, err := ExportValidation(c)
		if err == nil {
			t.Errorf("ExportValidation did not return an error")
		}
	})

	t.Run("No Delimiter", func(t *testing.T) {
		requestBody := `{
			"fileName": "test.csv",
			"records": [{"field1": "value1", "field2": "value2"}],
			"columns": ["field1", "field2"]
		}`
		req.Body = io.NopCloser(strings.NewReader(requestBody))

		// Call ExportValidation and capture the error
		_, delimiter, _, _ := ExportValidation(c)

		if delimiter != ';' {
			t.Errorf("Expected delimiter to be ';', but it was %c", delimiter)
		}
	})
}
