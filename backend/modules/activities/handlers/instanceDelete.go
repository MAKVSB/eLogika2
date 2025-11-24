package handlers

import (
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
)

type ActivityInstanceDeleteResponse struct {
	Success bool `json:"success"`
}

// @Summary Deletes test instance and its result
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param instanceId path int true "ID of the corresponding test instance"
// @Success 200 {object} TestGeneratorResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/activities/instance/{instanceId} [delete]
func ActivityInstanceDelete(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID   uint `uri:"courseId" binding:"required"`
			InstanceID uint `uri:"instanceId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	var activityInstance *models.ActivityInstance
	if err := initializers.DB.
		InnerJoins("Term").
		InnerJoins("CourseItem").
		InnerJoins("CourseItem.ActivityDetail").
		InnerJoins("Result").
		First(&activityInstance, params.InstanceID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    404,
			Message: "Activity instance does not exists",
			Details: err.Error(),
		}
	}

	// Check if tutor/garant can view/modify courseItem
	courseItemService := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, activityInstance.CourseItemID, userData.ID, userRole, nil, false, nil)
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

	if err := transaction.
		Delete(&activityInstance.Result).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Failed to delete instance result",
		}
	}

	if err := transaction.
		Delete(&activityInstance).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Failed to delete instance",
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

	c.JSON(200, ActivityInstanceDeleteResponse{
		Success: true,
	})

	return nil
}
