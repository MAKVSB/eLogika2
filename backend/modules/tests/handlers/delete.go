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

type TestDeleteResponse struct {
	Success bool `json:"success"`
}

// @Summary Deletes test and all its instances
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
// @Router /api/v2/courses/{courseId}/tests/{courseItemId}/instances/{testId} [delete]
func TestDelete(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			TestID       uint `uri:"testId" binding:"required"`
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
		var test models.Test
		if err := initializers.DB.
			Find(&test, params.TestID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to load test",
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
		var test models.Test
		if err := initializers.DB.
			Find(&test, params.TestID).Error; err != nil {
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
		return db.Preload("Instances").Preload("Instances.Result")
	}
	test, err := testRepo.GetTestByID(transaction, params.CourseID, params.TestID, userData.ID, &modifier, false)
	if err != nil {
		transaction.Rollback()
		return err
	}

	for _, v := range test.Instances {
		if err := transaction.
			Delete(&v.Result).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    403,
				Message: "Failed to delete instance result",
			}
		}

		if err := transaction.
			Delete(&v).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    403,
				Message: "Failed to delete instance",
			}
		}
	}

	if err := transaction.
		Delete(&test).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Failed to delete test",
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

	c.JSON(200, TestDeleteResponse{
		Success: true,
	})

	return nil
}
