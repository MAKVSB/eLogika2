package handlers

import (
	"encoding/json"
	"time"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TestInstanceGetResponse struct {
	InstanceData dtos.TestInstanceDTO `json:"instanceData"`
}

// @Summary Gets test instance for user
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param instanceId path int true "ID of the corresponding test instance"
// @Success 200 {object} TestInstanceGetResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/tests/{instanceId} [GET]
func TestInstanceGet(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
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

	if userRole != enums.CourseUserRoleStudent {
		return &common.ErrorResponse{
			Code:    403,
			Message: "User is not a student",
		}
	}

	// TODO validate from here

	testRepo := repositories.NewTestRepository()
	testInstance, err := testRepo.GetTestInstanceByID(initializers.DB, params.InstanceID, userData.ID, nil, true, true, nil, &userData.ID)
	if err != nil {
		return err
	}

	if testInstance.State == enums.TestInstanceStateActive && testInstance.CourseItem.TestDetail.IPRanges != "" {
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

	c.JSON(200, TestInstanceGetResponse{
		InstanceData: dtos.TestInstanceDTO{}.From(
			testInstance,
			false,
		),
	})

	return nil
}
