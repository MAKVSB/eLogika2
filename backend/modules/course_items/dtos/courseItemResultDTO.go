package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type CourseItemResultDTO struct {
	ID           uint   `json:"id"`
	StudentID    uint   `json:"studentId"`
	CourseItemID uint   `json:"courseItemId"`
	TermID       uint   `json:"termId"`
	UpdatedBy    string `json:"updatedBy"`

	TestInstanceID     *uint `json:"testInstanceId"`
	ActivityInstanceID *uint `json:"activityInstanceId"`

	Points   float64 `json:"points"`
	Final    bool    `json:"final"`
	Selected bool    `json:"selected"`
}

func (m CourseItemResultDTO) From(d *models.CourseItemResult) CourseItemResultDTO {
	dto := CourseItemResultDTO{
		ID:           d.ID,
		StudentID:    d.StudentID,
		CourseItemID: d.CourseItemID,
		TermID:       d.TermID,

		TestInstanceID:     d.TestInstanceID,
		ActivityInstanceID: d.ActivityInstanceID,

		Points:   d.Points,
		Final:    d.Final,
		Selected: d.Selected,
	}

	if d.UpdatedBy != nil {
		dto.UpdatedBy = d.UpdatedBy.FirstName + " " + d.UpdatedBy.FamilyName
	} else {
		dto.UpdatedBy = "System eLogika"
	}

	return dto
}
