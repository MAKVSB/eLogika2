package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type TemplateBlock struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``

	Title                 string                       ``
	ShowName              bool                         ``
	DifficultyFrom        uint                         ``
	DifficultyTo          uint                         ``
	Weight                uint                         ``
	QuestionFormat        enums.QuestionFormatEnum     ``
	QuestionCount         uint                         ``
	AnswerCount           uint                         ``
	AnswerDistribution    enums.AnswerDistributionEnum ``
	WrongAnswerPercentage uint                         ``
	AllowEmptyAnswers     bool                         ``
	MixInsideBlock        bool                         ``
	TemplateID            uint                         ``

	Segments []TemplateBlockSegment `gorm:"foreignKey:TemplateBlockID"`
}

func (TemplateBlock) TableName() string {
	return "template_blocks"
}
