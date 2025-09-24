package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type CourseListItemDTO struct {
	ID        uint               `json:"id"`
	Name      string             `json:"name"`
	Shortname string             `json:"shortname"`
	Public    bool               `json:"public"`
	Year      uint               `json:"year"`
	Semester  enums.SemesterEnum `json:"semester"`
}

func (m CourseListItemDTO) From(d *models.Course) CourseListItemDTO {
	dto := CourseListItemDTO{
		ID:   d.ID,
		Name: d.Name,

		Shortname: d.Shortname,
		Public:    d.Public,
		Year:      d.Year,
		Semester:  d.Semester,
	}

	return dto
}
