package enums

type AnswerDistributionEnum string

const (
	AnswerDistributionExactlyOneCorrect             AnswerDistributionEnum = "EXACTLY_ONE_CORRECT"               // Exactly 1 correct
	AnswerDistributionMinimumOneCorrect             AnswerDistributionEnum = "MINIMUM_ONE_CORRECT"               // At least 1 correct
	AnswerDistributionMinimumOneCorrectOneIncorrect AnswerDistributionEnum = "MINIMUM_ONE_CORRECT_ONE_INCORRECT" // At least 1 correct and 1 incorrect
)

var AnswerDistributionEnumAll = []AnswerDistributionEnum{
	AnswerDistributionExactlyOneCorrect,
	AnswerDistributionMinimumOneCorrect,
	AnswerDistributionMinimumOneCorrectOneIncorrect,
}

func (w AnswerDistributionEnum) TSName() string {
	switch w {
	case AnswerDistributionExactlyOneCorrect:
		return "EXACTLY_ONE_CORRECT"
	case AnswerDistributionMinimumOneCorrect:
		return "MINIMUM_ONE_CORRECT"
	case AnswerDistributionMinimumOneCorrectOneIncorrect:
		return "MINIMUM_ONE_CORRECT_ONE_INCORRECT"
	default:
		return "???"
	}
}
