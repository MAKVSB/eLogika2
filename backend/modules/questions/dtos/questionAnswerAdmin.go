package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type QuestionAnswerAdminDTO struct {
	ID      uint `json:"id"`
	Version uint `json:"version"`

	Content     *models.TipTapContent `json:"content" ts_type:"JSONContent"`
	Explanation *models.TipTapContent `json:"explanation" ts_type:"JSONContent"`
	TimeToSolve int                   `json:"timeToSolve"`
	Correct     bool                  `json:"correct"`
}

func (m QuestionAnswerAdminDTO) From(d *models.Answer) QuestionAnswerAdminDTO {
	dto := QuestionAnswerAdminDTO{
		ID:      d.ID,
		Version: d.Version,

		Content:     d.Content,
		Explanation: d.Explanation,
		TimeToSolve: d.TimeToSolve,
		Correct:     d.Correct,
	}

	return dto
}
