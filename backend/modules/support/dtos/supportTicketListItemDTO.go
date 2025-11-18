package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type SupportTicketListItemDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy UserDTO   `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy *UserDTO  `json:"updatedBy"`
	Solved    bool      `json:"solved"`
}

func (m SupportTicketListItemDTO) From(d *models.SupportTicket) SupportTicketListItemDTO {
	dto := SupportTicketListItemDTO{
		ID:        d.ID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		CreatedBy: UserDTO{}.From(&d.CreatedBy),
		UpdatedAt: d.UpdatedAt,
		Solved:    d.Solved,
	}

	if d.UpdatedBy != nil {
		ub := UserDTO{}.From(d.UpdatedBy)
		dto.UpdatedBy = &ub
	}

	return dto
}
