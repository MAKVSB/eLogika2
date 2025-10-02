package dtos

import "elogika.vsb.cz/backend/models"

type TestParticipantDTO struct {
	ID           uint   `json:"id"`
	DegreeBefore string `json:"degreeBefore"`
	FirstName    string `json:"firstName"`
	FamilyName   string `json:"familyName"`
	DegreeAfter  string `json:"degreeAfter"`
	Username     string `json:"username"`
	Email        string `json:"email"`
}

func (m TestParticipantDTO) From(d *models.User) TestParticipantDTO {
	dto := TestParticipantDTO{
		ID:           d.ID,
		DegreeBefore: d.DegreeBefore,
		FirstName:    d.FirstName,
		FamilyName:   d.FamilyName,
		DegreeAfter:  d.DegreeAfter,
		Username:     d.Username,
		Email:        d.Email,
	}

	return dto
}
