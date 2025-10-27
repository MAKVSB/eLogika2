package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Newly updated course item
type CourseItemDeleteResponse struct {
	Success bool `json:"success"`
}

// @Summary Deletes course item
// @Tags CourseItems
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body CourseItemUpdateRequest true "New data for question"
// @Success 200 {object} CourseItemUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items [put]
func Delete(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
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
	// If not admin, garant, or tutor
	if userRole != enums.CourseUserRoleGarant && userRole != enums.CourseUserRoleTutor {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	transaction := initializers.DB.Begin()

	courseItemService := services_course_item.CourseItemService{}
	courseItem, err := courseItemService.GetCourseItemByID(transaction, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}

	if !courseItem.Editable {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	// Get child count
	var childCount int64
	if err := transaction.
		Model(&models.CourseItem{}).
		Where("parent_id = ?", courseItem.ID).
		Count(&childCount).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get child count",
		}
	}
	if childCount != 0 {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "Failed to delete course item",
			Details: "Item has child objects",
		}
	}

	// Get term count
	var termCount int64
	if err := transaction.
		Model(&models.Term{}).
		Where("course_item_id = ?", courseItem.ID).
		Count(&termCount).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get child count",
		}
	}
	if termCount != 0 {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "Failed to delete course item",
			Details: "Item has terms",
		}
	}

	if err := transaction.Delete(&courseItem).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update course item",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	c.JSON(200, CourseItemDeleteResponse{
		Success: true,
	})

	return nil
}
