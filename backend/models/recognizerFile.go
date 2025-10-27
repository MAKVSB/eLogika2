package models

import (
	"time"

	"gorm.io/gorm"
)

type RecognizerFile struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``

	FileID         uint          ``
	File           *File         ``
	TestInstanceID uint          ``
	UniqueIdent    string        ``
	TestInstance   *TestInstance ``
}
