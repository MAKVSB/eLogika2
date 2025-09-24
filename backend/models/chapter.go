package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Chapter struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	CourseID uint            ``
	Name     string          ``
	ParentID *uint           ``
	Content  json.RawMessage ``
	Visible  bool            ``
	Order    uint            `gorm:"not null"`

	ContentFiles []*File     `gorm:"many2many:chapter_files;"`
	Parent       *Chapter    ``
	Childs       []*Chapter  `gorm:"foreignKey:ParentID"`
	Categories   []*Category ``
}

func (Chapter) TableName() string {
	return "chapters"
}
