package models

import (
	"time"

	"gorm.io/gorm"
)

type Term struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	CourseID     uint      ``
	CourseItemID uint      ``
	Name         string    ``
	ActiveFrom   time.Time ``
	ActiveTo     time.Time ``
	RequiresSign bool      ``
	SignInFrom   time.Time ``
	SignInTo     time.Time ``
	SignOutFrom  time.Time ``
	SignOutTo    time.Time ``
	OfflineTo    time.Time ``
	Classroom    string    ``
	StudentsMax  uint      ``
	PerClass     bool      ``
	Tries        uint      ``

	Students   []UserTerm  ``
	Course     *Course     ``
	CourseItem *CourseItem ``
}

func (Term) TableName() string {
	return "terms"
}
