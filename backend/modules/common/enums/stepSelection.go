package enums

type StepSelectionEnum string

const (
	StepSelectionNC StepSelectionEnum = "NC"
	StepSelectionSC StepSelectionEnum = "SC"
	StepSelectionND StepSelectionEnum = "ND"
	StepSelectionSD StepSelectionEnum = "SD"
)

var StepSelectionEnumAll = []StepSelectionEnum{
	StepSelectionNC,
	StepSelectionSC,
	StepSelectionND,
	StepSelectionSD,
}

func (w StepSelectionEnum) TSName() string {
	switch w {

	case StepSelectionNC:
		return "NC"
	case StepSelectionSC:
		return "SC"
	case StepSelectionND:
		return "ND"
	case StepSelectionSD:
		return "SD"
	default:
		return "???"
	}
}
