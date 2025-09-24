package dtos

import (
	"encoding/json"

	"elogika.vsb.cz/backend/models"
)

type ChapterDTO struct {
	ID       uint `json:"id"`
	CourseID uint `json:"courseId"`
	Version  uint `json:"version"`

	Name     string          `json:"name"`
	Content  json.RawMessage `json:"content" ts_type:"JSONContent"`
	Visible  bool            `json:"visible"`
	ParentID *uint           `json:"parentId"`
	Order    uint            `json:"order"`

	Childs []ChapterListItemDTO `json:"childs"`
}

func (m ChapterDTO) From(d *models.Chapter) ChapterDTO {
	dto := ChapterDTO{
		ID:       d.ID,
		CourseID: d.CourseID,
		Version:  d.Version,

		Name:     d.Name,
		Content:  d.Content,
		Visible:  d.Visible,
		ParentID: d.ParentID,
		Order:    d.Order,
		Childs:   make([]ChapterListItemDTO, len(d.Childs)),
	}

	for i, child := range d.Childs {
		dto.Childs[i] = ChapterListItemDTO{}.From(child)
	}

	return dto
}
