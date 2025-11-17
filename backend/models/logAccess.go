package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common"
	"github.com/gin-gonic/gin"
)

// LogAccess představuje záznam v logovací tabulce.
type LogAccess struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UUID        string    `gorm:"type:char(36);uniqueIndex"`
	URL         string    `gorm:"type:varchar(2048)"`
	Method      string    `gorm:"type:varchar(10)"`
	Time        time.Time ``
	IPAddress   string    `gorm:"type:varchar(50)"`
	UserID      *uint     `gorm:"index"`
	TokenID     *string   ``
	RequestBody *[]byte   ``
	HTTPStatus  int       ``
	Response    *string   ``
	Duration    float64   ``
}

func GetAccesLogEntry(c *gin.Context) *LogAccess {
	logEntryVal, ok := c.Get("accessLogEntry")
	if !ok {
		c.AbortWithStatusJSON(500, common.ErrorResponse{
			Message: "Failed to find access log entry",
		})
		return nil
	}

	logEntry, ok := logEntryVal.(*LogAccess)
	if !ok {
		return nil
	}
	return logEntry
}
