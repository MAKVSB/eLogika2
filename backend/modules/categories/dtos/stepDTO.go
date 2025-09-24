package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type StepDTO struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Difficulty uint   `json:"difficulty"`
	Deleted    bool   `json:"deleted"`
}

func (m StepDTO) From(d *models.Step) StepDTO {
	dto := StepDTO{
		ID:         d.ID,
		Name:       d.Name,
		Difficulty: d.Difficulty,
		Deleted:    false,
	}

	return dto
}
