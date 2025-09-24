package enums

type StudyFormEnum string

const (
	StudyFormFulltime StudyFormEnum = "FULLTIME"
	StudyFormCombined StudyFormEnum = "COMBINED"
)

var StudyFormEnumAll = []StudyFormEnum{
	StudyFormFulltime,
	StudyFormCombined,
}

func (w StudyFormEnum) TSName() string {
	switch w {
	case StudyFormFulltime:
		return "FULLTIME"
	case StudyFormCombined:
		return "COMBINED"
	default:
		return "???"
	}
}
