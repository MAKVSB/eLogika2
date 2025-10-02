package dtos

import "elogika.vsb.cz/backend/models"

type JoinedStudentDTO struct {
	ID           uint   `json:"id"`
	DegreeBefore string `json:"degreeBefore"`
	FirstName    string `json:"firstName"`
	FamilyName   string `json:"familyName"`
	DegreeAfter  string `json:"degreeAfter"`
	Username     string `json:"username"`
	Email        string `json:"email"`
}

func (m JoinedStudentDTO) From(d *models.UserTerm) JoinedStudentDTO {
	dto := JoinedStudentDTO{
		ID:           d.User.ID,
		DegreeBefore: d.User.DegreeBefore,
		FirstName:    d.User.FirstName,
		FamilyName:   d.User.FamilyName,
		DegreeAfter:  d.User.DegreeAfter,
		Username:     d.User.Username,
		Email:        d.User.Email,
	}

	return dto
}
