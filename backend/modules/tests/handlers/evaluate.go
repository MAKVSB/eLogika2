package handlers

import (
	"fmt"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestEvaluationRequest struct {
	CourseItemID uint  `json:"courseItemId" binding:"required"`
	TermID       *uint `json:"termId"`
	TestID       *uint `json:"testId"`
	InstanceID   *uint `json:"instanceId"`
}

type TestEvaluationResponse struct {
	Success bool `json:"success"`
}

// @Summary (Re)evaluates test
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param instanceId path int true "ID of the corresponding test instance"
// @Param body body TestEvaluationRequest true "Ability to filter results"
// @Success 200 {object} EvaluateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/evaluate [POST]

func TestEvaluate(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		TestEvaluationRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	// Check if tutor/garant can view/modify courseItem
	courseItemService := services_course_item.CourseItemService{}
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, reqData.CourseItemID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	if reqData.InstanceID != nil {
		err = EvaluateTestInstance(transaction, *reqData.InstanceID, &userData, true)
		if err != nil {
			return err
		}

		var instanceData models.TestInstance
		if err := transaction.Find(&instanceData, reqData.InstanceID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to load instance data",
				Details: err.Error(),
			}
		}

		rootCoureItem := courseItem.ID
		if courseItem.ParentID != nil {
			rootCoureItem = *courseItem.ParentID
		}

		services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		err = services_course_item.UpdateSelectedResults(transaction, params.CourseID, rootCoureItem, instanceData.ParticipantID)
		if err != nil {
			transaction.Rollback()
			return err
		}

	} else {
		var instanceIDs []uint
		query := transaction.
			Model(models.TestInstance{}).
			Where("course_item_id = ?", courseItem.ID)

		if reqData.TermID != nil {
			query = query.Where("term_id = ?", reqData.TermID)
		}

		if reqData.TestID != nil {
			query = query.Where("test_id = ?", reqData.TestID)
		}

		if err := query.Pluck("id", &instanceIDs).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch test",
				Details: err.Error(),
			}
		}

		for _, instanceID := range instanceIDs {

			var instanceData models.TestInstance
			if err := transaction.Find(&instanceData, instanceID).Error; err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to load instance data",
					Details: err.Error(),
				}
			}

			err = EvaluateTestInstance(transaction, instanceID, &userData, true)
			if err != nil {
				return err
			}

			rootCoureItem := courseItem.ID
			if courseItem.ParentID != nil {
				rootCoureItem = *courseItem.ParentID
			}

			services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
			err = services_course_item.UpdateSelectedResults(transaction, params.CourseID, rootCoureItem, instanceData.ParticipantID)
			if err != nil {
				transaction.Rollback()
				return err
			}
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
			Details: err.Error(),
		}
	}

	c.JSON(200, TestEvaluationResponse{
		Success: true,
	})

	return nil
}

type EvaluationBlock struct {
	Weight         uint
	QuestionFormat enums.QuestionFormatEnum

	// Question format = ABCD questions
	TotalAnswers          float64
	CorrectlyAnswered     float64
	IncorrectlyAnswered   float64
	WrongAnswerPercentage uint
	AllowEmptyAnswers     bool

	// Question format = Open questions
	TotalQuestions          float64
	TextAnswerPercentageSum float64
	TextAnswerReviewed      bool
}

