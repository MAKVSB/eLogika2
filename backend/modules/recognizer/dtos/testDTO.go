package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type TestDTO struct {
	ID uint `json:"id"`

	Questions []TestQuestionDTO `json:"questions"`
}

func (m TestDTO) From(d *models.Test) TestDTO {
	dto := TestDTO{
		ID: d.ID,
	}

	dto.Questions = make([]TestQuestionDTO, len(d.Questions))
	for q_i, q := range d.Questions {
		dto.Questions[q_i] = TestQuestionDTO{}.From(q)
	}

	return dto
}
