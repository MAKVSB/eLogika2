package enums

type QuestionTypeEnum string

const (
	QuestionTypeExam     QuestionTypeEnum = "EXAM"
	QuestionTypePractise QuestionTypeEnum = "PRACTISE"
	QuestionTypeSample   QuestionTypeEnum = "SAMPLE"
)

var QuestionTypeEnumAll = []QuestionTypeEnum{
	QuestionTypeExam,
	QuestionTypePractise,
	QuestionTypeSample,
}

func (w QuestionTypeEnum) TSName() string {
	switch w {
	case QuestionTypeExam:
		return "EXAM"
	case QuestionTypePractise:
		return "PRACTISE"
	case QuestionTypeSample:
		return "SAMPLE"
	default:
		return "???"
	}
}
