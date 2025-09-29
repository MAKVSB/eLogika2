package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

// Email represents an email to be sent.
type Email struct {
	ID        uint                       `gorm:"primaryKey"`                           //
	ToEmail   string                     `gorm:"type:text;not null"`                   //
	Subject   string                     `gorm:"type:text;not null"`                   //
	Body      string                     `gorm:"type:text;not null"`                   //
	Status    enums.EmailQueueStatusEnum `gorm:"type:text;not null;default:'PENDING'"` //
	Retries   int                        `gorm:"not null;default:0"`                   //
	Priority  int                        `gorm:"not null;default:0"`                   // higher number = higher priority
	CreatedAt time.Time                  `gorm:"autoCreateTime"`                       //
	UpdatedAt time.Time                  `gorm:"autoUpdateTime"`                       //
	DeletedAt gorm.DeletedAt             `gorm:"index"`                                // optional soft delete
}

func (Email) TableName() string {
	return "email_queue"
}
