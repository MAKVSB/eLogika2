package middlewares

import (
	"bytes"
	"io"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logEntry := &models.LogAccess{
			UUID:      uuid.New().String(),
			Time:      start,
			Method:    c.Request.Method,
			URL:       c.Request.URL.Path,
			IPAddress: c.ClientIP(),
		}

		// Read body and refresh
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				logEntry.RequestBody = &bodyBytes
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		c.Set("accessLogEntry", logEntry)
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		logEntry.HTTPStatus = c.Writer.Status()
		if logEntry.HTTPStatus >= 400 {
			responseString := blw.body.String()
			logEntry.Response = &responseString
		}
		logEntry.Duration = float64(time.Since(start).Milliseconds())

		go func() {
			if logEntry.Method == "OPTION" {
				return
			}
			// Place regexes of routes where i want to preserve body in logs here
			if true {
				logEntry.RequestBody = &[]byte{}
			}

			if err := initializers.DB.Create(logEntry).Error; err != nil {
				panic("failed to log entry")
			}
		}()
	}
}
