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
	Passed     bool    `json:"passed"` // External

	Results []*CourseItemResultDTO `json:"results"`
}

func (m CourseItemResultsDTO) From(d *models.User) CourseItemResultsDTO {
	dto := CourseItemResultsDTO{
		ID:         d.ID,
		FirstName:  d.FirstName,
		FamilyName: d.FamilyName,
		Username:   d.Username,
		Email:      d.Email,

		Results: make([]*CourseItemResultDTO, len(d.Results)),
	}

	for i, res := range d.Results {
		resDto := CourseItemResultDTO{}.From(res)
		dto.Results[i] = &resDto
	}

	return dto
}
