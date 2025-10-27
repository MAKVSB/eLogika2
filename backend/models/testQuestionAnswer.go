package models

type TestQuestionAnswer struct {
	CommonModel
	ID uint

	TestQuestionID uint ``
	AnswerID       uint ``
	Order          uint ``

	Answer *Answer ``
}
