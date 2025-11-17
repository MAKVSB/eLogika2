package models

import (
	"time"
)

// LogAccess představuje záznam v logovací tabulce.
type LogError struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	RequestUUID *string   `gorm:"type:char(36);uniqueIndex"` // Request UUID ,null if failed system call (cron)
	Time        time.Time ``
	RequestBody *[]byte   ``
	Trace       string    ``
}
