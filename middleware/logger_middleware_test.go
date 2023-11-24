package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware(t *testing.T) {
	// Create a new Gin engine with LoggerMiddleware
	router := gin.New()

	// Configure logrus to use a buffer for capturing logs
	var logBuffer bytes.Buffer
	log := logrus.New()
	log.Out = &logBuffer

	// Define a route for testing
	router.POST("/export", LoggerMiddleware(log), func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "Export handled")
	})

	// Create a recorder to capture the response
	w := httptest.NewRecorder()

	// Create a sample request with a JSON payload
	req, _ := http.NewRequest("POST", "/export", bytes.NewBufferString(`{"key": "value"}`))
	req.Header.Set("Origin", "http://example.com")

	// Serve the request using the custom router
	router.ServeHTTP(w, req)

	// Extract the log output from the buffer
	logOutput := logBuffer.String()

	// Check if the log contains specific information
	assert.Contains(t, logOutput, "Method=POST")                   // Check for the status code
	assert.Contains(t, logOutput, "Status=503")                    // Check for the request method and path
	assert.Contains(t, logOutput, "Path=/export")                  // Check for the Origin header
	assert.Contains(t, logOutput, "Origin=\"http://example.com\"") // Check for the JSON payload
}
