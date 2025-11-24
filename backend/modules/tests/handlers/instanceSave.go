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
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestInstanceSaveResponse struct {
}

// @Summary Starts test instance for user
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body helpers.TestInstanceQuestion true "Ability to filter results"
// @Success 200 {object} TestInstanceSaveResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/tests/{instanceId}/save [put]
func TestInstanceSave(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			InstanceID uint `uri:"instanceId" binding:"required"`
		},
		helpers.TestInstanceQuestion,
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

	var testInstanceQuestion models.TestInstanceQuestion
	if err := transaction.
		InnerJoins("TestInstance", initializers.DB.Where("TestInstance.participant_id = ? AND TestInstance.id = ? AND TestInstance.state = ?", userData.ID, params.InstanceID, enums.TestInstanceStateActive)).
		Joins("TestQuestion", initializers.DB.Unscoped()).
		Joins("TestQuestion.Question", initializers.DB.Unscoped()).
		Preload("Answers", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("Answers.TestQuestionAnswer", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		First(&testInstanceQuestion, reqData.QuestionID).Error; err != nil {
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
		First(&courseItem, testInstanceQuestion.TestInstance.CourseItemID).Error; err != nil {
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

	if testInstanceQuestion.TestInstance.IsExpired(time.Now()) {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Time expired",
		}
	}

	events := make([]*models.TestInstanceEvent, 0)

	switch testInstanceQuestion.TestQuestion.Question.QuestionFormat {
	case enums.QuestionFormatOpen:
		err = helpers.UpdateOpenQuestion(&testInstanceQuestion, reqData, transaction, userData.ID, false, &events)
		if err != nil {
			transaction.Rollback()
			return err
		}
	case enums.QuestionFormatTest:
		err = helpers.UpdateTestQuestion(&testInstanceQuestion, reqData, transaction, userData.ID, &events)
		if err != nil {
			transaction.Rollback()
			return err
		}
	default:
		panic(fmt.Sprintf("unexpected enums.QuestionFormatEnum: %#v", testInstanceQuestion.TestQuestion.Question.QuestionFormat))
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

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	c.JSON(200, TestInstanceSaveResponse{})

	return nil
}
