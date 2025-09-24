package dtos

import (
	"encoding/json"

	"elogika.vsb.cz/backend/models"
)

type TestInstanceQuestionAnswerDTO struct {
	ID       uint            `json:"id"`
	Selected bool            `json:"selected"`
	Order    uint            `json:"order"`
	Content  json.RawMessage `json:"content"`
	Correct  *bool           `json:"correct,omitempty"`
}

func (m TestInstanceQuestionAnswerDTO) From(d *models.TestInstanceQuestionAnswer, isTutor bool) TestInstanceQuestionAnswerDTO {
	dto := TestInstanceQuestionAnswerDTO{
		ID:       d.ID,
		Selected: d.Selected,
		Order:    d.TestQuestionAnswer.Order,
		Content:  d.TestQuestionAnswer.Answer.Content,
	}

	if isTutor {
		dto.Correct = &d.TestQuestionAnswer.Answer.Correct
	}

	return dto
}
