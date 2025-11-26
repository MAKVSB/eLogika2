package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type QuestionListItemDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Title           string                   `json:"title"`
	QuestionType    enums.QuestionTypeEnum   `json:"questionType"`
	QuestionFormat  enums.QuestionFormatEnum `json:"questionFormat"`
	CheckedBy       []QuestionCheckedByDTO   `json:"checkedBy"`
	CreatedBy       QuestionCreatedByDTO     `json:"createdBy"`
	UpdatedBy       QuestionCreatedByDTO     `json:"updatedBy"`
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
		UpdatedAt:       d.UpdatedAt,
		Title:           d.Title,
		QuestionType:    enums.QuestionTypeEnum(d.QuestionType),
		QuestionFormat:  enums.QuestionFormatEnum(d.QuestionFormat),
		CheckedBy:       make([]QuestionCheckedByDTO, len(d.CheckedBy)),
		CreatedBy:       QuestionCreatedByDTO{}.From(d.CreatedBy),
		UpdatedBy:       QuestionCreatedByDTO{}.From(d.UpdatedBy),
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
