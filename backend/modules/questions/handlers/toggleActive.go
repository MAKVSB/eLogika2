package handlers

import (
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Newly created question
type QuestionToggleActiveResponse struct {
	Active bool `json:"active"`
}

type CourseToggleActiveUri struct {
	CourseID   uint `uri:"courseId" binding:"required"`
	QuestionID uint `uri:"questionId" binding:"required"`
}

// @Summary Activate/Deactivate question
// @Tags Questions
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param questionId path int true "ID of the updated question"
// @Success 200 {object} QuestionToggleActiveResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/questions/{questionId}/toggleActive [patch]
func QuestionToggleActive(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
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
			Code:    404,
			Message: "Failed to load question",
		}
	}

	// Update only selected values
	question.Active = !question.Active
	question.UpdatedByID = userData.ID
	question.UpdatedAt = time.Now()

	if err := transaction.Save(&question).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update question",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	if err := initializers.DB.
		First(&question, question.ID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch updated data",
		}
	}

	c.JSON(200, QuestionToggleActiveResponse{
		Active: question.Active,
	})
	return nil
}
