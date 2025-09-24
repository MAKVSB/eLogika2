package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type LoggedUserDTO struct {
	ID         uint                  `json:"id"`
	FirstName  string                `json:"firstName"`
	FamilyName string                `json:"familyName"`
	Username   string                `json:"username"`
	Email      string                `json:"email"`
	Courses    []LoggedUserCourseDTO `json:"courses"`
	Type       enums.UserTypeEnum    `json:"type"`
}

func (m LoggedUserDTO) From(d *models.User) LoggedUserDTO {
	loggedUser := LoggedUserDTO{
		ID:         d.ID,
		FirstName:  d.FirstName,
		FamilyName: d.FamilyName,
		Username:   d.Username,
		Email:      d.Email,
		Courses:    make([]LoggedUserCourseDTO, len(d.UserCourses)),
		Type:       d.Type,
	}

	for i, course := range d.UserCourses {
		loggedUser.Courses[i] = LoggedUserCourseDTO{}.From(course)
	}

	return loggedUser
}

func (lu LoggedUserDTO) GetCourse(courseId uint) *LoggedUserCourseDTO {
	for _, item := range lu.Courses {
		if item.ID == courseId {
			return &item
		}
	}
	return nil
}
