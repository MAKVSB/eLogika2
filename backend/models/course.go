package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type CourseImportOptions struct {
	Date string `json:"date"`
	Code string `json:"code"`
}

type Course struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt ``
	Version   uint           ``

	Name          string              ``
	Content       *TipTapContent      `gorm:"serializer:json;type:varbinary"`
	ContentFiles  []*File             `gorm:"many2many:course_content_files;"`
	Shortname     string              ``
	Public        bool                ``
	Year          uint                ``
	Semester      enums.SemesterEnum  ``
	ChapterID     *uint               ``
	PointsMin     float64             ``
	PointsMax     float64             ``
	ImportOptions CourseImportOptions `gorm:"serializer:json"`

	Terms []Term ``
}
