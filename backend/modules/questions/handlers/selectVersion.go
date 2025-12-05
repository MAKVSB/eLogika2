package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/questions/dtos"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to update question

// @Description Newly created question
type QuestionSelectVersionResponse struct {
	Data dtos.QuestionAdminDTO `json:"data"`
}

// @Summary Select question version as active
// @Tags Questions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param questionId path int true "ID of the currently active question version"
// @Param newVersionQuestionId path int true "ID of the new question version"
// @Success 200 {object} QuestionUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/questions/{oldVersionQuestionId}/selecversion [patch]
func SelectVersion(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID   uint `uri:"courseId" binding:"required"`
			QuestionId uint `uri:"questionId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	questionService := services.QuestionService{}
	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	newQuestion, err := questionService.GetQuestionByID(initializers.DB, params.CourseID, params.QuestionId, userData.ID, userRole, nil, true, nil, true, false)
	if err != nil {
		return err
	}

	var versions []*models.CourseQuestion
	if err := transaction.
		Select("course_questions.id").
		Where("course_id = ?", params.CourseID).
		InnerJoins("Question", initializers.DB.Select("id").Where("question_group_id = ?", newQuestion.QuestionGroupID)).
		Find(&versions).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to unlink old question from course",
			Details: err.Error(),
		}
	}

	if err := transaction.Delete(&versions).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to unlink old question from course",
			Details: err.Error(),
		}
	}

	// Get question
	newQuestion.CourseLink.DeletedAt.Valid = false

	// Link new question instance to course
	if err := transaction.Save(&newQuestion.CourseLink).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to link new question to course",
			Details: err.Error(),
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	// Fetch updated data
	newQuestion, err = questionService.GetQuestionByID(initializers.DB, params.CourseID, params.QuestionId, userData.ID, userRole, nil, true, nil, false, true)
	if err != nil {
		return err
	}

	c.JSON(200, QuestionGetByIdResponse{
		Data: dtos.QuestionAdminDTO{}.From(newQuestion),
	})
	return nil
}
