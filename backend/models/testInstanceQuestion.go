package models

import "encoding/json"

type TestInstanceQuestion struct {
	CommonModel
	ID             uint
	TestInstanceID uint

	TestQuestionID         uint                         ``
	TextAnswer             json.RawMessage              ``
	TextAnswerReviewedByID *uint                        ``
	TextAnswerPercentage   uint                         ``
	Answers                []TestInstanceQuestionAnswer ``

	// TextAnswerReviewedBy *User         ``
	TestInstance *TestInstance ``
	TestQuestion *TestQuestion ``
}
