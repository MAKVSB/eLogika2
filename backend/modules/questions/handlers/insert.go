package handlers

import (
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/questions/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"elogika.vsb.cz/backend/utils/tiptap"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Description Request to insert new question
type QuestionInsertRequest struct {
	Title          string                        `json:"title" binding:"required" example:"Is number even ?"`    // Title of the question (for listing only)
	Content        *models.TipTapContent         `json:"content" binding:"required" ts_type:"JSONContent"`       // Question text in json (Using TipTap editor format)
	TimeToRead     int                           `json:"timeToRead"`                                             // Estimated time in seconds it takes a user to read the question.
	TimeToProcess  int                           `json:"timeToProcess"`                                          // Estimated time in seconds it takes to think about solution and evaluate common parts of solution (drawing graphs and other)
	QuestionType   enums.QuestionTypeEnum        `json:"questionType" binding:"required"`                        // Type of the question
	QuestionFormat enums.QuestionFormatEnum      `json:"questionFormat" binding:"required"`                      // Format of the question
	Active         bool                          `json:"active"`                                                 // Is the question in active pool for selection
	Answers        []dtos.QuestionAnswerAdminDTO `json:"answers"`                                                // All answers for this question
	ChapterID      uint                          `json:"chapterId" binding:"required"`                           // ID of the chapter
	CategoryID     *uint                         `json:"categoryId" validate:"optional" ts_type:"number | null"` // ID of the category
	Steps          []uint                        `json:"steps"`                                                  // Steps required for answering question
}

// @Description Newly created question
type QuestionInsertResponse struct {
	Data dtos.QuestionAdminDTO `json:"data"`
}

// @Summary Create new question
// @Tags Questions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body QuestionInsertRequest true "New data for question"
// @Success 200 {object} QuestionInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/questions [post]
func QuestionInsert(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		QuestionInsertRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	questionService := services.QuestionService{}
	questionRepo := repositories.QuestionRepository{}

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

	//Create question group
	questionGroup := models.QuestionGroup{
		ID:           0,
		OriginalName: reqData.Title,
	}

	transaction := initializers.DB.Begin()

	if err := transaction.Save(&questionGroup).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to create question group",
		}
	}

	//Create question
	question := &models.Question{
		ID:              0,
		Version:         1,
		Title:           reqData.Title,
		Content:         reqData.Content,
		TimeToRead:      reqData.TimeToRead,
		TimeToProcess:   reqData.TimeToProcess,
		QuestionType:    reqData.QuestionType,
		QuestionFormat:  reqData.QuestionFormat,
		CreatedAt:       time.Now(),
		CreatedByID:     userData.ID,
		UpdatedAt:       time.Now(),
		UpdatedByID:     userData.ID,
		ManagedBy:       userRole,
		Active:          reqData.Active,
		AnswerCount:     uint(len(reqData.Answers)),
		QuestionGroupID: questionGroup.ID,
	}

	err = tiptap.FindAndSaveRelations(transaction, userData.ID, reqData.Content, &question, "ContentFiles")
	if err != nil {
		return err
	}

	if err := transaction.Save(&question).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert question",
		}
	}

	// Link question to course
	question.CourseLink = &models.CourseQuestion{
		Version:    1,
		CourseID:   params.CourseID,
		QuestionID: question.ID,
		ChapterID:  reqData.ChapterID,
		CategoryID: reqData.CategoryID,
	}

	if err := transaction.Save(&question.CourseLink).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to connect question to course",
		}
	}

	// Sync steps
	question, err = questionRepo.SyncSteps(transaction, question, reqData.CategoryID, reqData.Steps)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// sync answers
	err = questionService.SyncAnswers(transaction, userData.ID, question, reqData.Answers)
	if err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.
		Joins("CreatedBy").
		Preload("CheckedBy").
		Preload("CheckedBy.User").
		Preload("Answers").
		Preload("Answers.Answer").
		Preload("CourseLink", func(db *gorm.DB) *gorm.DB {
			return db.Where("course_id = ?", params.CourseID)
		}).
		Preload("CourseLink.Steps").
		First(&question, question.ID).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch updated data",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	c.JSON(200, QuestionInsertResponse{
		Data: dtos.QuestionAdminDTO{}.From(question),
	})
	return nil
}
