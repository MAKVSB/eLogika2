package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type CourseItemResultsDTO struct {
	ID         uint    `json:"id"`
	Username   string  `json:"username"`
	FirstName  string  `json:"firstName"`
	FamilyName string  `json:"familyName"`
	Email      string  `json:"email"`
	Points     float64 `json:"points"` // External
	Final      bool    `json:"final"`  // External
	Passed     bool    `json:"passed"` // External

	Courses []CourseItemResultDTO `json:"courses"`
}

func (m CourseItemResultsDTO) From(d *models.User) CourseItemResultsDTO {
	dto := CourseItemResultsDTO{
		ID:         d.ID,
		FirstName:  d.FirstName,
		FamilyName: d.FamilyName,
		Username:   d.Username,
		Email:      d.Email,

		Courses: make([]CourseItemResultDTO, len(d.Results)),
	}

	for i, res := range d.Results {
		dto.Courses[i] = CourseItemResultDTO{}.From(res)
	}

	return dto
}
