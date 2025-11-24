package handlers

import (
	"encoding/json"
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/dtos"
	"elogika.vsb.cz/backend/modules/tests/helpers"
	"elogika.vsb.cz/backend/repositories"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TestInstanceTutorSaveRequest struct {
	Questions         []helpers.TestInstanceQuestion `json:"questions" binding:"required"`
	BonusPoints       float64                        `json:"bonusPoints"`
	BonusPointsReason string                         `json:"bonusPointsReason"`
}

type TestInstanceTutorSaveResponse struct {
	InstanceData dtos.TestInstanceDTO `json:"instanceData"`
}

// @Summary Saves test instance for user by tutor (advanced options)
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TestInstanceTutorSaveRequest true "Ability to filter results"
// @Success 200 {object} TestInstanceTutorSaveResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/{courseItemId}/instance/{instanceId}/tutorsave [put]
func TestInstanceTutorSave(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			InstanceID   uint `uri:"instanceId" binding:"required"`
		},
		TestInstanceTutorSaveRequest,
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
	courseItemService := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}
	if !courseItem.Editable {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	transaction := initializers.DB.Begin()

	testRepo := repositories.NewTestRepository()
	testInstance, err := testRepo.GetTestInstanceByID(transaction, params.InstanceID, userData.ID, nil, true, false, &params.CourseItemID, nil)
	if err != nil {
		transaction.Rollback()
		return err
	}

	events := make([]*models.TestInstanceEvent, 0)

	if testInstance.BonusPoints != reqData.BonusPoints || testInstance.BonusPointsReason != reqData.BonusPointsReason {

		data := map[string]interface{}{
			"points": reqData.BonusPoints,
			"reason": reqData.BonusPointsReason,
		}
		rawData, _ := json.Marshal(data)

		events = append(events, &models.TestInstanceEvent{
			TestInstanceID: testInstance.ID,
			UserID:         userData.ID,
			OccuredAt:      time.Now(),
			ReceivedAt:     time.Now(),
			EventSource:    enums.TestInstanceEventSourceServer,
			EventType:      enums.TestInstanceEventTypeBonusPointsModified,
			EventData:      rawData,
		})
	}
	testInstance.BonusPoints = reqData.BonusPoints
	testInstance.BonusPointsReason = reqData.BonusPointsReason
	switch testInstance.State {
	case enums.TestInstanceStateReady:
		testInstance.StartedAt = time.Now()
		testInstance.EndsAt = time.Now()
		testInstance.EndedAt = time.Now()
		testInstance.State = enums.TestInstanceStateFinished
	case enums.TestInstanceStateActive:
		testInstance.EndedAt = time.Now()
		testInstance.State = enums.TestInstanceStateFinished
	}

	for _, rd_q := range reqData.Questions {
		ti_q := helpers.FindQuestion(testInstance, rd_q.QuestionID)
		if ti_q == nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to find question",
			}
		}

		switch ti_q.TestQuestion.Question.QuestionFormat {
		case enums.QuestionFormatOpen:
			err = helpers.UpdateOpenQuestion(ti_q, &rd_q, transaction, userData.ID, true, &events)
			if err != nil {
				transaction.Rollback()
				return err
			}
			if rd_q.TextAnswerReviewed != nil {
				ti_q.TextAnswerReviewedByID = &userData.ID
				ti_q.TextAnswerPercentage = *rd_q.TextAnswerPercentage
			}
		case enums.QuestionFormatTest:
			err = helpers.UpdateTestQuestion(ti_q, &rd_q, transaction, userData.ID, &events)
			if err != nil {
				transaction.Rollback()
				return err
			}
		default:
			return &common.ErrorResponse{
				Code:    500,
				Message: "Invalid question format",
				Details: ti_q.TestQuestion.Question.QuestionFormat,
			}
		}
	}
	if len(events) != 0 {
		if err := transaction.Save(&events).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to save events for question",
				Details: err.Error(),
			}
		}
	}

	if err := transaction.Save(&testInstance).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to save answers for question",
			Details: err.Error(),
		}
	}

	err = EvaluateTestInstance(transaction, params.InstanceID, &userData, false)
	if err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to evaluate test",
			Details: err,
		}
	}

	rootCoureItem := courseItem.ID
	if courseItem.ParentID != nil {
		rootCoureItem = *courseItem.ParentID
	}

	services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
	err = services_course_item.UpdateSelectedResults(transaction, params.CourseID, rootCoureItem, testInstance.ParticipantID)
	if err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
			Details: err.Error(),
		}
	}

	testInstance, err = testRepo.GetTestInstanceByID(initializers.DB, params.InstanceID, userData.ID, nil, true, true, &params.CourseItemID, nil)
	if err != nil {
		return err
	}

	c.JSON(200, TestInstanceTutorSaveResponse{
		InstanceData: dtos.TestInstanceDTO{}.From(
			testInstance,
			true,
		),
	})

	return nil
}
