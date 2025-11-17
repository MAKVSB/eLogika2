package middlewares

import (
	"fmt"
	"runtime/debug"
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
)

func CustomRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error however you want
				fmt.Printf("panic recovered: %s\n%s\n", err, debug.Stack())

				// Replace default 500 response with custom JSON
				jsonObj := &common.ErrorResponse{
					Code:    500,
					Message: "Thread panicked",
				}
				c.Abort()
				c.JSON(jsonObj.Code, jsonObj)

				errorEntry := models.LogError{
					Time:  time.Now(),
					Trace: string(debug.Stack()),
				}

				logEntry := models.GetAccesLogEntry(c)
				if logEntry != nil {
					errorEntry.RequestUUID = &logEntry.UUID
					errorEntry.RequestBody = logEntry.RequestBody
				}

			}
		}()
		c.Next()
	}
}
