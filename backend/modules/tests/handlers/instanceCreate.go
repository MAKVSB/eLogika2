package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/helpers"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TestInstanceCreateRequest struct {
	Form enums.TestInstanceFormEnum `json:"form" binding:"required"`
}

type TestInstanceCreateResponse struct {
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
// @Router /api/v2/courses/{courseId}/tests/{courseItemId}/instances/{testId}/create [post]
func CreateInstance(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			TestID       uint `uri:"testId" binding:"required"`
		},
		TestInstanceCreateRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}
	// Check if tutor/garant can view/modify courseItem
	if userRole == enums.CourseUserRoleAdmin {
	} else if userRole == enums.CourseUserRoleGarant {
		var courseItem models.CourseItem
		if err := initializers.DB.
			Where("managed_by = ?", enums.CourseUserRoleGarant).
			Find(&courseItem, params.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else if userRole == enums.CourseUserRoleTutor {
		var courseItem models.CourseItem
		if err := initializers.DB.
			Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userData.ID).
			Find(&courseItem, params.CourseItemID).Error; err != nil {
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

	var test *models.Test
	if err := initializers.DB.
		Where("course_id = ?", params.CourseID).
		Where("course_item_id = ?", params.CourseItemID).
		Preload("Questions").
		Preload("Questions.Answers").
		First(&test, params.TestID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch term",
			Details: err.Error(),
		}
	}

	testInstance, err := helpers.CreateInstance(
		initializers.DB,
		test,
		userData.ID,
		test.TermID,
		params.CourseItemID,
		reqData.Form,
	)
	if err != nil {
		return err
	}

	c.JSON(200, TestInstanceCreateResponse{
		InstanceID: testInstance.ID,
	})

	return nil
}
