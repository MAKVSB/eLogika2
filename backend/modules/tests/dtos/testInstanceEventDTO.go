package dtos

import (
	"encoding/json"
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type TestInstanceEventDTO struct {
	ID          uint                              `json:"id"`
	OccuredAt   time.Time                         `json:"occuredAt"`
	EventSource enums.TestInstanceEventSourceEnum `json:"eventSource"`
	EventType   enums.TestInstanceEventTypeEnum   `json:"eventType"`
	EventData   json.RawMessage                   `json:"eventData"`
	PageID      string                            `json:"pageId"`

	User TestParticipantDTO `json:"user"`
}

func (m TestInstanceEventDTO) From(d *models.TestInstanceEvent) TestInstanceEventDTO {
	dto := TestInstanceEventDTO{
		ID:          d.ID,
		OccuredAt:   d.OccuredAt,
		EventSource: d.EventSource,
		EventType:   d.EventType,
		EventData:   d.EventData,
		PageID:      d.PageID,

		User: TestParticipantDTO{}.From(&d.User),
	}

	return dto
}
