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
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"elogika.vsb.cz/backend/utils/tiptap"
	"github.com/gin-gonic/gin"
)

// @Description Request to update question
type QuestionUpdateRequest struct {
	Title              string                        `json:"title" binding:"required" example:"Is number even ?"`    // Title of the question (for listing only)
	Content            *models.TipTapContent         `json:"content" binding:"required" ts_type:"JSONContent"`       // Question text in json (Using TipTap editor format)
	TimeToRead         int                           `json:"timeToRead"`                                             // Estimated time in seconds it takes a user to read the question.
	TimeToProcess      int                           `json:"timeToProcess"`                                          // Estimated time in seconds it takes to think about solution and evaluate common parts of solution (drawing graphs and other)
	QuestionType       enums.QuestionTypeEnum        `json:"questionType" binding:"required"`                        // Type of the question
	QuestionFormat     enums.QuestionFormatEnum      `json:"questionFormat" binding:"required"`                      // Format of the question
	IncludeAnswerSpace bool                          `json:"includeAnswerSpace"`                                     // Defines if a box of empty space should be included after open question
	Active             bool                          `json:"active"`                                                 // Is the question in active pool for selection
	Answers            []dtos.QuestionAnswerAdminDTO `json:"answers"`                                                // All answers for this question
	ChapterID          uint                          `json:"chapterId" binding:"required"`                           // ID of the chapter
	CategoryID         *uint                         `json:"categoryId" validate:"optional" ts_type:"number | null"` // ID of the category
	Steps              []uint                        `json:"steps"`                                                  // Steps required for answering question
	Version            uint                          `json:"version"`                                                // Version signature to prevent concurrency problems
	AsNewVersion       bool                          `json:"asNewVersion"`                                           // Indicates if this version of question should be edited or inserted as new version. (If false, can result in modifying already generated tests)
}

// @Description Newly created question
type QuestionUpdateResponse struct {
	Data dtos.QuestionAdminDTO `json:"data"`
}

type QuestionUpdateParams struct {
	CourseID   uint `uri:"courseId" binding:"required"`
	QuestionID uint `uri:"questionId" binding:"required"`
}

// @Summary Modify question
// @Tags Questions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param questionId path int true "ID of the updated question"
// @Param body body QuestionUpdateRequest true "New data for question"
// @Success 200 {object} QuestionUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/questions/{questionId} [put]
func QuestionUpdate(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		QuestionUpdateParams,
		QuestionUpdateRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	if reqData.AsNewVersion {
		return asNewVersion(c, userData, reqData, params, userRole)
	} else {
		return updateExisting(c, userData, reqData, params, userRole)
	}
}

//  TODO TODO TODO TODOO tady jsem skončil s předěláváním

