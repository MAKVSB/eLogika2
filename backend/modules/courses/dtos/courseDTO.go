package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type CourseDTO struct {
	ID      uint `json:"id"`
	Version uint `json:"version"`

	Name          string                     `json:"name"`
	Content       *models.TipTapContent      `json:"content" ts_type:"JSONContent"`
	Shortname     string                     `json:"shortname"`
	Public        bool                       `json:"public"`
	Year          uint                       `json:"year"`
	Semester      enums.SemesterEnum         `json:"semester"`
	PointsMin     float64                    `json:"pointsMin"`
	PointsMax     float64                    `json:"pointsMax"`
	ImportOptions models.CourseImportOptions `json:"importOptions"`
}

func (m CourseDTO) From(d *models.Course) CourseDTO {
	dto := CourseDTO{
		ID:      d.ID,
		Version: d.Version,

		Name:          d.Name,
		Content:       d.Content,
		Shortname:     d.Shortname,
		Public:        d.Public,
		Year:          d.Year,
		Semester:      d.Semester,
		PointsMin:     d.PointsMin,
		PointsMax:     d.PointsMax,
		ImportOptions: d.ImportOptions,
	}

	return dto
}
