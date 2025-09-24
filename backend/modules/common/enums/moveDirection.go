package enums

type MoveDirectionEnum string

const (
	MoveDirectionUp   MoveDirectionEnum = "UP"
	MoveDirectionDown MoveDirectionEnum = "DOWN"
)

var MoveDirectionEnumAll = []MoveDirectionEnum{
	MoveDirectionUp,
	MoveDirectionDown,
}

func (w MoveDirectionEnum) TSName() string {
	switch w {
	case MoveDirectionUp:
		return "UP"
	case MoveDirectionDown:
		return "DOWN"
	default:
		return "???"
	}
}
