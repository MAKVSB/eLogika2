package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type Template struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	CourseID      uint                     ``
	Title         string                   ``
	Description   string                   ``
	MixBlocks     bool                     ``
	MixEverything bool                     ``
	Blocks        []TemplateBlock          `gorm:"foreignKey:TemplateID"`
	CreatedByID   uint                     ``
	ManagedBy     enums.CourseUserRoleEnum `` // Role of user who manages it
	CreatedBy     *User                    ``

	Course *Course ``
}

func (Template) TableName() string {
	return "templates"
}
