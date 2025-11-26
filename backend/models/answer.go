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

	Content          *TipTapContent `gorm:"serializer:json;type:varbinary(max)"` // The text of the answer
	ContentFiles     []*File        `gorm:"many2many:answer_content_files;"`
	Explanation      *TipTapContent `gorm:"serializer:json;type:varbinary(max)"` // Text explaining the (in)correctness of the answer
	ExplanationFiles []*File        `gorm:"many2many:answer_explanation_files;"`
	TimeToSolve      int            ``
	Correct          bool           ``

	Question *QuestionAnswer `` // Question it relates to
}
