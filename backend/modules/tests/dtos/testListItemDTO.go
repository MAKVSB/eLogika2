package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type TestListItemDTO struct {
	ID        uint             `json:"id"`
	CreatedAt time.Time        `json:"createdAt"`
	Name      string           `json:"name"`
	Term      string           `json:"term"`
	Group     string           `json:"group"`
	CreatedBy TestCreatedByDTO `json:"createdBy"`
}

func (m TestListItemDTO) From(d *models.Test) TestListItemDTO {
	dto := TestListItemDTO{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		Name:      d.Name,
		Term:      d.Term.Name,
		Group:     d.Group,
		CreatedBy: TestCreatedByDTO{}.From(d.CreatedBy),
	}

	return dto
}
