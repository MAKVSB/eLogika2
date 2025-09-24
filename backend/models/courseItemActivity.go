package models

import "encoding/json"

type CourseItemActivity struct {
	CommonModel
	ID uint `gorm:"primarykey"`

	Description    json.RawMessage ``
	ExpectedResult json.RawMessage ``
	ContentFiles   []File          `gorm:"many2many:course_item_activity_files;"`
}

func (CourseItemActivity) TableName() string {
	return "course_item_activities"
}
