package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/templates/dtos"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TemplateBlockSegmentUpdateRequest struct {
	ID            uint                     `json:"id" binding:"required"`
	ChapterID     *uint                    `json:"chapterId"`
	CategoryID    *uint                    `json:"categoryId"`
	QuestionCount uint                     `json:"questionCount" binding:"required"`
	StepsMode     *enums.StepSelectionEnum `json:"stepsMode"`
	FilterBy      enums.CategoryFilterEnum `json:"filterBy" binding:"required"`
	Deleted       bool                     `json:"deleted"`

	Steps     []uint `json:"steps" binding:"required"`
	Questions []uint `json:"questions" binding:"required"`
}

type TemplateBlockModelUpdateRequest struct {
	ID                    uint                         `json:"id" binding:"required"`
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
	AllowEmptyAnswers     bool                         `json:"allowEmptyAnswers"`

	Segments []TemplateBlockSegmentUpdateRequest `json:"segments" binding:"required"`
}

// @Description Request to update template
type TemplateUpdateRequest struct {
	Version       uint                              `json:"version" binding:"required"`
	Title         string                            `json:"title" binding:"required"`
	Description   string                            `json:"description"`
	MixBlocks     bool                              `json:"mixBlocks"`
	MixEverything bool                              `json:"mixEverything"`
	Blocks        []TemplateBlockModelUpdateRequest `json:"blocks" binding:"required"`
}

// @Description Newly created template
type TemplateUpdateResponse struct {
	Data dtos.TemplateDTO `json:"data"`
}

// @Summary Create new template
// @Tags Categories
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TemplateUpdateRequest true "New data for template"
// @Success 200 {object} TemplateUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/templates [put]
func Update(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID   uint `uri:"courseId" binding:"required"`
			TemplateID uint `uri:"templateId" binding:"required"`
		},
		TemplateUpdateRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

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

	// Check question counts and block sums are valid
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

	templateService := services.TemplateService{}
	template, err := templateService.GetTemplateByID(initializers.DB, params.CourseID, params.TemplateID, userData.ID, userRole, nil, true, &reqData.Version)
	if err != nil {
		return err
	}

	template.Version = template.Version + 1
	template.Title = reqData.Title
	template.Description = reqData.Description
	template.MixBlocks = reqData.MixBlocks
	template.MixEverything = reqData.MixEverything

	transaction := initializers.DB.Begin()

	template.Blocks = make([]models.TemplateBlock, len(reqData.Blocks))
	blockIds := make([]uint, 0)
	for i, b := range reqData.Blocks {
		if b.ID != 0 {
			if err := transaction.
				First(&template.Blocks[i], b.ID).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to fetch template block data",
				}
			}
		} else {
			template.Blocks[i] = models.TemplateBlock{
				ID:         0,
				TemplateID: template.ID,
			}
		}

		template.Blocks[i].Title = b.Title
		template.Blocks[i].ShowName = b.ShowName
		template.Blocks[i].DifficultyFrom = b.DifficultyFrom
		template.Blocks[i].DifficultyTo = b.DifficultyTo
		template.Blocks[i].Weight = b.Weight
		template.Blocks[i].QuestionFormat = b.QuestionFormat
		template.Blocks[i].QuestionCount = b.QuestionCount
		template.Blocks[i].AnswerCount = b.AnswerCount
		template.Blocks[i].AnswerDistribution = b.AnswerDistribution
		template.Blocks[i].WrongAnswerPercentage = b.WrongAnswerPercentage
		template.Blocks[i].MixInsideBlock = b.MixInsideBlock
		template.Blocks[i].AllowEmptyAnswers = b.AllowEmptyAnswers

		if err := transaction.Save(&template.Blocks[i]).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to update template block",
			}
		}

		template.Blocks[i].Segments = make([]models.TemplateBlockSegment, len(b.Segments))
		chapterIds := make([]uint, 0)
		for j, seg := range b.Segments {
			if seg.ID != 0 {
				if err := transaction.
					First(&template.Blocks[i].Segments[j], seg.ID).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to fetch template block segment data",
					}
				}
			} else {
				template.Blocks[i].Segments[j] = models.TemplateBlockSegment{
					ID:              0,
					TemplateBlockID: template.Blocks[i].ID,
					ChapterID:       seg.ChapterID,
				}
			}

			template.Blocks[i].Segments[j].QuestionCount = seg.QuestionCount
			template.Blocks[i].Segments[j].ChapterID = seg.ChapterID
			// TODO check that category is a part of chapter
			template.Blocks[i].Segments[j].CategoryID = seg.CategoryID
			// TODO check that steps is a part of chapter
			template.Blocks[i].Segments[j].StepsMode = seg.StepsMode
			template.Blocks[i].Segments[j].FilterBy = seg.FilterBy

			if err := transaction.Save(&template.Blocks[i].Segments[j]).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to update template segment",
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

			chapterIds = append(chapterIds, template.Blocks[i].Segments[j].ID)
		}
		{
			var itemsToDelete []uint
			dq := transaction.
				Model(models.TemplateBlockSegment{}).
				Where("template_block_id = ?", b.ID)

			if len(chapterIds) != 0 {
				dq = dq.Where("id NOT IN (?)", chapterIds)
			}

			if err := dq.Pluck("id", &itemsToDelete).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to fetch segments to delete data",
				}
			}

			for _, id := range itemsToDelete {
				if err := transaction.
					Delete(&models.TemplateBlockSegment{}, id).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to delete segment",
					}
				}
			}
		}

		blockIds = append(blockIds, template.Blocks[i].ID)
	}
	{
		var itemsToDelete []uint
		dq := transaction.
			Model(models.TemplateBlock{}).
			Where("template_id = ?", params.TemplateID)

		if len(blockIds) != 0 {
			dq = dq.Where("id NOT IN (?)", blockIds)
		}

		if err := dq.Pluck("id", &itemsToDelete).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch blocks to delete data",
			}
		}

		for _, id := range itemsToDelete {
			if err := transaction.
				Delete(&models.TemplateBlock{}, id).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to delete block",
				}
			}
		}
	}

	if err := transaction.Save(&template).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update template",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	template, err = templateService.GetTemplateByID(initializers.DB, params.CourseID, params.TemplateID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, TemplateUpdateResponse{
		Data: dtos.TemplateDTO{}.From(template),
	})

	return nil
}
