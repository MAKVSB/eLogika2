package enums

type ClassTypeEnum string

const (
	ClassTypeP ClassTypeEnum = "P"
	ClassTypeC ClassTypeEnum = "C"
)

var ClassTypeEnumAll = []ClassTypeEnum{
	ClassTypeP,
	ClassTypeC,
}

func (w ClassTypeEnum) TSName() string {
	switch w {
	case ClassTypeP:
		return "P"
	case ClassTypeC:
		return "C"
	default:
		return "???"
	}
}
