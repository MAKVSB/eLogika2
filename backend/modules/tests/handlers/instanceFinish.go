package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/helpers"
	"elogika.vsb.cz/backend/repositories"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestInstanceFinishRequest struct {
	Questions []helpers.TestInstanceQuestion `json:"questions" binding:"required"`
}

type TestInstanceFinishResponse struct {
}

// @Summary Starts test instance for user
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TestInstanceFinishRequest true "Ability to filter results"
// @Success 200 {object} TestInstanceFinishResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/tests/{instanceId}/finish [put]
func TestInstanceFinish(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			InstanceID uint `uri:"instanceId" binding:"required"`
		},
		TestInstanceFinishRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	if userRole != enums.CourseUserRoleStudent {
		return &common.ErrorResponse{
			Code:    403,
			Message: "User is not a student",
		}
	}

	transaction := initializers.DB.Begin()
	var testInstance models.TestInstance
	if err := transaction.
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.
				Unscoped().
				Joins("TestQuestion").
				Joins("TestQuestion.Question")
		}).
		Preload("Questions.Answers", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped().InnerJoins("TestQuestionAnswer", initializers.DB.Unscoped())
		}).
		Where("participant_id = ?", userData.ID).
		Where("state = ?", enums.TestInstanceStateActive).
		First(&testInstance, params.InstanceID).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch test instance",
			Details: err.Error(),
		}
	}

	var courseItem *models.CourseItem
	if err := transaction.
		InnerJoins("TestDetail").
		First(&courseItem, testInstance.CourseItemID).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch test instance",
			Details: err.Error(),
		}
	}

	if courseItem.TestDetail.IPRanges != "" {
		if !utils.IsIPAllowed(courseItem.TestDetail.IPRanges, c.ClientIP()) {
			if err := initializers.DB.Create(&models.TestInstanceEvent{
				TestInstanceID: params.InstanceID,
				UserID:         userData.ID,
				OccuredAt:      time.Time{},
				ReceivedAt:     time.Time{},
				EventSource:    enums.TestInstanceEventSourceServer,
				EventType:      enums.TestInstanceEventTypeQuestionInvalidIP,
				EventData:      json.RawMessage(c.ClientIP()),
			}).Error; err != nil {
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to insert log report",
				}
			}
			return &common.ErrorResponse{
				Code:    403,
				Message: "IP is not in range of allowed IPs",
			}
		}
	}

	if testInstance.IsExpired(time.Now()) {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Time expired",
		}
	}

	testInstance.State = enums.TestInstanceStateFinished
	testInstance.EndedAt = time.Now()
	events := make([]*models.TestInstanceEvent, 0)
	for _, rd_q := range reqData.Questions {
		ti_q := helpers.FindQuestion(&testInstance, rd_q.QuestionID)
		if ti_q == nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to find question",
			}
		}

		switch ti_q.TestQuestion.Question.QuestionFormat {
		case enums.QuestionFormatOpen:
			err = helpers.UpdateOpenQuestion(ti_q, &rd_q, transaction, userData.ID, false, &events)
			if err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to find question",
				}
			}
		case enums.QuestionFormatTest:
			err = helpers.UpdateTestQuestion(ti_q, &rd_q, transaction, userData.ID, &events)
			if err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to find question",
				}
			}
		default:
			panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", ti_q.TestQuestion.Question.QuestionFormat))
		}
	}

	if len(events) != 0 {
		if err := transaction.Save(&events).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to save events",
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

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	err = EvaluateTestInstance(initializers.DB, params.InstanceID, nil, false)
	if err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to evaluate test",
			Details: "At least the answers have been correctly saved",
		}
	}

	rootCoureItem := courseItem.ID
	if courseItem.ParentID != nil {
		rootCoureItem = *courseItem.ParentID
	}

	services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
	err = services_course_item.UpdateSelectedResults(initializers.DB, courseItem.CourseID, rootCoureItem, testInstance.ParticipantID)
	if err != nil {
		return err
	}

	c.JSON(200, TestInstanceFinishResponse{})

	return nil
}
