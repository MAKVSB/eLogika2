package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/questions/dtos"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Newly created question
type QuestionGetByIdResponse struct {
	Data dtos.QuestionAdminDTO `json:"data"`
}

// @Summary Get question by id
// @Tags Questions
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param questionId path int true "ID of the requested question"
// @Success 200 {object} QuestionGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/questions/{questionId} [get]
func QuestionGetByID(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
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

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	questionService := services.QuestionService{}
	question, err := questionService.GetQuestionByID(initializers.DB, params.CourseID, params.QuestionID, userData.ID, userRole, nil, true, nil, true, true)
	if err != nil {
		return err
	}

	c.JSON(200, QuestionGetByIdResponse{
		Data: dtos.QuestionAdminDTO{}.From(question),
	})
	return nil
}
