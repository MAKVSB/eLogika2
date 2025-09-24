package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type CourseUserDTO struct {
	ID         uint                       `json:"id"`
	FirstName  string                     `json:"firstName"`
	FamilyName string                     `json:"familyName"`
	Username   string                     `json:"username"`
	Email      string                     `json:"email"`
	Roles      []enums.CourseUserRoleEnum `json:"roles"`
	StudyForm  *enums.StudyFormEnum       `json:"studyForm"`
}

func (m CourseUserDTO) From(d *models.CourseUser) CourseUserDTO {
	dto := CourseUserDTO{
		ID:         d.User.ID,
		FirstName:  d.User.FirstName,
		FamilyName: d.User.FamilyName,
		Username:   d.User.Username,
		Email:      d.User.Email,
		Roles:      d.Roles,
		StudyForm:  d.StudyForm,
	}

	return dto
}
