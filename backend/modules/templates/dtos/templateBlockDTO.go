package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type TemplateBlockDTO struct {
	ID                    uint                         `json:"id"`
	Title                 string                       `json:"title"`
	ShowName              bool                         `json:"showName"`
	DifficultyFrom        uint                         `json:"difficultyFrom"`
	DifficultyTo          uint                         `json:"difficultyTo"`
	Weight                uint                         `json:"weight"`
	QuestionFormat        enums.QuestionFormatEnum     `json:"questionFormat"`
	QuestionCount         uint                         `json:"questionCount"`
	AnswerCount           uint                         `json:"answerCount"`
	AnswerDistribution    enums.AnswerDistributionEnum `json:"answerDistribution"`
	WrongAnswerPercentage uint                         `json:"wrongAnswerPercentage"`
	MixInsideBlock        bool                         `json:"mixInsideBlock"`
	Segments              []TemplateBlockSegmentDTO    `json:"segments"`
}

func (TemplateBlockDTO) From(d *models.TemplateBlock) TemplateBlockDTO {
	dto := TemplateBlockDTO{
		ID:                    d.ID,
		Title:                 d.Title,
		ShowName:              d.ShowName,
		DifficultyFrom:        d.DifficultyFrom,
		DifficultyTo:          d.DifficultyTo,
		Weight:                d.Weight,
		QuestionFormat:        d.QuestionFormat,
		QuestionCount:         d.QuestionCount,
		AnswerCount:           d.AnswerCount,
		AnswerDistribution:    d.AnswerDistribution,
		WrongAnswerPercentage: d.WrongAnswerPercentage,
		MixInsideBlock:        d.MixInsideBlock,
		Segments:              make([]TemplateBlockSegmentDTO, len(d.Segments)),
	}

	for i, segment := range d.Segments {
		dto.Segments[i] = TemplateBlockSegmentDTO{}.From(&segment)
	}

	return dto
}
