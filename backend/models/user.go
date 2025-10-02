package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type NotificationLevel struct {
	Results  bool `json:"results"`
	Messages bool `json:"messages"`
	Terms    bool `json:"terms"`
}

type NotificationDiscord struct {
	Level  NotificationLevel `json:"level"`
	UserID string            `json:"userId"`
}

type NotificationEmail struct {
	Level NotificationLevel `json:"level"`
}

type NotificationPush struct {
	Level NotificationLevel `json:"level"`
	Token string
}

type UserNotification struct {
	Discord NotificationDiscord `json:"discord"`
	Email   NotificationEmail   `json:"email"`
	Push    NotificationPush    `json:"push"`
}

type User struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	DegreeBefore       string                     ``                       // Degree before name
	FirstName          string                     ``                       // First name of the user
	FamilyName         string                     ``                       // Family name of the user
	DegreeAfter        string                     ``                       // Degree after name
	Username           string                     ``                       // Username
	Password           string                     ``                       // Hashed password
	Email              string                     ``                       // Email of the user
	Notification       UserNotification           `gorm:"serializer:json"` // Notification setting
	Type               enums.UserTypeEnum         ``                       // System-wide user permissions
	IdentityProvider   enums.IdentityProviderEnum ``                       // Provider to link user data with
	IdentityProviderID string                     ``                       // Data to link provider entity

	UserCourses []*CourseUser ``

	// Temp helper data
	Results []*CourseItemResult `gorm:"foreignKey:StudentID"`
}

func (user User) FullName() string {
	return user.DegreeBefore + " " + user.FirstName + " " + user.FamilyName + " " + user.DegreeAfter
}

func (User) TableName() string {
	return "users"
}
