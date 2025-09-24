package enums

type TestInstanceEventTypeEnum string

// ANSWERCHANGE, COPY, RESIZE, FOCUS, BLUR, RESIZE, TABVISIBLE, TABHIDDEN, CUT, PASTE, SHORTCUT and other types from frontend

const (
	TestInstanceEventTypePageload     TestInstanceEventTypeEnum = "PAGELOAD"
	TestInstanceEventTypeNetworkerror TestInstanceEventTypeEnum = "NETWORKERROR"
	TestInstanceEventTypeTeststart    TestInstanceEventTypeEnum = "TESTSTART"
	TestInstanceEventTypeClipboard    TestInstanceEventTypeEnum = "CLIPBOARD"
	TestInstanceEventTypeContextMenu  TestInstanceEventTypeEnum = "CONTEXTMENU"
	TestInstanceEventTypeSelectStart  TestInstanceEventTypeEnum = "SELECTSTART"
	TestInstanceEventTypeDragStart    TestInstanceEventTypeEnum = "DRAGSTART"
	TestInstanceEventTypeDrop         TestInstanceEventTypeEnum = "DROP"
	TestInstanceEventTypePrint        TestInstanceEventTypeEnum = "PRINT"
	TestInstanceEventTypeFullscreen   TestInstanceEventTypeEnum = "FULLSCREEN"
	TestInstanceEventTypeBlur         TestInstanceEventTypeEnum = "BLUR"
	TestInstanceEventTypeFocus        TestInstanceEventTypeEnum = "FOCUS"
	TestInstanceEventTypeResize       TestInstanceEventTypeEnum = "RESIZE"
	TestInstanceEventTypeHide         TestInstanceEventTypeEnum = "HIDE"
	TestInstanceEventTypeUnload       TestInstanceEventTypeEnum = "UNLOAD"
	TestInstanceEventTypeTabHidden    TestInstanceEventTypeEnum = "TABVISIBLE"
	TestInstanceEventTypeTabVisible   TestInstanceEventTypeEnum = "TABHIDDEN"
	TestInstanceEventTypePrintscreen  TestInstanceEventTypeEnum = "PRINTSCREEN"
	TestInstanceEventTypeShortcut     TestInstanceEventTypeEnum = "SHORTCUT"

	TestInstanceEventTypeQuestionUpdate      TestInstanceEventTypeEnum = "QUESTIONUPDATE"
	TestInstanceEventTypeQuestionSwitched    TestInstanceEventTypeEnum = "QUESTIONSWITCHED"
	TestInstanceEventTypeQuestionInvalidIP   TestInstanceEventTypeEnum = "INVALIDIP"
	TestInstanceEventTypeBonusPointsModified TestInstanceEventTypeEnum = "BONUSPOINTS"
)

var TestInstanceEventTypeEnumAll = []TestInstanceEventTypeEnum{
	TestInstanceEventTypePageload,
	TestInstanceEventTypeNetworkerror,
	TestInstanceEventTypeTeststart,
	TestInstanceEventTypeClipboard,
	TestInstanceEventTypeContextMenu,
	TestInstanceEventTypeSelectStart,
	TestInstanceEventTypeDragStart,
	TestInstanceEventTypeDrop,
	TestInstanceEventTypePrint,
	TestInstanceEventTypeFullscreen,
	TestInstanceEventTypeBlur,
	TestInstanceEventTypeFocus,
	TestInstanceEventTypeResize,
	TestInstanceEventTypeHide,
	TestInstanceEventTypeUnload,
	TestInstanceEventTypeTabHidden,
	TestInstanceEventTypeTabVisible,
	TestInstanceEventTypePrintscreen,
	TestInstanceEventTypeShortcut,

	TestInstanceEventTypeQuestionUpdate,
	TestInstanceEventTypeQuestionSwitched,
	TestInstanceEventTypeQuestionInvalidIP,
	TestInstanceEventTypeBonusPointsModified,
}

func (w TestInstanceEventTypeEnum) TSName() string {
	switch w {

	case TestInstanceEventTypePageload:
		return "PAGELOAD"
	case TestInstanceEventTypeNetworkerror:
		return "NETWORKERROR"
	case TestInstanceEventTypeTeststart:
		return "TESTSTART"
	case TestInstanceEventTypeClipboard:
		return "CLIPBOARD"
	case TestInstanceEventTypeContextMenu:
		return "CONTEXTMENU"
	case TestInstanceEventTypeSelectStart:
		return "SELECTSTART"
	case TestInstanceEventTypeDragStart:
		return "DRAGSTART"
	case TestInstanceEventTypeDrop:
		return "DROP"
	case TestInstanceEventTypePrint:
		return "PRINT"
	case TestInstanceEventTypeFullscreen:
		return "FULLSCREEN"
	case TestInstanceEventTypeBlur:
		return "BLUR"
	case TestInstanceEventTypeFocus:
		return "FOCUS"
	case TestInstanceEventTypeResize:
		return "RESIZE"
	case TestInstanceEventTypeHide:
		return "HIDE"
	case TestInstanceEventTypeUnload:
		return "UNLOAD"
	case TestInstanceEventTypeTabHidden:
		return "TABVISIBLE"
	case TestInstanceEventTypeTabVisible:
		return "TABHIDDEN"
	case TestInstanceEventTypePrintscreen:
		return "PRINTSCREEN"
	case TestInstanceEventTypeShortcut:
		return "SHORTCUT"

	case TestInstanceEventTypeQuestionUpdate:
		return "QUESTIONUPDATE"
	case TestInstanceEventTypeQuestionSwitched:
		return "QUESTIONSWITCHED"
	case TestInstanceEventTypeQuestionInvalidIP:
		return "INVALIDIP"
	case TestInstanceEventTypeBonusPointsModified:
		return "BONUSPOINTS"
	default:
		return "???"
	}
}
