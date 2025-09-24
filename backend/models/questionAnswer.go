package models

import (
	"time"

	"gorm.io/gorm"
)

type QuestionAnswer struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	QuestionID uint ``
	AnswerID   uint ``

	Question *Question ``
	Answer   *Answer   ``
}

func (QuestionAnswer) TableName() string {
	return "question_answers"
}
