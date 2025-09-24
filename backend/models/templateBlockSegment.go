package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type TemplateBlockSegment struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``

	TemplateBlockID uint                     ``
	ChapterID       *uint                    ``
	CategoryID      *uint                    ``
	QuestionCount   uint                     ``
	StepsMode       *enums.StepSelectionEnum ``
	FilterBy        enums.CategoryFilterEnum ``

	Chapter   *Chapter        ``
	Category  *Category       ``
	Steps     []Step          `gorm:"many2many:template_segment_steps;"`
	Questions []QuestionGroup `gorm:"many2many:template_segment_questions;"`
}

func (TemplateBlockSegment) TableName() string {
	return "template_block_segments"
}

/*
{
    {
      "steps": [],
      "questions": [],
    },
    {
      "steps": [],
      "questions": [],
    },
    {
      "steps": [],
      "questions": [],
    },
    {
      "steps": [],
      "questions": [],
    },
    {
      "steps": [],
      "questions": [],
    }
  ]
}
*/
