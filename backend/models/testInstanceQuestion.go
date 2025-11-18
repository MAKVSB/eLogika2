package models

type TestInstanceQuestion struct {
	CommonModel
	ID             uint `gorm:"primarykey"`
	TestInstanceID uint

	TestQuestionID         uint                          ``
	TextAnswer             *TipTapContent                `gorm:"serializer:json;type:varbinary(max)"`
	TextAnswerReviewedByID *uint                         ``
	TextAnswerPercentage   float64                       ``
	Answers                []*TestInstanceQuestionAnswer ``

	// TextAnswerReviewedBy *User         ``
	TestInstance *TestInstance ``
	TestQuestion *TestQuestion ``
}
