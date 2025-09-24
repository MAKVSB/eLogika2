package dtos

import "elogika.vsb.cz/backend/models"

type CategoryDTO struct {
	ID        uint      `json:"id"`
	Version   uint      `json:"version"`
	Name      string    `json:"name"`
	Steps     []StepDTO `json:"steps"`
	ChapterID uint      `json:"chapterId"`
}

func (m CategoryDTO) From(d *models.Category) CategoryDTO {
	dto := CategoryDTO{
		ID:        d.ID,
		Version:   d.Version,
		Name:      d.Name,
		Steps:     make([]StepDTO, len(d.Steps)),
		ChapterID: d.ChapterID,
	}

	for i, s := range d.Steps {
		dto.Steps[i] = StepDTO{}.From(&s)
	}

	return dto
}
