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

// @Description Newly created courseItem
type CourseItemSelectResultResponse struct {
	Success bool `json:"success"`
}

// @Summary Get courseItem by id
// @Tags CourseItems
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param courseItemId path int true "ID of the requested item"
// @Success 200 {object} CourseItemListResultsResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId}/results/{resultId} [put]
func SelectResult(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			ResultID     uint `uri:"resultId" binding:"required"`
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

	courseItemService := services_course_item.CourseItemService{}
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	var nowActiveResult *models.CourseItemResult
	if err := transaction.
		Where("course_item_id = ?", courseItem.ID).
		First(&nowActiveResult, params.ResultID).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to find result",
			Details: err.Error(),
		}
	}

	if nowActiveResult.Selected {
		nowActiveResult.Selected = false
	} else {
		if err := transaction.
			Model(&models.CourseItemResult{}).
			Where("course_item_id = ?", params.CourseItemID).
			Where("student_id = ?", nowActiveResult.StudentID).
			Update("selected", false).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to unselect all other results",
				Details: err.Error(),
			}
		}

		nowActiveResult.Selected = true
	}

	if err := transaction.
		Save(&nowActiveResult).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to select result",
			Details: err.Error(),
		}
	}

	if err := transaction.
		Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
			Details: err.Error(),
		}
	}

	c.JSON(200, CourseItemSelectResultResponse{
		Success: true,
	})

	return nil
}
