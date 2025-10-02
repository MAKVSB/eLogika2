package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type QuestionCheckedByDTO struct {
	ID           uint      `json:"id"`
	DegreeBefore string    `json:"degreeBefore"`
	FirstName    string    `json:"firstName"`
	FamilyName   string    `json:"familyName"`
	DegreeAfter  string    `json:"degreeAfter"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	CheckedAt    time.Time `json:"checkedAt"`
}

func (m QuestionCheckedByDTO) From(d *models.QuestionCheck) QuestionCheckedByDTO {
	dto := QuestionCheckedByDTO{
		ID:           d.User.ID,
		DegreeBefore: d.User.DegreeBefore,
		FirstName:    d.User.FirstName,
		FamilyName:   d.User.FamilyName,
		DegreeAfter:  d.User.DegreeAfter,
		Username:     d.User.Username,
		Email:        d.User.Email,
		CheckedAt:    d.CreatedAt,
	}

	return dto
}
