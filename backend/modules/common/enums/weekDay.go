package enums

type WeekDayEnum string

const ( // Special VS normal user from the context of whole system (not just single course)
	WeekDayMonday    WeekDayEnum = "MONDAY"
	WeekDayTuesday   WeekDayEnum = "TUESDAY"
	WeekDayWednesday WeekDayEnum = "WEDNESDAY"
	WeekDayThursday  WeekDayEnum = "THURSDAY"
	WeekDayFriday    WeekDayEnum = "FRIDAY"
	WeekDaySaturday  WeekDayEnum = "SATURDAY"
	WeekDaySunday    WeekDayEnum = "SUNDAY"
)

var WeekDayEnumAll = []WeekDayEnum{
	WeekDayMonday,
	WeekDayTuesday,
	WeekDayWednesday,
	WeekDayThursday,
	WeekDayFriday,
	WeekDaySaturday,
	WeekDaySunday,
}

func (w WeekDayEnum) TSName() string {
	switch w {
	case WeekDayMonday:
		return "MONDAY"
	case WeekDayTuesday:
		return "TUESDAY"
	case WeekDayWednesday:
		return "WEDNESDAY"
	case WeekDayThursday:
		return "THURSDAY"
	case WeekDayFriday:
		return "FRIDAY"
	case WeekDaySaturday:
		return "SATURDAY"
	case WeekDaySunday:
		return "SUNDAY"
	default:
		return "???"
	}
}
