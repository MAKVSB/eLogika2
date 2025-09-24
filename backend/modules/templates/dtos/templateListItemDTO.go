package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type TemplateListItemDTO struct {
	ID          uint                 `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Blocks      int                  `json:"blocks"`
	CreatedBy   TemplateCreatedByDTO `json:"createdBy"`
	CreatedAt   time.Time            `json:"createdAt"`
}

func (TemplateListItemDTO) From(d *models.Template) TemplateListItemDTO {
	dto := TemplateListItemDTO{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		Blocks:      len(d.Blocks),
		CreatedBy:   TemplateCreatedByDTO{}.From(d.CreatedBy),
		CreatedAt:   d.CreatedAt,
	}

	return dto
}
