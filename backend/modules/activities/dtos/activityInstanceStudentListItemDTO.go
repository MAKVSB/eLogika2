package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type ActivityInstanceStudentListItemDTO struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	EditableUntil  time.Time `json:"editableUntil"`
	TermName       string    `json:"termName"`
	CourseItemName string    `json:"courseItemName"`
}

func (m ActivityInstanceStudentListItemDTO) From(d *models.ActivityInstance) ActivityInstanceStudentListItemDTO {
	dto := ActivityInstanceStudentListItemDTO{
		ID:             d.ID,
		CreatedAt:      d.CreatedAt,
		EditableUntil:  d.Term.ActiveTo,
		TermName:       d.Term.Name,
		CourseItemName: d.CourseItem.Name,
	}

	return dto
}
