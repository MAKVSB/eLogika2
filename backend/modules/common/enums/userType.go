package enums

type UserTypeEnum string

const ( // Special VS normal user from the context of whole system (not just single course)
	UserTypeNormal UserTypeEnum = "NORMAL"
	UserTypeAdmin  UserTypeEnum = "ADMIN"
)

var UserTypeEnumAll = []UserTypeEnum{
	UserTypeNormal,
	UserTypeAdmin,
}

func (w UserTypeEnum) TSName() string {
	switch w {
	case UserTypeNormal:
		return "NORMAL"
	case UserTypeAdmin:
		return "ADMIN"
	default:
		return "???"
	}
}
