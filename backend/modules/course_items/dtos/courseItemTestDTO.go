package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type CourseItemTestDTO struct {
	ID uint `json:"id"`

	TestType       enums.QuestionTypeEnum `json:"testType"`
	TestTemplateID uint                   `json:"testTemplateId"`
	TimeLimit      uint                   `json:"timeLimit"`
	ShowResults    bool                   `json:"showResults"`
	ShowTest       bool                   `json:"showTest"`
	AllowOffline   bool                   `json:"allowOffline"`
	IsPaper        bool                   `json:"isPaper"`
	IPRanges       string                 `json:"ipRanges"`
}

func (m CourseItemTestDTO) From(d *models.CourseItemTest) CourseItemTestDTO {
	dto := CourseItemTestDTO{
		ID: d.ID,

		TestType:       d.TestType,
		TestTemplateID: d.TestTemplateID,
		TimeLimit:      d.TimeLimit,
		ShowResults:    d.ShowResults,
		ShowTest:       d.ShowTest,
		AllowOffline:   d.AllowOffline,
		IsPaper:        d.IsPaper,
		IPRanges:       d.IPRanges,
	}

	return dto
}
