package models

import (
	"time"
)

type QuestionCheck struct {
	CommonModel
	CreatedAt  time.Time ``
	QuestionID uint      `gorm:"primaryKey"` // Question it relates to
	UserID     uint      `gorm:"primaryKey"` // User that checked question

	User User `` // User that checked question
}

func (QuestionCheck) TableName() string {
	return "question_checks"
}
