package middleware

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// bodyLogWriter is a custom response writer that captures the response body.
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body and writes it while passing it to the original ResponseWriter.
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware is a Gin middleware function that logs request information.
func LoggerMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a custom response writer that captures the response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

		// Replace the original ResponseWriter with the custom one
		c.Writer = blw

		c.Next()

		status := c.Writer.Status()

		if (status >= 200 && status <= 299) || status == 304 {
			return
		}

		clientIp := c.ClientIP()
		origin := c.Request.Header.Get("Origin")

		log.WithFields(logrus.Fields{
			"Status":   status,
			"ClientIP": clientIp,
			"Origin":   origin,
			"Method":   c.Request.Method,
			"Path":     c.Request.URL.Path,
		}).Info(blw.body.String())
	}
}