func asNewVersion(c *gin.Context, userData authdtos.LoggedUserDTO, reqData *QuestionUpdateRequest, params *QuestionUpdateParams, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	questionService := services.QuestionService{}
	questionRepo := repositories.QuestionRepository{}
	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	// Get question
	oldQuestion, err := questionService.GetQuestionByID(initializers.DB, params.CourseID, params.QuestionID, userData.ID, userRole, nil, true, &reqData.Version, false, false)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	maxVersion, err := questionRepo.GetMaxVersion(transaction, oldQuestion.QuestionGroupID)
	if err != nil {
		transaction.Rollback()
		return err
	}

	newQuestion := &models.Question{
		ID:                 0,
		Version:            maxVersion + 1,
		Title:              reqData.Title,
		Content:            reqData.Content,
		TimeToRead:         reqData.TimeToRead,
		TimeToProcess:      reqData.TimeToProcess,
		QuestionType:       reqData.QuestionType,
		QuestionFormat:     reqData.QuestionFormat,
		IncludeAnswerSpace: reqData.IncludeAnswerSpace,
		CreatedByID:        userData.ID,
		ManagedBy:          oldQuestion.ManagedBy,
		Active:             reqData.Active,
		QuestionGroupID:    oldQuestion.QuestionGroupID,
		AnswerCount:        uint(len(reqData.Answers)),
	}

	err = tiptap.FindAndSaveRelations(transaction, userData.ID, reqData.Content, &newQuestion, "ContentFiles")
	if err != nil {
		return err
	}

	if err := transaction.Save(&newQuestion).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert question",
			Details: err.Error(),
		}
	}

	// Unlink question from course
	if err := transaction.Model(&models.CourseQuestion{}).
		Where("course_id = ?", params.CourseID).
		Where("question_id = ?", params.QuestionID).
		Delete(&models.CourseQuestion{}).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to unlink old question from course",
			Details: err.Error(),
		}
	}

	// Link new question instance to course
	new_course_question := &models.CourseQuestion{
		ID:         0,
		Version:    maxVersion + 1,
		CourseID:   params.CourseID,
		QuestionID: newQuestion.ID,
		ChapterID:  reqData.ChapterID,
		CategoryID: reqData.CategoryID,
	}

	if err := transaction.Save(&new_course_question).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to link new question to course",
			Details: err.Error(),
		}
	}

	newQuestion.CourseLink = new_course_question

	// Sync steps
	newQuestion, err = questionRepo.SyncSteps(transaction, newQuestion, reqData.CategoryID, reqData.Steps)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// sync answers
	err = questionService.SyncAnswers(transaction, userData.ID, newQuestion, reqData.Answers, true)
	if err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	// Fetch updated data
	newQuestion, err = questionService.GetQuestionByID(initializers.DB, params.CourseID, newQuestion.ID, userData.ID, userRole, nil, true, nil, false, true)
	if err != nil {
		return err
	}

	c.JSON(200, QuestionUpdateResponse{
		Data: dtos.QuestionAdminDTO{}.From(newQuestion),
	})
	return nil
}

func updateExisting(c *gin.Context, userData authdtos.LoggedUserDTO, reqData *QuestionUpdateRequest, params *QuestionUpdateParams, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	questionService := services.QuestionService{}
	questionRepo := repositories.QuestionRepository{}
	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	// Get question
	question, err := questionService.GetQuestionByID(initializers.DB, params.CourseID, params.QuestionID, userData.ID, userRole, nil, true, &reqData.Version, true, false)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	maxVersion, err := questionRepo.GetMaxVersion(transaction, question.QuestionGroupID)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Update only selected values
	question.Version = maxVersion + 1
	question.Title = reqData.Title
	err = tiptap.FindAndSaveRelations(transaction, userData.ID, reqData.Content, &question, "ContentFiles")
	if err != nil {
		return err
	}
	question.Content = reqData.Content
	question.TimeToRead = reqData.TimeToRead
	question.TimeToProcess = reqData.TimeToProcess
	question.QuestionType = reqData.QuestionType
	question.QuestionFormat = reqData.QuestionFormat
	question.IncludeAnswerSpace = reqData.IncludeAnswerSpace
	question.Active = reqData.Active
	question.AnswerCount = uint(len(reqData.Answers))

	if err := transaction.Save(&question).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update question",
			Details: err.Error(),
		}
	}

	question.CourseLink.ChapterID = reqData.ChapterID
	question.CourseLink.CategoryID = reqData.CategoryID

	if err := transaction.Save(&question.CourseLink).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update question link",
			Details: err.Error(),
		}
	}

	// Sync steps
	question, err = questionRepo.SyncSteps(transaction, question, reqData.CategoryID, reqData.Steps)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// sync answers
	err = questionService.SyncAnswers(transaction, userData.ID, question, reqData.Answers, false)
	if err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	// Fetch updated data
	question, err = questionService.GetQuestionByID(initializers.DB, params.CourseID, params.QuestionID, userData.ID, userRole, nil, true, nil, true, true)
	if err != nil {
		return err
	}

	c.JSON(200, QuestionUpdateResponse{
		Data: dtos.QuestionAdminDTO{}.From(question),
	})
	return nil
}
