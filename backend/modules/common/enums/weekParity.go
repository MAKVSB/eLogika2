package enums

type WeekParityEnum string

const (
	WeekParityBoth   WeekParityEnum = "BOTH"
	WeekParityEven   WeekParityEnum = "EVEN"
	WeekParityOdd    WeekParityEnum = "ODD"
	WeekParityCustom WeekParityEnum = "CUSTOM"
)

var WeekParityEnumAll = []WeekParityEnum{
	WeekParityBoth,
	WeekParityEven,
	WeekParityOdd,
	WeekParityCustom,
}

func (w WeekParityEnum) TSName() string {
	switch w {
	case WeekParityBoth:
		return "BOTH"
	case WeekParityEven:
		return "EVEN"
	case WeekParityOdd:
		return "ODD"
	case WeekParityCustom:
		return "CUSTOM"
	default:
		return "???"
	}
}
