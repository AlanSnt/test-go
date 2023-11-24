package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unyc/json-csv-converter/translations"
)

func RemoveUselessChars(value string) string {
	stringValue := strings.ReplaceAll(value, "\n", " ")
	stringValue = strings.ReplaceAll(stringValue, "\r\n", " ")
	stringValue = strings.TrimSpace(stringValue)

	return stringValue
}

func FormatValue(value interface{}) (string, error) {
	var stringValue string

	switch typedValue := value.(type) {
	case string:
		if typedValue == "" {
			stringValue = " "
		} else {
			stringValue = typedValue
		}
	case bool:
		stringValue = strconv.FormatBool(typedValue)
	case float64:
		stringValue = strconv.FormatFloat(typedValue, 'f', -1, 64)
	case map[string]interface{}:
		if dateValue, dateOK := typedValue["date"].(string); dateOK {
			t, err := time.Parse("2006-01-02 15:04:05.000000", dateValue)
			if err != nil {
				return "", fmt.Errorf(translations.GetTranslation("errors.formatValue.dateFormatNotSupported"))
			}
			stringValue = t.Format("2006/01/02")
		} else {
			return "", fmt.Errorf(translations.GetTranslation("errors.formatValue.formatNotSupported"))
		}
	case int:
		stringValue = strconv.Itoa(typedValue)
	case nil:
		stringValue = " "
	default:
		return "", fmt.Errorf(translations.GetTranslation("errors.formatValue.formatNotSupported"))
	}

	return RemoveUselessChars(stringValue), nil
}
