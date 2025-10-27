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
	PointsFinal       bool               `json:"pointsFinal"`
	PointsMin         float64            `json:"pointsMin"`
	PointsMax         float64            `json:"pointsMax"`
	BonusPoints       float64            `json:"bonusPoints"`
	BonusPointsReason string             `json:"bonusPointsReason"`

	ShowContent     bool `json:"showContent"`
	ShowCorrectness bool `json:"showCorrectness"`

	RecognizerFiles []FileDTO
}

func (m TestInstanceDTO) From(
	d *models.TestInstance,
	isTutor bool,
) TestInstanceDTO {
	showResults := true
	showCorrectness := d.CourseItem.TestDetail.ShowCorrectness
	showTestContent := false
	showLayout := d.CourseItem.TestDetail.IsPaper

	if !d.CourseItem.TestDetail.ShowResults {
		if d.Term.ActiveTo.After(time.Now()) {
			showResults = false
		}
	}

	if d.State == enums.TestInstanceStateFinished {
		showTestContent = d.CourseItem.TestDetail.ShowTest
	}

	if d.State == enums.TestInstanceStateActive {
		showTestContent = true
		showCorrectness = false
		showResults = false
		showLayout = true
	}

	if isTutor {
		showTestContent = true
		showCorrectness = true
		showResults = true
		showLayout = true
	}

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

	if showLayout || showTestContent {
		dto.Questions = make([]TestInstanceQuestionDTO, len(d.Questions))
		for q_i, q := range d.Questions {
			dto.Questions[q_i] = TestInstanceQuestionDTO{}.From(q, isTutor, showTestContent, showCorrectness)
		}
	}

	if showResults {
		dto.Points = d.Result.Points
		dto.PointsMax = float64(d.CourseItem.PointsMax)
		dto.PointsMin = float64(d.CourseItem.PointsMin)
		dto.PointsFinal = d.Result.Final
		dto.BonusPoints = d.BonusPoints
		dto.BonusPointsReason = d.BonusPointsReason
		dto.RecognizerFiles = make([]FileDTO, len(d.RecognizerFiles))

		for i_i, image := range d.RecognizerFiles {
			dto.RecognizerFiles[i_i] = FileDTO{}.From(image.File)
		}
	}

	dto.ShowContent = showTestContent
	dto.ShowCorrectness = showCorrectness

	return dto
}
