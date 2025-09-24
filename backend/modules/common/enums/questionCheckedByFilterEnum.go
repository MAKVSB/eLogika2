package enums

type QuestionCheckedByFilterEnum string

const (
	QuestionCheckedByFilterChecked       QuestionCheckedByFilterEnum = "CHECKED"
	QuestionCheckedByFilterCheckedByMe   QuestionCheckedByFilterEnum = "CHECKEDBYME"
	QuestionCheckedByFilterUnchecked     QuestionCheckedByFilterEnum = "UNCHECKED"
	QuestionCheckedByFilterUncheckedByMe QuestionCheckedByFilterEnum = "UNCHECKEDBYME"
)

var QuestionCheckedByFilterEnumAll = []QuestionCheckedByFilterEnum{
	QuestionCheckedByFilterChecked,
	QuestionCheckedByFilterCheckedByMe,
	QuestionCheckedByFilterUnchecked,
	QuestionCheckedByFilterUncheckedByMe,
}

func (w QuestionCheckedByFilterEnum) TSName() string {
	switch w {
	case QuestionCheckedByFilterChecked:
		return "CHECKED"
	case QuestionCheckedByFilterCheckedByMe:
		return "CHECKEDBYME"
	case QuestionCheckedByFilterUnchecked:
		return "UNCHECKED"
	case QuestionCheckedByFilterUncheckedByMe:
		return "UNCHECKEDBYME"
	default:
		return "???"
	}
}
