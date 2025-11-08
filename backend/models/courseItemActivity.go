package models

type CourseItemActivity struct {
	CommonModel
	ID uint `gorm:"primarykey"`

	Description         *TipTapContent `gorm:"serializer:json;type:varbinary"`
	ExpectedResult      *TipTapContent `gorm:"serializer:json;type:varbinary"`
	DescriptionFiles    []*File        `gorm:"many2many:course_item_activity_files_description;"`
	ExpectedResultFiles []*File        `gorm:"many2many:course_item_activity_files_result;"`
}

func (CourseItemActivity) TableName() string {
	return "course_item_activities"
}
