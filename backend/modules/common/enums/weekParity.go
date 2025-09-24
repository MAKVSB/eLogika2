package enums

type WeekParityEnum string

const (
	WeekParityBoth WeekParityEnum = "BOTH"
	WeekParityEven WeekParityEnum = "EVEN"
	WeekParityOdd  WeekParityEnum = "ODD"
)

var WeekParityEnumAll = []WeekParityEnum{
	WeekParityBoth,
	WeekParityEven,
	WeekParityOdd,
}

func (w WeekParityEnum) TSName() string {
	switch w {
	case WeekParityBoth:
		return "BOTH"
	case WeekParityEven:
		return "EVEN"
	case WeekParityOdd:
		return "ODD"
	default:
		return "???"
	}
}
