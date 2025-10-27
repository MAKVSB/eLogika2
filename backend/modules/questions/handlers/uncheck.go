package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/questions/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Marks that question correctness has not been checked by invocating user
// @Tags Questions
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param questionId path int true "ID of the unchecked question"
// @Success 200 {object} QuestionCheckResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/questions/{questionId}/check [delete]
func Uncheck(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID   uint `uri:"courseId" binding:"required"`
			QuestionID uint `uri:"questionId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	transaction := initializers.DB.Begin()

	query := transaction

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		transaction.Rollback()
		return err
	}
	// If not admin, garant, or tutor
	if userRole == enums.CourseUserRoleAdmin {
		// TODO? Admin can get anything
	} else if userRole == enums.CourseUserRoleGarant {
		// Filter questions that are created by any garant
		courseRepo := repositories.CourseRepository{}
		garants, err := courseRepo.GetCourseGarantsIds(initializers.DB, params.CourseID)
		if err != nil {
			transaction.Rollback()
			return err
		}
		query = query.Where("created_by_id in ?", garants)
		// And has ManagedBy GARANT
		query = query.Where("managed_by = ?", enums.CourseUserRoleGarant)

	} else if userRole == enums.CourseUserRoleTutor {
		// Filter questions created by tutor
		query = query.Where("created_by_id = ?", userData.ID)
		// And has ManagedBy TUTOR
		query = query.Where("managed_by = ?", enums.CourseUserRoleTutor)
	} else {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	var question models.Question
	if err := query.First(&question, params.QuestionID).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to load question",
		}
	}

	if err := transaction.Where("question_id = ? AND user_id = ?", params.QuestionID, userData.ID).Delete(&models.QuestionCheck{}).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "failed to uncheck question",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	var questionChecks []models.QuestionCheck
	if err := initializers.DB.Preload("User").
		Where("question_id = ?", params.QuestionID).
		Find(&questionChecks).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to read updated data",
		}
	}

	checkedBy := make([]dtos.QuestionCheckedByDTO, len(questionChecks))
	for i, userCheck := range questionChecks {
		checkedBy[i] = dtos.QuestionCheckedByDTO{}.From(&userCheck)
	}

	c.JSON(200, QuestionCheckResponse{
		CheckedBy: checkedBy,
	})
	return nil
}
