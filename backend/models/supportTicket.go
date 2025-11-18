package models

import (
	"time"

	"gorm.io/gorm"
)

type SupportTicket struct {
	CommonModel
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      ``
	CreatedByID uint           ``
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	UpdatedByID *uint          ``
	DeletedAt   gorm.DeletedAt ``
	DeletedByID *uint          ``

	CreatedBy User  ``
	UpdatedBy *User ``
	DeletedBy *User ``

	// Real content
	Name         string         ``
	Content      *TipTapContent `gorm:"serializer:json;type:varbinary(max)"` // The text of the answer
	ContentFiles []*File        `gorm:"many2many:support_ticket_content_files;"`
	URL          string         ``

	Solved   bool                    ``
	Comments []*SupportTicketComment ``
}
