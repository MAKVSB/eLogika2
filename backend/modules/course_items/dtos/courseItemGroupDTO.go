package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type CourseItemGroupDTO struct {
	Choice    bool `json:"choice"`
	ChooseMin uint `json:"chooseMin"`
	ChooseMax uint `json:"chooseMax"`
}

func (m CourseItemGroupDTO) From(d *models.CourseItemGroup) CourseItemGroupDTO {
	dto := CourseItemGroupDTO{
		Choice:    d.Choice,
		ChooseMin: d.ChooseMin,
		ChooseMax: d.ChooseMax,
	}

	return dto
}
