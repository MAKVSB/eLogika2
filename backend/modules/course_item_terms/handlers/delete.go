package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_item_terms/dtos"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Get Term by id
type TermsDeleteResponse struct {
	Data dtos.TermDTO `json:"data"`
}

// @Summary Get term by id
// @Tags Terms
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param courseItemId path int true "ID of the corresponding course item"
// @Param termId path int true "ID of the requested item"
// @Success 200 {object} TermsGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId}/terms/{termId} [get]
func Delete(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			TermID       uint `uri:"termId" binding:"required"`
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

	transaction := initializers.DB.Begin()

	termService := services.TermService{}
	term, err := termService.GetTermByID(transaction, params.CourseID, params.CourseItemID, params.TermID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Check if students exists
	var studentsCount int64
	if err := transaction.
		Model(&models.UserTerm{}).
		Where("term_id = ?", params.TermID).
		Count(&studentsCount).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get number of joined students",
			Details: err.Error(),
		}
	}
	if studentsCount != 0 {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "Cannot delete term with students signed in",
		}
	}

	// Check generated tests
	var testsCount int64
	if err := transaction.
		Model(&models.Test{}).
		Where("term_id = ?", params.TermID).
		Count(&testsCount).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get number of generated tests",
			Details: err.Error(),
		}
	}
	if testsCount != 0 {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "Cannot delete term with test already generated",
		}
	}

	// Check activities
	var activityCount int64
	if err := transaction.
		Model(&models.ActivityInstance{}).
		Where("term_id = ?", params.TermID).
		Count(&activityCount).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get number of submitted activities",
			Details: err.Error(),
		}
	}
	if activityCount != 0 {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    409,
			Message: "Cannot delete term with test already generated",
		}
	}

	if err := transaction.Delete(&term).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to delete term",
			Details: err.Error(),
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

	c.JSON(200, TermsGetByIdResponse{
		Data: dtos.TermDTO{}.From(term),
	})
	return nil
}
