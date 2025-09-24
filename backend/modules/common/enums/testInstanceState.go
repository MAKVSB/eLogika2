package enums

type TestInstanceStateEnum string

const (
	TestInstanceStateReady    TestInstanceStateEnum = "READY"
	TestInstanceStateActive   TestInstanceStateEnum = "ACTIVE"
	TestInstanceStateFinished TestInstanceStateEnum = "FINISHED"
	TestInstanceStateExpired  TestInstanceStateEnum = "EXPIRED"
)

var TestInstanceStateEnumAll = []TestInstanceStateEnum{
	TestInstanceStateReady,
	TestInstanceStateActive,
	TestInstanceStateFinished,
	TestInstanceStateExpired,
}

func (w TestInstanceStateEnum) TSName() string {
	switch w {
	case TestInstanceStateReady:
		return "READY"
	case TestInstanceStateActive:
		return "ACTIVE"
	case TestInstanceStateFinished:
		return "FINISHED"
	case TestInstanceStateExpired:
		return "EXPIRED"
	default:
		return "???"
	}
}
