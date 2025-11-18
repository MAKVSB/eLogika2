package models

import (
	"time"

	"gorm.io/gorm"
)

type SupportTicketComment struct {
	CommonModel
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      ``
	CreatedByID uint           ``
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt ``

	CreatedBy User ``

	SupportTicketID uint           ``
	Content         *TipTapContent `gorm:"serializer:json;type:varbinary(max)"` // The text of the answer
	ContentFiles    []*File        `gorm:"many2many:support_ticket_comment_content_files;"`
}
