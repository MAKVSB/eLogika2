package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type JoinedStudentDTO struct {
	ID           uint       `json:"id"`
	UserID       uint       `json:"userId"`
	DegreeBefore string     `json:"degreeBefore"`
	FirstName    string     `json:"firstName"`
	FamilyName   string     `json:"familyName"`
	DegreeAfter  string     `json:"degreeAfter"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	CreatedAt    time.Time  `json:"createdAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
}

func (m JoinedStudentDTO) From(d *models.UserTerm) JoinedStudentDTO {
	dto := JoinedStudentDTO{
		ID:           d.ID,
		UserID:       d.User.ID,
		DegreeBefore: d.User.DegreeBefore,
		FirstName:    d.User.FirstName,
		FamilyName:   d.User.FamilyName,
		DegreeAfter:  d.User.DegreeAfter,
		Username:     d.User.Username,
		Email:        d.User.Email,
		CreatedAt:    d.CreatedAt,
	}

	if d.DeletedAt.Valid {
		dto.DeletedAt = &d.DeletedAt.Time
	}

	return dto
}
