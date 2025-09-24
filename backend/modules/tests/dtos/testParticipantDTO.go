package dtos

import "elogika.vsb.cz/backend/models"

type TestParticipantDTO struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstName"`
	FamilyName string `json:"familyName"`
	Username   string `json:"username"`
	Email      string `json:"email"`
}

func (m TestParticipantDTO) From(d *models.User) TestParticipantDTO {
	dto := TestParticipantDTO{
		ID:         d.ID,
		FirstName:  d.FirstName,
		FamilyName: d.FamilyName,
		Username:   d.Username,
		Email:      d.Email,
	}

	return dto
}
