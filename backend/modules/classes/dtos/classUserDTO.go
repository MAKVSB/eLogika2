package dtos

import "elogika.vsb.cz/backend/models"

type ClassUserDTO struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstName"`
	FamilyName string `json:"familyName"`
	Username   string `json:"username"`
	Email      string `json:"email"`
}

func (m ClassUserDTO) From(d *models.User) ClassUserDTO {
	dto := ClassUserDTO{
		ID:         d.ID,
		FirstName:  d.FirstName,
		FamilyName: d.FamilyName,
		Username:   d.Username,
		Email:      d.Email,
	}

	return dto
}
