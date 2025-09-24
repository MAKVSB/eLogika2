package models

import (
	"encoding/json"
	"time"

	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type TestInstanceEvent struct {
	CommonModel
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      ``
	DeletedAt gorm.DeletedAt ``

	TestInstanceID uint ``
	UserID         uint ``

	OccuredAt   time.Time                         ``
	EventSource enums.TestInstanceEventSourceEnum ``
	EventType   enums.TestInstanceEventTypeEnum   ``
	EventData   json.RawMessage                   ``
	PageID      string                            ``

	TestInstance TestInstance ``
	User         User         ``
}
