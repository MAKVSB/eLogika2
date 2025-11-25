package handlers

import (
	"os"
	"path/filepath"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/print/helpers"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// @Summary Print questions
// @Tags Print
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param questionId path int false "ID of the corresponding course"
// @Success 200 {file} file "PDF file of tests"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/print/question/{questionId} [post]
func PrintQuestion(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID   uint  `uri:"courseId" binding:"required"`
			QuestionID *uint `uri:"questionId"`
		},
		any,
	](c, "search")
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	var questions []*models.Question
	if params.QuestionID != nil {
		questionServ := services.NewQuestionService(repositories.NewQuestionRepository())
		modifier := func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("CourseLink.Chapter").
				Preload("CourseLink.Category").
				Preload("Answers")
		}
		question, err := questionServ.GetQuestionByID(initializers.DB, params.CourseID, *params.QuestionID, userData.ID, userRole, &modifier, true, nil)
		if err != nil {
			return err
		}
		questions = []*models.Question{question}
	} else {
		searchParams.Pagination = nil
		questionServ := services.NewQuestionService(repositories.NewQuestionRepository())
		modifier := func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("CourseLink.Chapter").
				Preload("CourseLink.Category").
				Preload("Answers")
		}
		questions, _, err = questionServ.ListQuestions(initializers.DB, params.CourseID, userData.ID, userRole, &modifier, true, searchParams)
		if err != nil {
			return err
		}
	}

	workDir, err2 := os.Getwd()
	if err2 != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to print questions",
			Details: err2.Error(),
		}
	}

	tmpFolder, err2 := utils.CreateTmpFolder(workDir)
	if err2 != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to print questions",
			Details: err2.Error(),
		}
	}

	testOutputDir, err2 := utils.CreateFolder(tmpFolder)
	if err2 != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to print questions",
			Details: err2.Error(),
		}
	}

	assetDir, err2 := utils.CreateFolder(filepath.Join(tmpFolder, "assets"))
	if err2 != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to print questions",
			Details: err2.Error(),
		}
	}

	filepath, err2 := helpers.PrintQuestions(questions, workDir, assetDir, testOutputDir)
	if err2 != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to print questions",
			Details: err2.Error(),
		}
	}
	c.FileAttachment(filepath, uuid.NewString())

	return nil
}
