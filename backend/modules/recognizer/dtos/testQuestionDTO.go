package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type TestQuestionDTO struct {
	ID             uint                     `json:"id"`
	AnswerCount    uint                     `json:"answerCount"`
	Order          uint                     `json:"order"`
	QuestionFormat enums.QuestionFormatEnum `json:"questionFormat"`
}

func (m TestQuestionDTO) From(
	d *models.TestQuestion,
) TestQuestionDTO {
	dto := TestQuestionDTO{
		ID:             d.ID,
		Order:          d.Order,
		AnswerCount:    uint(len(d.Answers)),
		QuestionFormat: d.Question.QuestionFormat,
	}

	if d.Question.QuestionFormat == enums.QuestionFormatOpen {
		dto.AnswerCount = 11
	}

	return dto
}
