package dtos

import "elogika.vsb.cz/backend/models"

type TemplateCreatorCategoryDTO struct {
	ID    uint                             `json:"id"`
	Name  string                           `json:"name"`
	Steps []TemplateCreatorCategoryStepDTO `json:"steps"`
}

func (TemplateCreatorCategoryDTO) From(d *models.Category) TemplateCreatorCategoryDTO {
	dto := TemplateCreatorCategoryDTO{
		ID:    d.ID,
		Name:  d.Name,
		Steps: make([]TemplateCreatorCategoryStepDTO, len(d.Steps)),
	}

	for i, step := range d.Steps {
		dto.Steps[i] = TemplateCreatorCategoryStepDTO{}.From(&step)
	}

	return dto
}
