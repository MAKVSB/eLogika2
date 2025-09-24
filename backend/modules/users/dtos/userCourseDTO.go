package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type UserCourseDTO struct {
	ID       uint                       `json:"id"`
	Name     string                     `json:"name"`
	Year     uint                       `json:"year"`
	Semester enums.SemesterEnum         `json:"semester"`
	Roles    []enums.CourseUserRoleEnum `json:"roles"`
}

func (m UserCourseDTO) From(d *models.CourseUser) UserCourseDTO {
	dto := UserCourseDTO{
		ID:       d.Course.ID,
		Name:     d.Course.Name,
		Roles:    d.Roles,
		Year:     d.Course.Year,
		Semester: d.Course.Semester,
	}

	return dto
}
