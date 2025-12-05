package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type QuestionVersionDTO struct {
	ID      uint `json:"id"`
	Version uint `json:"version"`

	Title     string               `json:"title"`
	CreatedBy QuestionCreatedByDTO `json:"createdBy"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`

	// Modifying old versions
	IsArchiveVersion bool `json:"isArchiveVersion"`
}

func (m QuestionVersionDTO) From(d *models.QuestionVersion) QuestionVersionDTO {
	questionDTO := QuestionVersionDTO{
		ID:      d.ID,
		Version: d.Version,

		Title:            d.Title,
		UpdatedAt:        d.UpdatedAt,
		CreatedAt:        d.CreatedAt,
		CreatedBy:        QuestionCreatedByDTO{}.From(d.CreatedBy),
		IsArchiveVersion: d.CourseLink.DeletedAt.Valid,
	}

	return questionDTO
}
