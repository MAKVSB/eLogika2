package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/templates/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TemplateBlockSegmentInsertRequest struct {
	ChapterID     *uint                    `json:"chapterId"`
	CategoryID    *uint                    `json:"categoryId"`
	QuestionCount uint                     `json:"questionCount" binding:"required"`
	StepsMode     *enums.StepSelectionEnum `json:"stepsMode"`
	FilterBy      enums.CategoryFilterEnum `json:"filterBy" binding:"required"`

	Steps     []uint `json:"steps" binding:"required"`
	Questions []uint `json:"questions" binding:"required"`
}

type TemplateBlockInsertRequest struct {
	Title                 string                       `json:"title" binding:"required"`
	ShowName              bool                         `json:"showName"`
	DifficultyFrom        uint                         `json:"difficultyFrom" binding:"required"`
	DifficultyTo          uint                         `json:"difficultyTo" binding:"required"`
	Weight                uint                         `json:"weight" binding:"required"`
	QuestionFormat        enums.QuestionFormatEnum     `json:"questionFormat" binding:"required"`
	QuestionCount         uint                         `json:"questionCount" binding:"required"`
	AnswerCount           uint                         `json:"answerCount" binding:"required"`
	AnswerDistribution    enums.AnswerDistributionEnum `json:"answerDistribution" binding:"required"`
	WrongAnswerPercentage uint                         `json:"wrongAnswerPercentage" binding:"required"`
	MixInsideBlock        bool                         `json:"mixInsideBlock"`

	Segments []TemplateBlockSegmentInsertRequest `json:"segments" binding:"required"`
}

// @Description Request to insert new template
type TemplateInsertRequest struct {
	Title         string                       `json:"title" binding:"required"`
	Description   string                       `json:"description"`
	MixBlocks     bool                         `json:"mixBlocks"`
	MixEverything bool                         `json:"mixEverything"`
	Blocks        []TemplateBlockInsertRequest `json:"blocks" binding:"required"`
}

// @Description Newly created template
type TemplateInsertResponse struct {
	Data dtos.TemplateDTO `json:"data"`
}

// @Summary Create new template
// @Tags Categories
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TemplateInsertRequest true "New data for template"
// @Success 200 {object} TemplateInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/templates [post]
func Insert(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		TemplateInsertRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin, garant, or tutor
	if userRole != enums.CourseUserRoleGarant && userRole != enums.CourseUserRoleTutor {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	// Check question counts are valid
	blockWeight := 0
	for _, block := range reqData.Blocks {
		blockWeight += int(block.Weight)
		if len(block.Segments) != 0 {
			segmentQuestionsSum := uint(0)
			for _, segment := range block.Segments {
				segmentQuestionsSum += segment.QuestionCount
			}
			if segmentQuestionsSum != block.QuestionCount {
				return &common.ErrorResponse{
					Code:    422,
					Message: "Question counts are not consistent",
				}
			}
		}
	}
	if blockWeight != 100 {
		return &common.ErrorResponse{
			Code:    422,
			Message: "Block weights must sum to 100%",
		}
	}

	template := &models.Template{
		ID:            0,
		Title:         reqData.Title,
		Description:   reqData.Description,
		MixBlocks:     reqData.MixBlocks,
		MixEverything: reqData.MixEverything,
		CreatedByID:   userData.ID,
		ManagedBy:     userRole,
		CourseID:      params.CourseID,
		Version:       1,
	}

	transaction := initializers.DB.Begin()

	if err := transaction.Save(&template).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert template",
		}
	}

	template.Blocks = make([]models.TemplateBlock, len(reqData.Blocks))
	for i, b := range reqData.Blocks {
		template.Blocks[i] = models.TemplateBlock{
			ID:                    0,
			Title:                 b.Title,
			ShowName:              b.ShowName,
			DifficultyFrom:        b.DifficultyFrom,
			DifficultyTo:          b.DifficultyTo,
			Weight:                b.Weight,
			QuestionFormat:        b.QuestionFormat,
			QuestionCount:         b.QuestionCount,
			AnswerCount:           b.AnswerCount,
			AnswerDistribution:    b.AnswerDistribution,
			WrongAnswerPercentage: b.WrongAnswerPercentage,
			MixInsideBlock:        b.MixInsideBlock,
			TemplateID:            template.ID,
		}
		if err := transaction.Save(&template.Blocks[i]).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to insert template block",
			}
		}

		template.Blocks[i].Segments = make([]models.TemplateBlockSegment, len(b.Segments))
		for j, seg := range b.Segments {
			template.Blocks[i].Segments[j] = models.TemplateBlockSegment{
				ID:              0,
				TemplateBlockID: template.Blocks[i].ID,
				QuestionCount:   b.QuestionCount,
				ChapterID:       seg.ChapterID,
				CategoryID:      seg.CategoryID,
				StepsMode:       seg.StepsMode,
				FilterBy:        seg.FilterBy,
			}

			if err := transaction.Save(&template.Blocks[i].Segments[j]).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to insert template block segment",
				}
			}

			// Sync steps
			var steps []models.Step
			if err := transaction.Model(&models.Step{}).Where("id IN ? AND category_id = ?", seg.Steps, seg.CategoryID).Find(&steps).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to load steps",
				}
			}
			template.Blocks[i].Segments[j].Steps = steps
			if err := transaction.Model(&template.Blocks[i].Segments[j]).Association("Steps").Replace(&template.Blocks[i].Segments[j].Steps); err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to update steps",
				}
			}

			// Sync questions
			var questions []models.QuestionGroup
			if err := transaction.Model(&models.QuestionGroup{}).Where("id IN ?", seg.Questions).Find(&questions).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to load questions",
				}
			}
			template.Blocks[i].Segments[j].Questions = questions
			if err := transaction.Model(&template.Blocks[i].Segments[j]).Association("Questions").Replace(&template.Blocks[i].Segments[j].Questions); err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to update questions",
				}
			}
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	templateServ := services.NewTemplateService(repositories.NewTemplateRepository())
	template, err = templateServ.GetTemplateByID(initializers.DB, params.CourseID, template.ID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, TemplateInsertResponse{
		Data: dtos.TemplateDTO{}.From(template),
	})
	return nil
}
