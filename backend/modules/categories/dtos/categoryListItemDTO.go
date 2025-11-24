package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type CategoryListItemDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	StepsCount  int    `json:"stepsCount"`
	ChapterID   uint   `json:"chapterId"`
	ChapterName string `json:"chapterName"`
}

func (m CategoryListItemDTO) From(d *models.Category) CategoryListItemDTO {
	dto := CategoryListItemDTO{
		ID:          d.ID,
		Name:        d.Name,
		StepsCount:  len(d.Steps),
		ChapterID:   d.Chapter.ID,
		ChapterName: d.Chapter.Name,
	}

	return dto
}
