package helpers

import (
	"encoding/json"
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

type TestInstanceAnswer struct {
	AnswerID uint `json:"id" binding:"required"`
	Selected bool `json:"selected" binding:"required"`
}

type TestInstanceQuestion struct {
	QuestionID           uint                  `json:"id" binding:"required"`
	TextAnswer           *models.TipTapContent `json:"textAnswer"`
	TextAnswerPercentage *float64              `json:"textAnswerPercentage"` // Only for teacher endpoint
	TextAnswerReviewed   *bool                 `json:"textAnswerReviewed"`   // Only for teacher endpoint
	Answers              []TestInstanceAnswer  `json:"answers"`
}

func UpdateOpenQuestion(ti_q *models.TestInstanceQuestion, rd_q *TestInstanceQuestion, transaction *gorm.DB, userId uint, isTutor bool, events *[]*models.TestInstanceEvent) *common.ErrorResponse {
	if rd_q.TextAnswer != nil {
		if ti_q.TextAnswer == nil {
			ti_q.TextAnswer = rd_q.TextAnswer
		} else {
			reqAnswer, err := rd_q.TextAnswer.Hash()
			if err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to calculate hash",
					Details: err.Error(),
				}
			}

			tiAnswer, err := ti_q.TextAnswer.Hash()
			if err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to calculate hash",
					Details: err.Error(),
				}
			}

			if tiAnswer != reqAnswer {
				ti_q.TextAnswer = rd_q.TextAnswer
			}
		}

		if rd_q.TextAnswerReviewed != nil {
			ti_q.TextAnswerReviewedByID = &userId
			ti_q.TextAnswerPercentage = *rd_q.TextAnswerPercentage
		}

		// TODO check why is it not saving objects inside objects
		if err := transaction.Save(&ti_q).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to save text answer for question",
				Details: err.Error(),
			}
		}

		eventData, _ := json.Marshal(map[string]interface{}{
			"QuestionOrder": ti_q.TestQuestion.Order,
			"AnswerData":    rd_q.TextAnswer,
		})

		*events = append(*events, &models.TestInstanceEvent{
			TestInstanceID: ti_q.TestInstanceID,
			UserID:         userId,
			OccuredAt:      time.Now(),
			EventSource:    enums.TestInstanceEventSourceServer,
			EventType:      enums.TestInstanceEventTypeQuestionUpdate,
			EventData:      eventData,
			PageID:         "",
		})
	}
	return nil
}

func UpdateTestQuestion(ti_q *models.TestInstanceQuestion, rd_q *TestInstanceQuestion, transaction *gorm.DB, userId uint, events *[]*models.TestInstanceEvent) *common.ErrorResponse {
	for _, ra := range rd_q.Answers {
		ta := FindAnswer(ti_q, ra.AnswerID)
		if ta == nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to find answer",
			}
		}

		if ta.Selected != ra.Selected {
			ta.Selected = ra.Selected

			// TODO fix why global save is not saving
			if err := transaction.Save(&ta).Error; err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to save answers for question",
					Details: err.Error(),
				}
			}

			eventData, _ := json.Marshal(map[string]interface{}{
				"QuestionOrder": ti_q.TestQuestion.Order,
				"AnswerOrder":   ta.TestQuestionAnswer.Order,
				"AnswerData":    ra.Selected,
			})

			*events = append(*events, &models.TestInstanceEvent{
				TestInstanceID: ti_q.TestInstanceID,
				UserID:         userId,
				OccuredAt:      time.Now(),
				EventSource:    enums.TestInstanceEventSourceServer,
				EventType:      enums.TestInstanceEventTypeQuestionUpdate,
				EventData:      eventData,
				PageID:         "",
			})
		}
	}
	return nil
}

func FindAnswer(ti_q *models.TestInstanceQuestion, a_id uint) *models.TestInstanceQuestionAnswer {
	for _, ti_a := range ti_q.Answers {
		if ti_a.ID == a_id {
			return ti_a
		}
	}

	return nil
}

func FindQuestion(ti *models.TestInstance, q_id uint) *models.TestInstanceQuestion {
	for _, ti_a := range ti.Questions {
		if ti_a.ID == q_id {
			return ti_a
		}
	}

	return nil
}
