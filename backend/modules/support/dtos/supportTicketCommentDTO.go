package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
	"gorm.io/gorm"
)

type SupportTicketCommentDTO struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	CreatedBy UserDTO        `json:"createdBy"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`

	Content *models.TipTapContent `json:"content"`
}

func (m SupportTicketCommentDTO) From(d *models.SupportTicketComment) SupportTicketCommentDTO {
	dto := SupportTicketCommentDTO{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		CreatedBy: UserDTO{}.From(&d.CreatedBy),
		UpdatedAt: d.UpdatedAt,
		DeletedAt: d.DeletedAt,
		Content:   d.Content,
	}

	return dto
}
