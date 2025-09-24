package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type QuestionListItemDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	Title           string                   `json:"title"`
	QuestionType    enums.QuestionTypeEnum   `json:"questionType"`
	QuestionFormat  enums.QuestionFormatEnum `json:"questionFormat"`
	CheckedBy       []QuestionCheckedByDTO   `json:"checkedBy"`
	CreatedBy       QuestionCreatedByDTO     `json:"createdBy"`
	Active          bool                     `json:"active"`
	QuestionGroupID uint                     `json:"questionGroupId"`
	// ChapterID      uint                     `json:"chapterId"`
	// CategoryID     *uint                    `json:"categoryId"`
}

func (m QuestionListItemDTO) From(d *models.Question) QuestionListItemDTO {
	dto := QuestionListItemDTO{
		ID:              d.ID,
		CreatedAt:       d.CreatedAt,
		Title:           d.Title,
		QuestionType:    enums.QuestionTypeEnum(d.QuestionType),
		QuestionFormat:  enums.QuestionFormatEnum(d.QuestionFormat),
		CheckedBy:       make([]QuestionCheckedByDTO, len(d.CheckedBy)),
		CreatedBy:       QuestionCreatedByDTO{}.From(d.CreatedBy),
		Active:          d.Active,
		QuestionGroupID: d.QuestionGroupID,
		// ChapterID: d.ChapterID,
		// CategoryID: d.CategoryID,
	}

	for i, userCheck := range d.CheckedBy {
		dto.CheckedBy[i] = QuestionCheckedByDTO{}.From(&userCheck)
	}

	return dto
}
