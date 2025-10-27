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
	TextAnswerPercentage float64                         `json:"textAnswerPercentage"`
	Answers              []TestInstanceQuestionAnswerDTO `json:"answers"`
	OpenAnswers          []TestInstanceOpenAnswerDTO     `json:"openAnswers,omitempty"`
	BlockID              uint                            `json:"blockId"`
	Order                uint                            `json:"order"`
	Title                *string                         `json:"title,omitempty"`
	Content              json.RawMessage                 `json:"content,omitempty"`
	QuestionFormat       enums.QuestionFormatEnum        `json:"questionFormat"`
}

func (m TestInstanceQuestionDTO) From(
	d *models.TestInstanceQuestion,
	showTutor bool,
	showTest bool,
	showCorrectness bool,
) TestInstanceQuestionDTO {
	dto := TestInstanceQuestionDTO{
		ID:             d.ID,
		TextAnswer:     d.TextAnswer,
		BlockID:        d.TestQuestion.BlockID,
		Order:          d.TestQuestion.Order,
		QuestionFormat: d.TestQuestion.Question.QuestionFormat,
		Answers:        make([]TestInstanceQuestionAnswerDTO, len(d.Answers)),
	}

	if showTest {
		dto.Content = d.TestQuestion.Question.Content
	}

	for a_i, a := range d.Answers {
		dto.Answers[a_i] = TestInstanceQuestionAnswerDTO{}.From(a, showTest, showCorrectness)
	}

	if showTutor {
		dto.Title = &d.TestQuestion.Question.Title
		dto.TextAnswerPercentage = d.TextAnswerPercentage
		dto.OpenAnswers = make([]TestInstanceOpenAnswerDTO, len(d.TestQuestion.OpenAnswers))
		for a_i, a := range d.TestQuestion.OpenAnswers {
			dto.OpenAnswers[a_i] = TestInstanceOpenAnswerDTO{}.From(&a)
		}
	}

	return dto
}
