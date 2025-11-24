package handlers

import (
	"encoding/json"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TestInstanceStartResponse struct {
	InstanceID uint `json:"instanceId"`
}

// @Summary Starts test instance for user
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TestGeneratorRequest true "Ability to filter results"
// @Success 200 {object} TestGeneratorResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/instance/start [post]
func TestInstanceStart(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			InstanceID uint `uri:"instanceId" binding:"required"`
		},
		any,
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

	timeFreeze := time.Now()

	transaction := initializers.DB.Begin()

	testRepo := repositories.NewTestRepository()
	testInstance, err := testRepo.GetTestInstanceByID(transaction, params.InstanceID, userData.ID, nil, true, true, nil, &userData.ID)
	if err != nil {
		transaction.Rollback()
		return err
	}

	if testInstance.CourseItem.TestDetail.IPRanges != "" {
		if !utils.IsIPAllowed(testInstance.CourseItem.TestDetail.IPRanges, c.ClientIP()) {
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

	if testInstance.Form != enums.TestInstanceFormOnline {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "Instance is not set as online",
		}
	}

	if testInstance.State != enums.TestInstanceStateReady {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "Instance not in prepared state",
		}
	}

	testInstance.State = enums.TestInstanceStateActive
	testInstance.StartedAt = timeFreeze

	endsAt1 := timeFreeze.Add(time.Minute * time.Duration(testInstance.CourseItem.TestDetail.TimeLimit))
	endsAt2 := testInstance.Term.ActiveTo

	// TODO CONSULT Jak určovat čas konce ? takto ?
	if endsAt1.Before(endsAt2) {
		testInstance.EndsAt = endsAt1
	} else {
		testInstance.EndsAt = endsAt2
	}

	var runningUserInstances []*models.TestInstance
	if err := transaction.
		Where("state = ?", enums.TestInstanceStateActive).
		Where("participant_id = ?", userData.ID).
		Find(&runningUserInstances).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to check for running instances",
			Details: err.Error(),
		}
	}

	// Check that no test is actively running
	// TODO maybe lock table ?
	if len(runningUserInstances) != 0 {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Another test instance is already running",
		}
	}

	if err := transaction.Save(&testInstance).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to create test instance",
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

	c.JSON(200, TestInstanceStartResponse{
		InstanceID: testInstance.ID,
	})

	return nil
}