func EvaluateTestInstance(dbRef *gorm.DB, instanceID uint, userData *authdtos.LoggedUserDTO, silentSkip bool) *common.ErrorResponse {
	var testInstance models.TestInstance
	if err := dbRef.
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("TestQuestion").
				Joins("TestQuestion.Question")
		}).
		Preload("Questions.Answers", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("TestQuestionAnswer").
				Joins("TestQuestionAnswer.Answer")
		}).
		Joins("Test").
		Joins("CourseItem").
		Joins("CourseItem.Parent").
		Joins("Result").
		Find(&testInstance, instanceID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to load instance data",
			Details: err.Error(),
		}
	}
	if testInstance.State != enums.TestInstanceStateFinished && testInstance.State != enums.TestInstanceStateExpired {
		if silentSkip {
			return nil
		} else {
			return &common.ErrorResponse{
				Code:    401,
				Message: "Cannot evaluate points for active test",
			}
		}
	}

	// remap back into blocks
	evaluationMap := make(map[uint]*EvaluationBlock)

	for _, block := range testInstance.Test.Blocks {
		evaluationMap[block.ID] = &EvaluationBlock{
			Weight:                block.Weight,
			WrongAnswerPercentage: block.WrongAnswerPercentage,
			AllowEmptyAnswers:     block.AllowEmptyAnswers,
		}
	}

	for _, q := range testInstance.Questions {
		block := evaluationMap[q.TestQuestion.BlockID]

		switch q.TestQuestion.Question.QuestionFormat {
		case enums.QuestionFormatOpen:
			block.TotalQuestions++
			block.TextAnswerPercentageSum += float64(q.TextAnswerPercentage)
			block.TextAnswerReviewed = block.TextAnswerReviewed && q.TextAnswerReviewedByID != nil
			block.QuestionFormat = q.TestQuestion.Question.QuestionFormat
		case enums.QuestionFormatTest:
			total, correct, incorrect := QuestionAnswersScore(q.Answers, block.AllowEmptyAnswers)
			block.TotalAnswers += total
			block.CorrectlyAnswered += correct
			block.IncorrectlyAnswered += incorrect
			block.QuestionFormat = q.TestQuestion.Question.QuestionFormat
		default:
			panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", q.TestQuestion.Question.QuestionFormat))
		}
	}

	// calculate points
	points := float64(0)
	final := true

	for _, block := range evaluationMap {
		blockWeight := utils.ToPercentage(block.Weight)
		testPointsMax := float64(testInstance.CourseItem.PointsMax)

		switch block.QuestionFormat {
		case enums.QuestionFormatOpen:
			if !block.TextAnswerReviewed {
				final = false
			}
			textAnswerPercentage := block.TextAnswerPercentageSum / 100 * block.TotalQuestions

			points += textAnswerPercentage * (testPointsMax * blockWeight)
		case enums.QuestionFormatTest:
			wrongAnswerPercentage := utils.ToPercentage(block.WrongAnswerPercentage)

			ratio := (block.CorrectlyAnswered - wrongAnswerPercentage*block.IncorrectlyAnswered) / block.TotalAnswers
			points += ratio * (testPointsMax * blockWeight)
		default:
			utils.DebugPrintJSON(evaluationMap)
			panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", block.QuestionFormat))
		}
	}

	points += testInstance.BonusPoints

	testInstance.Result.Version = testInstance.Result.Version + 1
	testInstance.Result.Points = utils.RoundToEven(points, 2)
	testInstance.Result.Final = final
	if userData != nil {
		testInstance.Result.UpdatedByID = &userData.ID
	} else {
		testInstance.Result.UpdatedByID = nil
	}
	if err := dbRef.Save(&testInstance.Result).Error; err != nil {
		return &common.ErrorResponse{
			Message: "Failed to save result",
			Details: err.Error(),
		}
	}

	return nil
}

func QuestionAnswersScore(answers []*models.TestInstanceQuestionAnswer, allowEmptyAnswers bool) (float64, float64, float64) {
	var numAnswers = float64(0)
	var numChecked = float64(0)
	var numCorrect = float64(0)
	var numIncorrect = float64(0)

	// Assings answer statistics for question
	for _, qa := range answers {
		numAnswers++
		if qa.Selected {
			numChecked++
		}

		if qa.TestQuestionAnswer.Answer.Correct == qa.Selected {
			numCorrect++
		} else {
			numIncorrect++
		}
	}

	if allowEmptyAnswers && (numChecked == 0 || numChecked == numAnswers) {
		return numAnswers, 0, 0
	}

	return numAnswers, numCorrect, numIncorrect
}
