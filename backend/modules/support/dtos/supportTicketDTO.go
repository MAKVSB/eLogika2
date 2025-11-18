package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type SupportTicketDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy UserDTO   `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy *UserDTO  `json:"updatedBy"`

	// Real content
	Name     string                    `json:"name"`
	Content  *models.TipTapContent     `json:"content"`
	Solved   bool                      `json:"solved"`
	Comments []SupportTicketCommentDTO `json:"comments"`
	URL      string                    `json:"url"`
	Editable bool                      `json:"editable"`
}

func (m SupportTicketDTO) From(d *models.SupportTicket, editable bool) SupportTicketDTO {
	dto := SupportTicketDTO{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		CreatedBy: UserDTO{}.From(&d.CreatedBy),
		UpdatedAt: d.UpdatedAt,

		// Real content
		Name:     d.Name,
		Content:  d.Content,
		Solved:   d.Solved,
		URL:      d.URL,
		Comments: make([]SupportTicketCommentDTO, len(d.Comments)),
		Editable: editable,
	}

	if d.UpdatedBy != nil {
		ub := UserDTO{}.From(d.UpdatedBy)
		dto.UpdatedBy = &ub
	}

	for i, comment := range d.Comments {
		dto.Comments[i] = SupportTicketCommentDTO{}.From(comment)
	}

	return dto
}
