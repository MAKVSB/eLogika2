package enums

type TestInstanceEventSourceEnum string

const (
	TestInstanceEventSourceServer TestInstanceEventSourceEnum = "SERVER"
	TestInstanceEventSourceClient TestInstanceEventSourceEnum = "CLIENT"
	TestInstanceEventSourceSystem TestInstanceEventSourceEnum = "SYSTEM"
)

var TestInstanceEventSourceEnumAll = []TestInstanceEventSourceEnum{
	TestInstanceEventSourceServer,
	TestInstanceEventSourceClient,
	TestInstanceEventSourceSystem,
}

func (w TestInstanceEventSourceEnum) TSName() string {
	switch w {

	case TestInstanceEventSourceServer:
		return "SERVER"
	case TestInstanceEventSourceClient:
		return "CLIENT"
	case TestInstanceEventSourceSystem:
		return "SYSTEM"
	default:
		return "???"
	}
}
