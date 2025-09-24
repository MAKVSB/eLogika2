package models

import (
	"time"

	"gorm.io/gorm"
)

type ClassTutor struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``

	UserID  uint ``
	ClassID uint ``

	User  *User  ``
	Class *Class ``
}
