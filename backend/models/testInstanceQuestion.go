package models

import "encoding/json"

type TestInstanceQuestion struct {
	CommonModel
	ID             uint
	TestInstanceID uint

	TestQuestionID         uint                         ``
	TextAnswer             json.RawMessage              ``
	TextAnswerReviewedByID *uint                        ``
	TextAnswerPercentage   float64                      ``
	Answers                []TestInstanceQuestionAnswer ``

	// TextAnswerReviewedBy *User         ``
	TestInstance *TestInstance ``
	TestQuestion *TestQuestion ``
}
