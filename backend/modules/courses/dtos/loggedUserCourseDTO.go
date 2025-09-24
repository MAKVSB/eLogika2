package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type LoggedUserCourseDTO2 struct {
	ID uint `json:"id"`

	Name      string                     `json:"name"`
	Shortname string                     `json:"shortName"`
	Public    bool                       `json:"public"`
	Roles     []enums.CourseUserRoleEnum `json:"roles"`
	Year      uint                       `json:"year"`
	Semester  enums.SemesterEnum         `json:"semester"`
	ChapterID uint                       `json:"chapterId"`
}

func (m LoggedUserCourseDTO2) From(d *models.CourseUser) LoggedUserCourseDTO2 {
	adminCourseDTO := LoggedUserCourseDTO2{
		ID:        d.Course.ID,
		Name:      d.Course.Name,
		Shortname: d.Course.Shortname,
		Public:    d.Course.Public,
		Roles:     d.Roles,
		Year:      d.Course.Year,
		Semester:  d.Course.Semester,
	}

	if d.Course.ChapterID != nil {
		adminCourseDTO.ChapterID = *d.Course.ChapterID
	}

	return adminCourseDTO
}

func (m LoggedUserCourseDTO2) IsStudent() bool {
	for _, r := range m.Roles {
		if r == enums.CourseUserRoleStudent {
			return true
		}
	}
	return false
}

func (m LoggedUserCourseDTO2) IsTutor() bool {
	for _, r := range m.Roles {
		if r == enums.CourseUserRoleTutor {
			return true
		}
	}
	return false
}

func (m LoggedUserCourseDTO2) IsGarant() bool {
	for _, r := range m.Roles {
		if r == enums.CourseUserRoleGarant {
			return true
		}
	}
	return false
}

// func (m LoggedUserCourseDTO2) IsSecretary() bool {
// 	for _, r := range m.Roles {
// 		if r == enums.CourseUserRoleSecretary {
// 			return true
// 		}
// 	}
// 	return false
// }

func (m LoggedUserCourseDTO2) IsAdmin() bool {
	for _, r := range m.Roles {
		if r == enums.CourseUserRoleAdmin {
			return true
		}
	}
	return false
}
