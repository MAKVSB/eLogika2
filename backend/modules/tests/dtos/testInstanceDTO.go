package dtos

import (
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common/enums"
)

type TestInstanceDTO struct {
	ID uint `json:"id"`

	State     enums.TestInstanceStateEnum `json:"state"`
	StartedAt time.Time                   `json:"startedAt"`
	EndedAt   time.Time                   `json:"endedAt"`
	EndsAt    time.Time                   `json:"endsAt"`

	Group string `json:"group"`

	Questions     []TestInstanceQuestionDTO `json:"questions,omitempty"`
	QuestionCount uint                      `json:"questionCount"`
	TimeLimit     uint                      `json:"timeLimit"`

	Participant       TestParticipantDTO `json:"participant"`
	Points            float64            `json:"points"`
	BonusPoints       float64            `json:"bonusPoints"`
	BonusPointsReason string             `json:"bonusPointsReason"`
}

func (m TestInstanceDTO) From(d *models.TestInstance, isTutor bool) TestInstanceDTO {
	dto := TestInstanceDTO{
		ID:            d.ID,
		State:         d.State,
		StartedAt:     d.StartedAt,
		EndedAt:       d.EndedAt,
		EndsAt:        d.EndsAt,
		QuestionCount: uint(len(d.Questions)),
		TimeLimit:     d.CourseItem.TestDetail.TimeLimit,
		Participant:   TestParticipantDTO{}.From(d.Participant),
	}

	if isTutor || d.State == enums.TestInstanceStateActive {
		dto.Questions = make([]TestInstanceQuestionDTO, len(d.Questions))
		for q_i, q := range d.Questions {
			dto.Questions[q_i] = TestInstanceQuestionDTO{}.From(&q, isTutor)
		}
	}

	if isTutor {
		dto.Points = d.Result.Points
		dto.BonusPoints = d.BonusPoints
		dto.BonusPointsReason = d.BonusPointsReason
	}

	return dto
}
