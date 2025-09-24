package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	CourseID  uint   ``
	ChapterID uint   ``
	Name      string ``

	Steps   []Step  ``
	Chapter Chapter ``
}

func (Category) TableName() string {
	return "categories"
}
