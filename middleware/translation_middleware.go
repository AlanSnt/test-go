package middleware

import (
	"unyc/json-csv-converter/translations"

	"github.com/gin-gonic/gin"
)

// I18nMiddleware
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Request.FormValue("lang")
		accept := c.Request.Header.Get("Accept-Language")

		translations.NewLocalizer(lang, accept)

		c.Next()
	}
}
