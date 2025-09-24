package models

import (
	"time"

	"gorm.io/gorm"
)

type Test struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	DeletedAt gorm.DeletedAt ``

	CourseID     uint ``
	CourseItemID uint ``
	TermID       uint ``
	CreatedByID  uint ``

	Name  string ``
	Group string ``

	Course     *Course        ``
	CourseItem *CourseItem    ``
	Term       *Term          ``
	CreatedBy  *User          ``
	Blocks     []TestBlock    `gorm:"serializer:json"`
	Questions  []TestQuestion ``
	Instances  []TestInstance ``
}

type TestBlock struct {
	ID                    uint   `json:"id"`
	Title                 string `json:"title"`
	ShowName              bool   `json:"showName"`
	Weight                uint   `json:"weight"`
	WrongAnswerPercentage uint   `json:"wrongAnswerPercentage"`
}
