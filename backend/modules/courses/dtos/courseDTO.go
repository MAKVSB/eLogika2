package dtos

import (
	"encoding/json"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type CourseDTO struct {
	ID      uint `json:"id"`
	Version uint `json:"version"`

	Name          string                     `json:"name"`
	Content       json.RawMessage            `json:"content" ts_type:"JSONContent"`
	Shortname     string                     `json:"shortname"`
	Public        bool                       `json:"public"`
	Year          uint                       `json:"year"`
	Semester      enums.SemesterEnum         `json:"semester"`
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
		ImportOptions: d.ImportOptions,
	}

	return dto
}
