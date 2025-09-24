package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type QuestionCheckedByDTO struct {
	ID         uint      `json:"id"`
	FirstName  string    `json:"firstName"`
	FamilyName string    `json:"familyName"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	CheckedAt  time.Time `json:"checkedAt"`
}

func (m QuestionCheckedByDTO) From(d *models.QuestionCheck) QuestionCheckedByDTO {
	dto := QuestionCheckedByDTO{
		ID:         d.User.ID,
		FirstName:  d.User.FirstName,
		FamilyName: d.User.FamilyName,
		Username:   d.User.Username,
		Email:      d.User.Email,
		CheckedAt:  d.CreatedAt,
	}

	return dto
}
