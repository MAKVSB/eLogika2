package dtos

import "elogika.vsb.cz/backend/models"

type TestCreatedByDTO struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstName"`
	FamilyName string `json:"familyName"`
	Username   string `json:"username"`
	Email      string `json:"email"`
}

func (m TestCreatedByDTO) From(d *models.User) TestCreatedByDTO {
	dto := TestCreatedByDTO{
		ID:         d.ID,
		FirstName:  d.FirstName,
		FamilyName: d.FamilyName,
		Username:   d.Username,
		Email:      d.Email,
	}

	return dto
}
