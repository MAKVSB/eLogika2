package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type StudentCourseItemDTO struct {
	ID          uint                     `json:"id"`
	Version     uint                     `json:"version"`
	Name        string                   `json:"name"`
	Type        enums.CourseItemTypeEnum `json:"type"`
	PointsMin   uint                     `json:"pointsMin"`
	PointsMax   uint                     `json:"pointsMax"`
	Mandatory   bool                     `json:"mandatory"`
	StudyForm   enums.StudyFormEnum      `json:"studyForm"`
	MaxAttempts uint                     `json:"maxAttempts"`
	ParentID    *uint                    `json:"parentId"`

	Childs []StudentCourseItemDTO `json:"childs"` // Needs to be assigned outside
	Points float64                `json:"points"` // Needs to be assigned outside
	Passed bool                   `json:"passed"` // Needs to be assigned outside
}

func (m StudentCourseItemDTO) From(d *models.CourseItem) StudentCourseItemDTO {
	dto := StudentCourseItemDTO{
		ID:      d.ID,
		Version: d.Version,

		Name:        d.Name,
		Type:        d.Type,
		PointsMin:   d.PointsMin,
		PointsMax:   d.PointsMax,
		Mandatory:   d.Mandatory,
		StudyForm:   d.StudyForm,
		MaxAttempts: d.MaxAttempts,
		ParentID:    d.ParentID,
		Points:      0,
	}

	return dto
}
