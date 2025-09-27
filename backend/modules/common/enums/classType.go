package enums

type ClassTypeEnum string

const (
	ClassTypeP ClassTypeEnum = "P"
	ClassTypeC ClassTypeEnum = "C"
	ClassTypeT ClassTypeEnum = "T"
)

var ClassTypeEnumAll = []ClassTypeEnum{
	ClassTypeP,
	ClassTypeC,
	ClassTypeT,
}

func (w ClassTypeEnum) TSName() string {
	switch w {
	case ClassTypeP:
		return "P"
	case ClassTypeC:
		return "C"
	case ClassTypeT:
		return "T"
	default:
		return "???"
	}
}
