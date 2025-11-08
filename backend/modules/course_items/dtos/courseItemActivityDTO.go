package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type CourseItemActivityDTO struct {
	Description    *models.TipTapContent `json:"description" ts_type:"JSONContent"`
	ExpectedResult *models.TipTapContent `json:"expectedResult" ts_type:"JSONContent"`
}

func (m CourseItemActivityDTO) From(d *models.CourseItemActivity) CourseItemActivityDTO {
	dto := CourseItemActivityDTO{
		Description:    d.Description,
		ExpectedResult: d.ExpectedResult,
	}

	return dto
}
