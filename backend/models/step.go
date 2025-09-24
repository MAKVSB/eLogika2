package models

import (
	"time"

	"gorm.io/gorm"
)

type Step struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``

	CategoryID uint   ``
	Name       string ``
	Difficulty uint   ``
}

func (Step) TableName() string {
	return "steps"
}
