package enums

type CourseUserRoleEnum string

const (
	CourseUserRoleAdmin CourseUserRoleEnum = "ADMIN"
	// CourseUserRoleSecretary CourseUserRoleEnum = "SECRETARY"
	CourseUserRoleGarant  CourseUserRoleEnum = "GARANT"
	CourseUserRoleTutor   CourseUserRoleEnum = "TUTOR"
	CourseUserRoleStudent CourseUserRoleEnum = "STUDENT"
)

var CourseUserRoleEnumAll = []CourseUserRoleEnum{
	CourseUserRoleAdmin,
	// CourseUserRoleSecretary,
	CourseUserRoleGarant,
	CourseUserRoleTutor,
	CourseUserRoleStudent,
}

func (w CourseUserRoleEnum) TSName() string {
	switch w {
	case CourseUserRoleAdmin:
		return "ADMIN"
	// case CourseUserRoleSecretary:
	// 	return "SECRETARY"
	case CourseUserRoleGarant:
		return "GARANT"
	case CourseUserRoleTutor:
		return "TUTOR"
	case CourseUserRoleStudent:
		return "STUDENT"
	default:
		return "???"
	}
}
