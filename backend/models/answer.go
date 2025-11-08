package models

import (
	"time"

	"gorm.io/gorm"
)

type Answer struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	Content      *TipTapContent `gorm:"serializer:json;type:varbinary"` // The text of the answer
	ContentFiles []*File        `gorm:"many2many:answer_content_files;"`
	Explanation  *TipTapContent `gorm:"serializer:json;type:varbinary"` // Text explaining the (in)correctness of the answer
	TimeToSolve  int            ``                                      // Time it takes to check the correctness of the answer
	Correct      bool           ``                                      // If answer is true

	Question *QuestionAnswer `` // Question it relates to
}

func (Answer) TableName() string {
	return "answers"
}
