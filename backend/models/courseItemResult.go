package models

import (
	"time"

	"gorm.io/gorm"
)

type CourseItemResult struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	StudentID    uint  ``
	CourseItemID uint  ``
	TermID       uint  ``
	UpdatedByID  *uint ``

	TestInstanceID     *uint `` // If is a result of test. Will contain test reference
	ActivityInstanceID *uint `` // If is a result of activity. Will contain test reference

	Points   float64 ``
	Final    bool    `` // Result is not waiting for manual intervention
	Selected bool    `` // Result is selected as active

	CourseItem *CourseItem ``
	Term       *Term       ``
	UpdatedBy  *User       ``
	Student    *User       ``

	TestInstance     *TestInstance     `` // If is a result of test. Will contain test reference
	ActivityInstance *ActivityInstance `` // If is a result of activity. Will contain test reference
}

func (CourseItemResult) TableName() string {
	return "course_item_results"
}
