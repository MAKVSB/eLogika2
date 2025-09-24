package enums

type CategoryFilterEnum string

const (
	CategoryFilterALL  CategoryFilterEnum = "ALL"
	CategoryFilterS    CategoryFilterEnum = "S"
	CategoryFilterQ    CategoryFilterEnum = "Q"
	CategoryFilterSQOR CategoryFilterEnum = "SQOR"
)

var CategoryFilterEnumAll = []CategoryFilterEnum{
	CategoryFilterALL,
	CategoryFilterS,
	CategoryFilterQ,
	CategoryFilterSQOR,
}

func (w CategoryFilterEnum) TSName() string {
	switch w {
	case CategoryFilterALL:
		return "ALL"
	case CategoryFilterS:
		return "S"
	case CategoryFilterQ:
		return "Q"
	case CategoryFilterSQOR:
		return "SQOR"
	default:
		return "???"
	}
}
