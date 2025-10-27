package models

type TestQuestion struct {
	CommonModel
	ID      uint
	TestID  uint
	BlockID uint

	QuestionID uint ``
	Order      uint ``

	Question *Question             ``
	Answers  []*TestQuestionAnswer ``

	// Temporary helper data
	OpenAnswers   []Answer `gorm:"-"`
	QuestionCount int      `gorm:"-"`
}
