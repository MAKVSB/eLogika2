package dtos

import (
	"encoding/json"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type TestInstanceQuestionDTO struct {
	ID                   uint                            `json:"id"`
	TextAnswer           json.RawMessage                 `json:"textAnswer"`
	TextAnswerReviewed   bool                            `json:"textAnswerReviewed,omitempty"`
	TextAnswerPercentage uint                            `json:"textAnswerPercentage"`
	Answers              []TestInstanceQuestionAnswerDTO `json:"answers"`
	OpenAnswers          []TestInstanceOpenAnswerDTO     `json:"openAnswers,omitempty"`
	BlockID              uint                            `json:"blockId"`
	Order                uint                            `json:"order"`
	Title                *string                         `json:"title"`
	Content              json.RawMessage                 `json:"content"`
	QuestionFormat       enums.QuestionFormatEnum        `json:"questionFormat"`
}

func (m TestInstanceQuestionDTO) From(d *models.TestInstanceQuestion, isTutor bool) TestInstanceQuestionDTO {
	dto := TestInstanceQuestionDTO{
		ID:             d.ID,
		TextAnswer:     d.TextAnswer,
		BlockID:        d.TestQuestion.BlockID,
		Order:          d.TestQuestion.Order,
		Content:        d.TestQuestion.Question.Content,
		QuestionFormat: d.TestQuestion.Question.QuestionFormat,
		Answers:        make([]TestInstanceQuestionAnswerDTO, len(d.Answers)),
	}

	for a_i, a := range d.Answers {
		dto.Answers[a_i] = TestInstanceQuestionAnswerDTO{}.From(&a, isTutor)
	}

	if isTutor {
		dto.Title = &d.TestQuestion.Question.Title
		dto.TextAnswerPercentage = d.TextAnswerPercentage
		dto.OpenAnswers = make([]TestInstanceOpenAnswerDTO, len(d.TestQuestion.OpenAnswers))
		for a_i, a := range d.TestQuestion.OpenAnswers {
			dto.OpenAnswers[a_i] = TestInstanceOpenAnswerDTO{}.From(&a)
		}
	}

	return dto
}
