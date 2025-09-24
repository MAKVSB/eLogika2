package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type TemplateBlockSegmentDTO struct {
	ID            uint                     `json:"id"`
	ChapterID     *uint                    `json:"chapterId"`
	CategoryID    *uint                    `json:"categoryId"`
	QuestionCount uint                     `json:"questionCount"`
	StepsMode     *enums.StepSelectionEnum `json:"stepsMode"`
	FilterBy      enums.CategoryFilterEnum `json:"filterBy"`
	Steps         []uint                   `json:"steps"`
	Questions     []uint                   `json:"questions"`
}

func (TemplateBlockSegmentDTO) From(d *models.TemplateBlockSegment) TemplateBlockSegmentDTO {
	dto := TemplateBlockSegmentDTO{
		ID:            d.ID,
		ChapterID:     d.ChapterID,
		CategoryID:    d.CategoryID,
		QuestionCount: d.QuestionCount,
		StepsMode:     d.StepsMode,
		FilterBy:      d.FilterBy,
		Steps:         make([]uint, len(d.Steps)),
		Questions:     make([]uint, len(d.Questions)),
	}

	for i, step := range d.Steps {
		dto.Steps[i] = step.ID
	}

	for i, question := range d.Questions {
		dto.Questions[i] = question.ID
	}

	return dto
}
