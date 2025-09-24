package dtos

import "elogika.vsb.cz/backend/models"

type TemplateCreatorCategoryStepDTO struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Difficulty uint   `json:"difficulty"`
}

func (TemplateCreatorCategoryStepDTO) From(d *models.Step) TemplateCreatorCategoryStepDTO {
	dto := TemplateCreatorCategoryStepDTO{
		ID:         d.ID,
		Name:       d.Name,
		Difficulty: d.Difficulty,
	}

	return dto
}
