package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type StudentCourseItemDTO struct {
	ID          uint                     `json:"id"`
	Name        string                   `json:"name"`
	Type        enums.CourseItemTypeEnum `json:"type"`
	PointsMin   uint                     `json:"pointsMin"`
	PointsMax   uint                     `json:"pointsMax"`
	Mandatory   bool                     `json:"mandatory"`
	MaxAttempts uint                     `json:"maxAttempts"`

	Childs  []*StudentCourseItemDTO `json:"childs"`  // Needs to be assigned outside
	Results []*CourseItemResultDTO  `json:"results"` // Needs to be assigned outside
	Points  float64                 `json:"points"`  // Needs to be assigned outside
	Passed  bool                    `json:"passed"`  // Needs to be assigned outside
}

func (m StudentCourseItemDTO) From(d *models.CourseItem) StudentCourseItemDTO {
	dto := StudentCourseItemDTO{
		ID:          d.ID,
		Name:        d.Name,
		Type:        d.Type,
		PointsMin:   d.PointsMin,
		PointsMax:   d.PointsMax,
		Mandatory:   d.Mandatory,
		MaxAttempts: d.MaxAttempts,
		Points:      0,
	}

	return dto
}
