package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type CourseItemResultsDTO struct {
	ID           uint    `json:"id"`
	Username     string  `json:"username"`
	DegreeBefore string  `json:"degreeBefore"`
	FirstName    string  `json:"firstName"`
	FamilyName   string  `json:"familyName"`
	DegreeAfter  string  `json:"degreeAfter"`
	Email        string  `json:"email"`
	Points       float64 `json:"points"` // External
	Passed       bool    `json:"passed"` // External

	Results []*CourseItemResultDTO `json:"results"`
}

func (m CourseItemResultsDTO) From(d *models.User) CourseItemResultsDTO {
	dto := CourseItemResultsDTO{
		ID:           d.ID,
		DegreeBefore: d.DegreeBefore,
		FirstName:    d.FirstName,
		FamilyName:   d.FamilyName,
		DegreeAfter:  d.DegreeAfter,
		Username:     d.Username,
		Email:        d.Email,

		Results: make([]*CourseItemResultDTO, len(d.Results)),
	}

	for i, res := range d.Results {
		resDto := CourseItemResultDTO{}.From(res)
		dto.Results[i] = &resDto
	}

	return dto
}
