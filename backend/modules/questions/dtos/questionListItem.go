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
	ChapterID       uint                     `json:"chapterId"`
	ChapterName     string                   `json:"chapterName"`
	CategoryID      *uint                    `json:"categoryId"`
	CategoryName    *string                  `json:"categoryName"`
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
		ChapterID:       d.CourseLink.ChapterID,
		ChapterName:     d.CourseLink.Chapter.Name,
		CategoryID:      d.CourseLink.CategoryID,
	}

	if d.CourseLink.CategoryID != nil {
		dto.CategoryName = &d.CourseLink.Category.Name
	}

	for i, userCheck := range d.CheckedBy {
		dto.CheckedBy[i] = QuestionCheckedByDTO{}.From(&userCheck)
	}

	return dto
}
