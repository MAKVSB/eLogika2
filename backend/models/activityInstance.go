package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityInstance struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``

	ParticipantID uint ``
	TermID        uint ``
	CourseItemID  uint ``

	Content      *TipTapContent `gorm:"serializer:json;type:varbinary"`
	ContentFiles []*File        `gorm:"many2many:activity_instance_content_files;"`

	Participant *User       ``
	Term        *Term       ``
	CourseItem  *CourseItem ``

	Result *CourseItemResult ``
}
