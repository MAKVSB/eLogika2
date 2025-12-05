package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type QuestionAdminDTO struct {
	ID      uint `json:"id"`
	Version uint `json:"version"`

	Title              string                   `json:"title"`
	Content            *models.TipTapContent    `json:"content" ts_type:"JSONContent"`
	TimeToRead         int                      `json:"timeToRead"`
	TimeToProcess      int                      `json:"timeToProcess"`
	QuestionType       enums.QuestionTypeEnum   `json:"questionType"`
	QuestionFormat     enums.QuestionFormatEnum `json:"questionFormat"`
	IncludeAnswerSpace bool                     `json:"includeAnswerSpace"`
	CreatedBy          QuestionCreatedByDTO     `json:"createdBy"`
	Active             bool                     `json:"active"`
	ChapterID          uint                     `json:"chapterId"`
	CategoryID         *uint                    `json:"categoryId" ts_type:"number | null"`
	Steps              []uint                   `json:"steps"`

	Answers   []QuestionAnswerAdminDTO `json:"answers"`
	CheckedBy []QuestionCheckedByDTO   `json:"checkedBy"`
}

func (m QuestionAdminDTO) From(d *models.Question) QuestionAdminDTO {
	questionDTO := QuestionAdminDTO{
		ID:      d.ID,
		Version: d.Version,

		Title:              d.Title,
		Content:            d.Content,
		TimeToRead:         d.TimeToRead,
		TimeToProcess:      d.TimeToProcess,
		QuestionType:       enums.QuestionTypeEnum(d.QuestionType),
		QuestionFormat:     enums.QuestionFormatEnum(d.QuestionFormat),
		IncludeAnswerSpace: d.IncludeAnswerSpace,
		ChapterID:          d.CourseLink.ChapterID,
		CategoryID:         d.CourseLink.CategoryID,
		CreatedBy:          QuestionCreatedByDTO{}.From(d.CreatedBy),
		Active:             d.Active,
		Steps:              make([]uint, len(d.CourseLink.Steps)),

		CheckedBy: make([]QuestionCheckedByDTO, len(d.CheckedBy)),
		Answers:   make([]QuestionAnswerAdminDTO, len(d.Answers)),
	}

	for i, step := range d.CourseLink.Steps {
		questionDTO.Steps[i] = step.ID
	}

	for i, userCheck := range d.CheckedBy {
		questionDTO.CheckedBy[i] = QuestionCheckedByDTO{}.From(&userCheck)
	}

	for i, answer := range d.Answers {
		questionDTO.Answers[i] = QuestionAnswerAdminDTO{}.From(answer.Answer)
	}

	return questionDTO
}
