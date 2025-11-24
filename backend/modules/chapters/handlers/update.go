package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/chapters/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"elogika.vsb.cz/backend/utils/tiptap"
	"github.com/gin-gonic/gin"
)

// @Description Request to insert new chapter
type ChapterUpdateRequest struct {
	Name    string                `json:"name" binding:"required"`                          // Name of the chapter
	Content *models.TipTapContent `json:"content" binding:"required" ts_type:"JSONContent"` // Content of chapter
	Visible bool                  `json:"visible"`                                          // Should chapter be visible to students
	Version uint                  `json:"version"`                                          // Version signature to prevent concurrency problems
}

// @Description Updated chapter
type ChapterUpdateResponse struct {
	Data dtos.ChapterDTO `json:"data"`
}

type CourseChapterUri struct {
	CourseID  uint `uri:"courseId" binding:"required"`
	ChapterID uint `uri:"chapterId" binding:"required"`
}

// @Summary Update chapter
// @Description Updates chapter content
// @Tags Chapters
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param chapterId path int true "ID of the edited chapter"
// @Param body body ChapterUpdateRequest true "New data for chapter"
// @Success 200 {object} ChapterUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/chapters/{chapterId} [put]
func ChapterUpdate(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID  uint `uri:"courseId" binding:"required"`
			ChapterID uint `uri:"chapterId" binding:"required"`
		},
		ChapterUpdateRequest,
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
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant && userRole != enums.CourseUserRoleTutor {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	transaction := initializers.DB.Begin()

	chapterService := services.NewChapterService(repositories.NewChapterRepository())
	chapter, err := chapterService.GetChapterByID(transaction, params.CourseID, params.ChapterID, userRole, nil, true, &reqData.Version)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Partially modify data
	chapter.Version = chapter.Version + 1
	chapter.Name = reqData.Name
	err = tiptap.FindAndSaveRelations(transaction, userData.ID, reqData.Content, &chapter, "ContentFiles")
	if err != nil {
		return err
	}
	chapter.Content = reqData.Content
	chapter.Visible = reqData.Visible

	if err := transaction.Save(&chapter).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update chapter",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	chapter, err = chapterService.GetChapterByID(transaction, params.CourseID, params.ChapterID, userRole, nil, true, nil)
	if err != nil {
		transaction.Rollback()
		return err
	}

	c.JSON(200, ChapterUpdateResponse{
		Data: dtos.ChapterDTO{}.From(chapter),
	})
	return nil
}
