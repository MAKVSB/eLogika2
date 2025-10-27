package models

import (
	"time"

	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type TestInstance struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	DeletedAt gorm.DeletedAt ``

	State     enums.TestInstanceStateEnum ``
	Form      enums.TestInstanceFormEnum  ``
	StartedAt time.Time                   `` // Time that the test started at
	EndsAt    time.Time                   `` // Time until the student must end the instance
	EndedAt   time.Time                   `` // Time that the intance actually ended (Student finishes test or automatic job marks as finished)

	TestID        uint ``
	ParticipantID uint ``
	TermID        uint ``
	CourseItemID  uint ``

	Test        *Test       ``
	Participant *User       ``
	Term        *Term       ``
	CourseItem  *CourseItem ``

	Questions         []*TestInstanceQuestion ``
	Result            *CourseItemResult       ``
	BonusPoints       float64                 ``
	BonusPointsReason string                  ``
	RecognizerFiles   []*RecognizerFile       ``
}

func (ti TestInstance) IsExpired(timeNow time.Time) bool {
	return ti.EndedAt.Add(time.Minute*1).Compare(timeNow) > 0
}
