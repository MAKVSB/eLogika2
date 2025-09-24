package models

import (
	"time"

	"gorm.io/gorm"
)

type QuestionGroup struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``

	OriginalName string     `` // Original name of the question
	Questions    []Question ``
}

func (QuestionGroup) TableName() string {
	return "question_groups"
}
