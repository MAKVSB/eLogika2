package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type TestInstanceOpenAnswerDTO struct {
	ID      uint                  `json:"id"`
	Content *models.TipTapContent `json:"content"`
	Correct bool                  `json:"correct,omitempty"`
}

func (m TestInstanceOpenAnswerDTO) From(d *models.Answer) TestInstanceOpenAnswerDTO {
	dto := TestInstanceOpenAnswerDTO{
		ID:      d.ID,
		Content: d.Content,
		Correct: d.Correct,
	}

	return dto
}
