package dtos

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type CourseItemDTO struct {
	ID      uint `json:"id"`
	Version uint `json:"version"`

	Name              string                      `json:"name"`
	Type              enums.CourseItemTypeEnum    `json:"type"`
	PointsMin         uint                        `json:"pointsMin"`
	PointsMax         uint                        `json:"pointsMax"`
	Mandatory         bool                        `json:"mandatory"`
	StudyForm         enums.StudyFormEnum         `json:"studyForm"`
	MaxAttempts       uint                        `json:"maxAttempts"`
	AllowNegative     bool                        `json:"allowNegative"`
	Editable          bool                        `json:"editable"`
	ManagedBy         enums.CourseUserRoleEnum    `json:"managedBy"`
	EvaluateByAttempt enums.EvaluateByAttemptEnum `json:"evaluateByAttempt"`

	ActivityDetail *CourseItemActivityDTO `json:"activityDetail,omitempty"`
	TestDetail     *CourseItemTestDTO     `json:"testDetail,omitempty"`
	GroupDetail    *CourseItemGroupDTO    `json:"groupDetail,omitempty"`
}

func (m CourseItemDTO) From(d *models.CourseItem) CourseItemDTO {
	dto := CourseItemDTO{
		ID:      d.ID,
		Version: d.Version,

		Name:              d.Name,
		Type:              d.Type,
		PointsMin:         d.PointsMin,
		PointsMax:         d.PointsMax,
		Mandatory:         d.Mandatory,
		StudyForm:         d.StudyForm,
		MaxAttempts:       d.MaxAttempts,
		AllowNegative:     d.AllowNegative,
		Editable:          d.Editable,
		ManagedBy:         d.ManagedBy,
		EvaluateByAttempt: d.EvaluateByAttempt,
	}
	if d.ActivityDetail != nil {
		a := CourseItemActivityDTO{}.From(d.ActivityDetail)
		dto.ActivityDetail = &a
	}

	if d.TestDetail != nil {
		a := CourseItemTestDTO{}.From(d.TestDetail)
		dto.TestDetail = &a
	}

	if d.GroupDetail != nil {
		a := CourseItemGroupDTO{}.From(d.GroupDetail)
		dto.GroupDetail = &a
	}

	return dto
}
