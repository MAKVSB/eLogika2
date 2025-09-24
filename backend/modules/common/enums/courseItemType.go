package enums

type CourseItemTypeEnum string

const (
	CourseItemTypeGroup    CourseItemTypeEnum = "GROUP"
	CourseItemTypeTest     CourseItemTypeEnum = "TEST"
	CourseItemTypeActivity CourseItemTypeEnum = "ACTIVITY"
)

var CourseItemTypeEnumAll = []CourseItemTypeEnum{
	CourseItemTypeGroup,
	CourseItemTypeTest,
	CourseItemTypeActivity,
}

func (w CourseItemTypeEnum) TSName() string {
	switch w {
	case CourseItemTypeGroup:
		return "GROUP"
	case CourseItemTypeTest:
		return "TEST"
	case CourseItemTypeActivity:
		return "ACTIVITY"
	default:
		return "???"
	}
}
