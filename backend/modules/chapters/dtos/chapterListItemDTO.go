package dtos

import "elogika.vsb.cz/backend/models"

type ChapterListItemDTO struct {
	ID       uint   `json:"id"`
	CourseID uint   `json:"course_id"`
	Name     string `json:"name"`
	Visible  bool   `json:"visible"`
	Order    uint   `json:"order"`
}

func (m ChapterListItemDTO) From(d *models.Chapter) ChapterListItemDTO {
	dto := ChapterListItemDTO{
		ID:       d.ID,
		CourseID: d.CourseID,
		Name:     d.Name,
		Visible:  d.Visible,
		Order:    d.Order,
	}

	return dto
}
