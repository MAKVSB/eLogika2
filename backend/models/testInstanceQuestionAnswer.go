package models

type TestInstanceQuestionAnswer struct {
	CommonModel
	ID                     uint ``
	TestInstanceQuestionID uint ``
	TestQuestionAnswerID   uint ``

	Selected bool ``

	TestQuestionAnswer *TestQuestionAnswer ``
}
