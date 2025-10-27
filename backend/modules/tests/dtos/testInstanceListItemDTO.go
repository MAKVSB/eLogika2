package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type TestInstanceListItemDTO struct {
	ID          uint                        `json:"id"`
	CreatedAt   time.Time                   `json:"createdAt"`
	State       enums.TestInstanceStateEnum `json:"state"`
	Form        enums.TestInstanceFormEnum  `json:"form"`
	StartedAt   time.Time                   `json:"startedAt"`
	EndedAt     time.Time                   `json:"endedAt"`
	Participant TestParticipantDTO          `json:"participant"`
	Points      float64                     `json:"points"`
}

func (m TestInstanceListItemDTO) From(d *models.TestInstance) TestInstanceListItemDTO {
	dto := TestInstanceListItemDTO{
		ID:          d.ID,
		CreatedAt:   d.CreatedAt,
		State:       d.State,
		Form:        d.Form,
		StartedAt:   d.StartedAt,
		EndedAt:     d.EndedAt,
		Participant: TestParticipantDTO{}.From(d.Participant),
		Points:      d.Result.Points,
	}

	return dto
}
