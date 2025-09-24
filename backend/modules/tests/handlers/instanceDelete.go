package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestInstanceDeleteResponse struct {
	Success bool `json:"success"`
}

// @Summary Deletes test instance and its result
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param instanceId path int true "ID of the corresponding test instance"
// @Success 200 {object} TestInstanceDeleteResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/{courseItemId}/instance/{instanceId} [delete]
func TestInstanceDelete(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
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
				Message: "Not enough permission for this item",
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

	transaction := initializers.DB.Begin()

	testRepo := repositories.NewTestRepository()
	modifier := func(db *gorm.DB) *gorm.DB {
		return db.Preload("Result")
	}
	testInstance, err := testRepo.GetTestInstanceByID(transaction, params.InstanceID, userData.ID, &modifier, false, false, &courseItem.ID, nil)
	if err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.
		Delete(&testInstance.Result).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Failed to delete instance result",
		}
	}

	if err := transaction.
		Delete(&testInstance).Error; err != nil {
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

	c.JSON(200, TestInstanceDeleteResponse{
		Success: true,
	})

	return nil
}
