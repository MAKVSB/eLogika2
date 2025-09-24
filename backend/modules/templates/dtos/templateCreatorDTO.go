package dtos

import "elogika.vsb.cz/backend/models"

type TemplateCreatorDTO struct {
	ID         uint                         `json:"id"`
	Name       string                       `json:"name"`
	ParentID   *uint                        `json:"parentId"`
	Order      uint                         `json:"order"`
	Categories []TemplateCreatorCategoryDTO `json:"categories"`
}

func (TemplateCreatorDTO) From(d *models.Chapter) TemplateCreatorDTO {
	dto := TemplateCreatorDTO{
		ID:         d.ID,
		Name:       d.Name,
		ParentID:   d.ParentID,
		Order:      d.Order,
		Categories: make([]TemplateCreatorCategoryDTO, len(d.Categories)),
	}

	for i, category := range d.Categories {
		dto.Categories[i] = TemplateCreatorCategoryDTO{}.From(category)
	}

	return dto
}
