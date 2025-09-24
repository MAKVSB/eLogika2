package enums

type QuestionFormatEnum string

const (
	QuestionFormatTest QuestionFormatEnum = "ABCD"
	QuestionFormatOpen QuestionFormatEnum = "OPEN"
)

var QuestionFormatEnumAll = []QuestionFormatEnum{
	QuestionFormatTest,
	QuestionFormatOpen,
}

func (w QuestionFormatEnum) TSName() string {
	switch w {
	case QuestionFormatTest:
		return "ABCD"
	case QuestionFormatOpen:
		return "OPEN"
	default:
		return "???"
	}
}
