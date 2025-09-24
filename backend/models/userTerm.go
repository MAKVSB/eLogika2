package models

import (
	"time"

	"gorm.io/gorm"
)

type UserTerm struct {
	CommonModel
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      ``
	CreatedByID uint           ``
	CreatedBy   *User          ``
	DeletedAt   gorm.DeletedAt ``
	DeletedByID *uint          ``
	DeletedBy   *User          ``

	UserID uint ``
	TermID uint ``

	User *User ``
	Term *Term ``
}

func (UserTerm) TableName() string {
	return "user_terms"
}
