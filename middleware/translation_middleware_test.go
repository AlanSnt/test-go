package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"unyc/json-csv-converter/translations"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
)

func TestI18nMiddleware(t *testing.T) {
	router := gin.Default()

	router.Use(I18nMiddleware())

	router.GET("/test", func(c *gin.Context) {
		localizer := translations.GetLocalizer()
		translation := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "errors.exportValidation.noFormat",
		})
		c.JSON(http.StatusOK, gin.H{"translation": translation})
	})

	t.Run("fr-FR", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Accept-Language", "fr-FR")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Format non spécifié")
	})

	t.Run("en-US", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Accept-Language", "en-US")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Format not specified")
	})

}
