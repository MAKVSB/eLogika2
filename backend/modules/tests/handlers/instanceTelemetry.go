package handlers

import (
	"encoding/json"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TestInstanceEventData struct {
	OccuredAt time.Time                       `json:"occuredAt"`
	EventType enums.TestInstanceEventTypeEnum `json:"eventType"`
	EventData json.RawMessage                 `json:"eventData"`
	PageID    string                          `json:"pageId"`
}

type TestInstanceTelemetryRequest struct {
	Events []models.TestInstanceEvent `json:"events"`
}

type TestInstanceTelemetryResponse struct {
}

// @Summary Saves telemetry for online test
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TestInstanceTelemetryRequest true "Ability to filter results"
// @Success 200 {object} TestInstanceTelemetryResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/tests/{instanceId}/telemetry [put]
func TestInstanceTelemetry(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			InstanceID uint `uri:"instanceId" binding:"required"`
		},
		TestInstanceTelemetryRequest,
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

	var testInstance models.TestInstance
	if err := initializers.DB.
		Select("id").
		Where("participant_id = ?", userData.ID).
		Find(&testInstance, params.InstanceID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    404,
			Message: "Failed to find instance",
		}
	}

	events := make([]models.TestInstanceEvent, len(reqData.Events))

	for i, event := range reqData.Events {
		events[i] = models.TestInstanceEvent{
			TestInstanceID: params.InstanceID,
			UserID:         userData.ID,
			OccuredAt:      event.OccuredAt,
			ReceivedAt:     time.Now(),
			EventSource:    enums.TestInstanceEventSourceClient,
			EventType:      event.EventType,
			EventData:      event.EventData,
			PageID:         event.PageID,
		}
	}

	if err := initializers.DB.Save(&events).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to save test events",
			Details: err.Error(),
		}
	}

	c.JSON(200, TestInstanceTelemetryResponse{})

	return nil
}
