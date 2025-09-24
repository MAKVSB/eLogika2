package dtos

import (
	"encoding/json"

	"elogika.vsb.cz/backend/models"
)

type CourseItemActivityDTO struct {
	ID uint `json:"id"`

	Description json.RawMessage `json:"description" ts_type:"JSONContent"`
	// DescriptionFile    []File                     ``
	ExpectedResult json.RawMessage `json:"expectedResult" ts_type:"JSONContent"`
	// ExpectedResultFiles []File                     ``
}

func (m CourseItemActivityDTO) From(d *models.CourseItemActivity) CourseItemActivityDTO {
	dto := CourseItemActivityDTO{
		ID: d.ID,

		Description: d.Description,
		// DescriptionFile    []File                     ``
		ExpectedResult: d.ExpectedResult,
		// ExpectedResultFiles []File                     ``
	}

	return dto
}
