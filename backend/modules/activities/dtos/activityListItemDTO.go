package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type ActivityListItemDTO struct {
	ID          uint                   `json:"id"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	Participant ActivityParticipantDTO `json:"participant"`
	// CourseItem  *CourseItem ``
}

func (m ActivityListItemDTO) From(d *models.ActivityInstance) ActivityListItemDTO {
	dto := ActivityListItemDTO{
		ID:          d.ID,
		UpdatedAt:   d.UpdatedAt,
		Participant: ActivityParticipantDTO{}.From(d.Participant),
		// CourseItemName: d.CourseItem.Name,
	}

	return dto
}
