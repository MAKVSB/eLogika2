package dtos

import "elogika.vsb.cz/backend/models"

type JoinedStudentDTO struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"firstName"`
	FamilyName string `json:"familyName"`
	Username   string `json:"username"`
	Email      string `json:"email"`
}

func (m JoinedStudentDTO) From(d *models.UserTerm) JoinedStudentDTO {
	dto := JoinedStudentDTO{
		ID:         d.User.ID,
		FirstName:  d.User.FirstName,
		FamilyName: d.User.FamilyName,
		Username:   d.User.Username,
		Email:      d.User.Email,
	}

	return dto
}
