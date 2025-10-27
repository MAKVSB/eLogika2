package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type ClassDTO struct {
	ID      uint `json:"id"`
	Version uint `json:"version"`

	Name          string                    `json:"name"`
	Room          string                    `json:"room"`
	Type          enums.ClassTypeEnum       `json:"type"`
	StudyForm     enums.StudyFormEnum       `json:"studyForm"`
	TimeFrom      string                    `json:"timeFrom"`
	TimeTo        string                    `json:"timeTo"`
	Day           enums.WeekDayEnum         `json:"day"`
	WeekParity    enums.WeekParityEnum      `json:"weekParity"`
	StudentLimit  uint                      `json:"studentLimit"`
	Tutors        []ClassUserDTO            `json:"tutors"`
	ImportOptions models.ClassImportOptions `json:"importOptions"`
}

func (m ClassDTO) From(d *models.Class) ClassDTO {
	dto := ClassDTO{
		ID:      d.ID,
		Version: d.Version,

		Name:         d.Name,
		Room:         d.Room,
		Type:         d.Type,
		StudyForm:    d.StudyForm,
		TimeFrom:     d.TimeFrom,
		TimeTo:       d.TimeTo,
		Day:          d.Day,
		WeekParity:   d.WeekParity,
		StudentLimit: d.StudentLimit,
		Tutors:       make([]ClassUserDTO, len(d.Tutors)),
	}

	for i, t := range d.Tutors {
		dto.Tutors[i] = ClassUserDTO{}.From(t.User)
	}

	return dto
}
