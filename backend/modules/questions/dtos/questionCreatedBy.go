package dtos

import "elogika.vsb.cz/backend/models"

type QuestionCreatedByDTO struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstName"`
	FamilyName string `json:"familyName"`
	Username   string `json:"username"`
	Email      string `json:"email"`
}

func (m QuestionCreatedByDTO) From(d *models.User) QuestionCreatedByDTO {
	dto := QuestionCreatedByDTO{
		ID:         d.ID,
		FirstName:  d.FirstName,
		FamilyName: d.FamilyName,
		Username:   d.Username,
		Email:      d.Email,
	}

	return dto
}
