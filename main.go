package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"unyc/json-csv-converter/middleware"
	"unyc/json-csv-converter/translations"
	"unyc/json-csv-converter/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Payload represents the data payload for the export request.
type Payload struct {
	FileName  string      `json:"fileName"`
	Records   interface{} `json:"records"`
	Columns   []string    `json:"columns"`
	Delimiter string      `json:"delimiter"`
}

var FILE_PATH string

// sendFile sends a file as an HTTP response and removes it afterward.
func sendFile(path string, fileType string, name string, c *gin.Context) {
	fileName := fmt.Sprintf("%s/%s.%s", path, name, fileType)

	file, err := os.ReadFile(fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "text/"+fileType)
	c.Header("Content-Length", strconv.Itoa(len(file)))
	c.Data(http.StatusOK, "text/"+fileType, file)

	// Remove the final file after sending it
	removeFileErr := os.Remove(fileName)
	if removeFileErr != nil {
		fmt.Println(removeFileErr.Error())
		return
	}
}

// ExportValidation validates the export request payload.
func ExportValidation(c *gin.Context) (string, rune, Payload, error) {
	var data Payload
	delimiter := ';'
	format := c.Query("type")

	if format == "" {
		return format, delimiter, data, fmt.Errorf(translations.GetTranslation("errors.exportValidation.noFormat"))
	}

	if format != "csv" && format != "xlsx" {
		return format, delimiter, data, fmt.Errorf(translations.GetTranslation("errors.exportValidation.invalidFormat"))
	}

	if err := c.BindJSON(&data); err != nil {
		return format, delimiter, data, err
	}

	if data.FileName == "" {
		return format, delimiter, data, fmt.Errorf(translations.GetTranslation("errors.exportValidation.emptyFileName"))
	}

	if data.Records == nil || len(data.Records.([]interface{})) == 0 {
		return format, delimiter, data, fmt.Errorf(translations.GetTranslation("errors.exportValidation.emptyRecords"))
	}

	if data.Columns == nil || len(data.Columns) == 0 {
		return format, delimiter, data, fmt.Errorf(translations.GetTranslation("errors.exportValidation.emptyColumns"))
	}

	if data.Delimiter == "," {
		delimiter = ','
	}

	return format, delimiter, data, nil
}

// exportHandler handles the export request.
func exportHandler(c *gin.Context) {
	format, delimiter, data, err := ExportValidation(c)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}

	processErr := utils.ProcessExport(FILE_PATH, format, data.FileName, data.Records.([]interface{}), data.Columns, delimiter)
	if processErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": processErr.Error()})
		return
	}

	sendFile(FILE_PATH, format, data.FileName, c)
}

func init() {
	filePath := os.Getenv("FILE_PATH")

	if filePath == "" {
		FILE_PATH, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		filePath = FILE_PATH
	}

	FILE_PATH = filePath
}

func main() {
	log := logrus.New()
	router := gin.Default()

	gin.SetMode(os.Getenv("GIN_MODE"))

	router.Use(middleware.I18nMiddleware())

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length", "Content-Disposition", "Content-Type"},
		MaxAge:        12 * time.Hour,
	}))

	// Define routes
	router.POST("/export", middleware.LoggerMiddleware(log), exportHandler)
	router.GET("/healthz/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/healthz/readiness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// Run the HTTP server
	router.Run(":8000")
}
