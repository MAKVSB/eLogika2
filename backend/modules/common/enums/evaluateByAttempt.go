package enums

type EvaluateByAttemptEnum string

const (
	EvaluateByAttemptBest EvaluateByAttemptEnum = "BEST"
	EvaluateByAttemptLast EvaluateByAttemptEnum = "LAST"
)

var EvaluateByAttemptEnumAll = []EvaluateByAttemptEnum{
	EvaluateByAttemptBest,
	EvaluateByAttemptLast,
}

func (w EvaluateByAttemptEnum) TSName() string {
	switch w {
	case EvaluateByAttemptBest:
		return "BEST"
	case EvaluateByAttemptLast:
		return "LAST"
	default:
		return "???"
	}
}
