package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type ClassImportOptions struct {
	Code string `json:"code"`
}

type Class struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	CourseID      uint                 ``
	Name          string               ``
	Room          string               ``
	Type          enums.ClassTypeEnum  ``
	StudyForm     enums.StudyFormEnum  ``
	TimeFrom      string               ``
	TimeTo        string               ``
	Day           enums.WeekDayEnum    ``
	WeekParity    enums.WeekParityEnum ``
	StudentLimit  uint                 ``
	ImportOptions ClassImportOptions   `gorm:"serializer:json"`

	Course   *Course        ``
	Students []ClassStudent ``
	Tutors   []ClassTutor   ``
}
