package translations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTranslation(t *testing.T) {
	// Set up the localizer with the desired language
	NewLocalizer("fr-FR", "fr-FR")

	// Test a translation for a specific MessageID
	result := GetTranslation("errors.exportValidation.noFormat")
	assert.Equal(t, "Format non spécifié", result, "Expected translation did not match")

	// You can add more similar test cases for other MessageIDs
}

func TestGetTranslationWithArgs(t *testing.T) {
	// Set up the localizer with the desired language
	NewLocalizer("fr-FR", "fr-FR")

	// Test a translation for a specific MessageID with arguments
	args := map[string]interface{}{
		"Column": "ColumnName",
	}

	result := GetTranslationWithArgs("errors.csvExport.columnNotFound", args)
	assert.Equal(t, "Export CSV : Colonne ColumnName non trouvée", result, "Expected translation with args did not match")
}
