package dtos

import "time"

type StudentTestDTO struct {
	TermID         uint      `json:"termId"`
	TermName       string    `json:"termName"`
	CourseItemId   uint      `json:"courseItemId"`
	CourseItemName string    `json:"courseItemName"`
	TriesLeft      uint      `json:"triesLeft"`
	ActiveFrom     time.Time `json:"activeFrom"`
	ActiveTo       time.Time `json:"activeTo"`
	CanStart       bool      `json:"canStart"`
}
