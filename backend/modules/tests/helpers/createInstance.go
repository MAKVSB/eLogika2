package helpers

import (
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"gorm.io/gorm"
)

func CreateInstance(
	transaction *gorm.DB,
	generatedTest *models.Test,
	participantId uint,
	termId uint,
	courseItemId uint,
	form enums.TestInstanceFormEnum,
) (*models.TestInstance, *common.ErrorResponse) {
	testInstance := &models.TestInstance{
		State:         enums.TestInstanceStateReady,
		Form:          form,
		ParticipantID: participantId,
		TestID:        generatedTest.ID,
		TermID:        termId,
		CourseItemID:  courseItemId,
		Questions:     make([]models.TestInstanceQuestion, 0),
	}

	for _, generatedTestQuestion := range generatedTest.Questions {
		testInstanceQuestion := models.TestInstanceQuestion{
			TestQuestionID: generatedTestQuestion.ID,
			Answers:        make([]models.TestInstanceQuestionAnswer, 0),
		}

		for _, generatedTestQuestionAnswer := range generatedTestQuestion.Answers {
			testInstanceQuestionAnswer := models.TestInstanceQuestionAnswer{
				TestQuestionAnswerID: generatedTestQuestionAnswer.ID,
				Selected:             false,
			}

			testInstanceQuestion.Answers = append(testInstanceQuestion.Answers, testInstanceQuestionAnswer)
		}

		testInstance.Questions = append(testInstance.Questions, testInstanceQuestion)
	}

	// Add result object
	testInstance.Result = &models.CourseItemResult{
		Version:      0,
		CourseItemID: testInstance.CourseItemID,
		TermID:       testInstance.TermID,
		StudentID:    testInstance.ParticipantID,
	}

	if err := transaction.Save(&testInstance).Error; err != nil {
		return nil, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to create test instance",
			Details: err.Error(),
		}
	}

	return testInstance, nil
}
