package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type LoggedUserCourseDTO struct {
	ID    uint                       `json:"id"`
	Roles []enums.CourseUserRoleEnum `json:"roles"`
}

func (m LoggedUserCourseDTO) From(d *models.CourseUser) LoggedUserCourseDTO {
	adminCourseDTO := LoggedUserCourseDTO{
		ID:    d.Course.ID,
		Roles: d.Roles,
	}

	return adminCourseDTO
}

func (m LoggedUserCourseDTO) HasRole(role enums.CourseUserRoleEnum) bool {
	for _, r := range m.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (m LoggedUserCourseDTO) IsStudent() bool {
	return m.HasRole(enums.CourseUserRoleStudent)
}

func (m LoggedUserCourseDTO) IsTutor() bool {
	return m.HasRole(enums.CourseUserRoleTutor)
}

func (m LoggedUserCourseDTO) IsGarant() bool {
	return m.HasRole(enums.CourseUserRoleGarant)
}

// func (m LoggedUserCourseDTO) IsSecretary() bool {
// 	return m.HasRole(enums.CourseUserRoleSecretary)
// }

func (m LoggedUserCourseDTO) IsAdmin() bool {
	return m.HasRole(enums.CourseUserRoleAdmin)
}
