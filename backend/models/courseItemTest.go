package models

import (
	"elogika.vsb.cz/backend/modules/common/enums"
)

type CourseItemTest struct {
	CommonModel
	ID uint `gorm:"primarykey"`

	TestType        enums.QuestionTypeEnum ``
	TestTemplateID  uint                   ``
	TimeLimit       uint                   ``
	ShowResults     bool                   ``
	ShowTest        bool                   ``
	ShowCorrectness bool                   ``
	AllowOffline    bool                   ``
	IsPaper         bool                   ``
	IPRanges        string                 ``

	TestTemplate *Template ``
}

func (CourseItemTest) TableName() string {
	return "course_item_tests"
}
