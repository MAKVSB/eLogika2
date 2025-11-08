package dtos

import (
	"elogika.vsb.cz/backend/models"
)

type ActivityInstanceDTO struct {
	ID      uint                  `json:"id"`
	Content *models.TipTapContent `json:"content" ts_type:"JSONContent"`

	AssignmentName           string                `json:"assignmentName"`
	AssignmentDescription    *models.TipTapContent `json:"assignmentDescription" ts_type:"JSONContent"`
	AssignmentExpectedResult *models.TipTapContent `json:"assignmentExpectedResult" ts_type:"JSONContent"`

	Points    float64 `json:"points"`
	PointsMin uint    `json:"pointsMin"`
	PointsMax uint    `json:"pointsMax"`

	Editable bool `json:"editable"`
}

func (m ActivityInstanceDTO) From(d *models.ActivityInstance, editable bool) ActivityInstanceDTO {
	dto := ActivityInstanceDTO{
		ID:                       d.ID,
		Content:                  d.Content,
		AssignmentName:           d.CourseItem.Name,
		AssignmentDescription:    d.CourseItem.ActivityDetail.Description,
		AssignmentExpectedResult: d.CourseItem.ActivityDetail.ExpectedResult,
		Points:                   d.Result.Points,
		PointsMin:                d.CourseItem.PointsMin,
		PointsMax:                d.CourseItem.PointsMax,
		Editable:                 editable,
	}

	return dto
}
