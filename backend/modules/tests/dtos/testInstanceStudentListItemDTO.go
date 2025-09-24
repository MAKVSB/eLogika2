package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type TestInstanceStudentListItemDTO struct {
	ID             uint                        `json:"id"`
	CreatedAt      time.Time                   `json:"createdAt"`
	State          enums.TestInstanceStateEnum `json:"state"`
	StartUntil     time.Time                   `json:"startUntil"`
	StartedAt      time.Time                   `json:"startedAt"`
	EndsAt         time.Time                   `json:"endsAt"`
	TermName       string                      `json:"termName"`
	CourseItemName string                      `json:"courseItemName"`
}

func (m TestInstanceStudentListItemDTO) From(d *models.TestInstance) TestInstanceStudentListItemDTO {
	dto := TestInstanceStudentListItemDTO{
		ID:             d.ID,
		CreatedAt:      d.CreatedAt,
		State:          d.State,
		StartUntil:     d.Term.ActiveTo,
		StartedAt:      d.StartedAt,
		EndsAt:         d.EndsAt,
		TermName:       d.Term.Name,
		CourseItemName: d.CourseItem.Name,
	}

	return dto
}
