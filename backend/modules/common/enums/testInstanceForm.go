package enums

type TestInstanceFormEnum string

const (
	TestInstanceFormOnline  TestInstanceFormEnum = "ONLINE"
	TestInstanceFormOffline TestInstanceFormEnum = "OFFLINE"
)

var TestInstanceFormEnumAll = []TestInstanceFormEnum{
	TestInstanceFormOnline,
	TestInstanceFormOffline,
}

func (w TestInstanceFormEnum) TSName() string {
	switch w {
	case TestInstanceFormOnline:
		return "ONLINE"
	case TestInstanceFormOffline:
		return "OFFLINE"
	default:
		return "???"
	}
}
