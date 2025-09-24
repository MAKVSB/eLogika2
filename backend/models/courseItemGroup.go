package models

type CourseItemGroup struct {
	CommonModel
	ID uint `gorm:"primarykey"`

	Choice    bool ``
	ChooseMin uint ``
	ChooseMax uint ``
}

func (CourseItemGroup) TableName() string {
	return "course_item_groups"
}
