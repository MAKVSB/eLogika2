package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/questions/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type QuestionListResponse struct {
	Items      []dtos.QuestionListItemDTO `json:"items"`
	ItemsCount int64                      `json:"itemsCount"`
}

type QuestionListRequest struct {
	common.SearchRequest
}

// @Summary List all available questions in course
// @Tags Questions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body QuestionListRequest true "Ability to filter results"
// @Success 200 {object} QuestionListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/questions [get]
func List(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		any,
	](c, "search")
	if err != nil {
		return err
	}

	// TODO validate from here 2
	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}
	questionServ := services.NewQuestionService(repositories.NewQuestionRepository())
	questions, questionCount, err := questionServ.ListQuestions(initializers.DB, params.CourseID, userData.ID, userRole, nil, false, searchParams)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.QuestionListItemDTO, len(questions))
	for i, q := range questions {
		dtoList[i] = dtos.QuestionListItemDTO{}.From(q)
	}

	c.JSON(200, QuestionListResponse{
		Items:      dtoList,
		ItemsCount: questionCount,
	})
	return nil
}
