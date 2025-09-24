package enums

type SemesterEnum string

const (
	SemesterSummer SemesterEnum = "SUMMER"
	SemesterWinter SemesterEnum = "WINTER"
)

var SemesterEnumAll = []SemesterEnum{
	SemesterSummer,
	SemesterWinter,
}

func (w SemesterEnum) TSName() string {
	switch w {
	case SemesterSummer:
		return "SUMMER"
	case SemesterWinter:
		return "WINTER"
	default:
		return "???"
	}
}
