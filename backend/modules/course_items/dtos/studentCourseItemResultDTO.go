package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
)

type StudentCourseItemResultDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`

	CourseItemName      string `json:"courseItemName"`
	CourseItemGroupName string `json:"courseItemGroupName"`

	TestInstanceID     *uint     `json:"testInstanceId"`
	ActivityInstanceID *uint     `json:"activityInstanceId"`
	InstanceStartTime  time.Time `json:"instanceStartTime"`

	Points   float64 `json:"points"`
	Final    bool    `json:"final"`
	Selected bool    `json:"selected"`

	TermName       string    `json:"termName"`
	TermActiveFrom time.Time `json:"termActiveFrom"`
	TermActiveTo   time.Time `json:"termActiveTo"`
}

func (m StudentCourseItemResultDTO) From(d *models.CourseItemResult) StudentCourseItemResultDTO {
	dto := StudentCourseItemResultDTO{
		ID:                 d.ID,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
		TestInstanceID:     d.TestInstanceID,
		ActivityInstanceID: d.ActivityInstanceID,
		InstanceStartTime:  d.CreatedAt,
		Points:             d.Points,
		Final:              d.Final,
		Selected:           d.Selected,
		TermName:           d.Term.Name,
		TermActiveFrom:     d.Term.ActiveFrom,
		TermActiveTo:       d.Term.ActiveTo,
		CourseItemName:     d.CourseItem.Name,
	}

	if d.CourseItem.Parent != nil {
		dto.CourseItemGroupName = d.CourseItem.Parent.Name
	}

	if d.UpdatedBy != nil {
		dto.UpdatedBy = d.UpdatedBy.FirstName + " " + d.UpdatedBy.FamilyName
	} else {
		dto.UpdatedBy = "System eLogika"
	}

	return dto
}
