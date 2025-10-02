package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type CourseItemTestDTO struct {
	TestType        enums.QuestionTypeEnum `json:"testType"`
	TestTemplateID  uint                   `json:"testTemplateId"`
	TimeLimit       uint                   `json:"timeLimit"`
	ShowResults     bool                   `json:"showResults"`
	ShowTest        bool                   `json:"showTest"`
	ShowCorrectness bool                   `json:"showCorrectness"`
	AllowOffline    bool                   `json:"allowOffline"`
	IsPaper         bool                   `json:"isPaper"`
	IPRanges        string                 `json:"ipRanges"`
}

func (m CourseItemTestDTO) From(d *models.CourseItemTest) CourseItemTestDTO {
	dto := CourseItemTestDTO{
		TestType:        d.TestType,
		TestTemplateID:  d.TestTemplateID,
		TimeLimit:       d.TimeLimit,
		ShowResults:     d.ShowResults,
		ShowTest:        d.ShowTest,
		ShowCorrectness: d.ShowCorrectness,
		AllowOffline:    d.AllowOffline,
		IsPaper:         d.IsPaper,
		IPRanges:        d.IPRanges,
	}

	return dto
}
