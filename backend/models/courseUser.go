package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type CourseUser struct {
	CommonModel
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``

	CourseID  uint                       ``
	UserID    uint                       ``
	Roles     []enums.CourseUserRoleEnum `gorm:"type:nvarchar(max);serializer:json"` // Role that this user has in course
	StudyForm *enums.StudyFormEnum       ``

	Course *Course ``
	User   *User   ``
}

func (CourseUser) TableName() string {
	return "course_users"
}

func (cu CourseUser) HasRole(role enums.CourseUserRoleEnum) bool {
	for _, r := range cu.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (cu CourseUser) NotHasRole(role enums.CourseUserRoleEnum) bool {
	for _, r := range cu.Roles {
		if r == role {
			return false
		}
	}
	return true
}
