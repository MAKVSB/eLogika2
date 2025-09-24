package models

import (
	"time"

	"gorm.io/gorm"
)

type CourseQuestion struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	CourseID   uint  ``
	QuestionID uint  ``
	ChapterID  uint  `` // ID of the chapter
	CategoryID *uint `` // ID of the category
	Difficulty int   ``

	Course   *Course   ``
	Question *Question ``
	Chapter  *Chapter  ``
	Category *Category ``
	Steps    []Step    `gorm:"many2many:question_steps;"`
}

func (CourseQuestion) TableName() string {
	return "course_questions"
}
