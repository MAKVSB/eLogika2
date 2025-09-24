package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/dtos"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TestInstanceGetTelemetryResponse struct {
	Items []dtos.TestInstanceEventDTO `json:"items"`
}

// @Summary Gets events associated with test instance
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param instanceId path int true "ID of the corresponding test instance"
// @Success 200 {object} TestInstanceGetTelemetryResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/{courseItemId}/instance/{instanceId}/telemetry [GET]
func TestInstanceGetTelemetry(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			InstanceID   uint `uri:"instanceId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}
	var courseItem models.CourseItem
	// Check if tutor/garant can view/modify courseItem
	if userRole == enums.CourseUserRoleAdmin {
	} else if userRole == enums.CourseUserRoleGarant {
		var test models.TestInstance
		if err := initializers.DB.
			InnerJoins("Test").
			Find(&test, params.InstanceID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to commit changes",
			}
		}

		if err := initializers.DB.
			Preload("TestDetail").
			Where("managed_by = ?", enums.CourseUserRoleGarant).
			Find(&courseItem, test.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else if userRole == enums.CourseUserRoleTutor {
		var test models.TestInstance
		if err := initializers.DB.
			InnerJoins("Test").
			Find(&test, params.InstanceID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to commit changes",
			}
		}

		if err := initializers.DB.
			Preload("TestDetail").
			Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userData.ID).
			Find(&courseItem, test.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	var testInstanceEvents []models.TestInstanceEvent

	if err := initializers.DB.
		Preload("User").
		Where("test_instance_id = ?", params.InstanceID).
		Find(&testInstanceEvents).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.TestInstanceEventDTO, len(testInstanceEvents))
	for i, e := range testInstanceEvents {
		dtoList[i] = dtos.TestInstanceEventDTO{}.From(&e)
	}

	c.JSON(200, TestInstanceGetTelemetryResponse{
		Items: dtoList,
	})

	return nil
}
