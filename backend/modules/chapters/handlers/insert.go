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
type ChapterInsertRequest struct {
	Name    string                `json:"name" binding:"required"`                          // Name of the chapter
	Content *models.TipTapContent `json:"content" binding:"required" ts_type:"JSONContent"` // Content of chapter
	Visible bool                  `json:"visible"`                                          // Should chapter be visible to students
}

// @Description Newly created chapter
type ChapterInsertResponse struct {
	Data dtos.ChapterDTO `json:"data"`
}

// @Summary Create new chapter
// @Description Adds new chapter as a child of `parentId`
// @Tags Chapters
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param parentId path int true "ID of the parent chapter"
// @Param body body ChapterInsertRequest true "New data for chapter"
// @Success 200 {object} ChapterInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/chapters/{parentId} [post]
func ChapterInsert(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
			ParentID uint `uri:"parentId" binding:"required"`
		},
		ChapterInsertRequest,
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

	chapterRepo := repositories.NewChapterRepository()
	lastOrder := chapterRepo.LastChapterOrder(transaction, params.ParentID)

	chapter := &models.Chapter{
		ID:       0,
		Version:  1,
		CourseID: params.CourseID,
		Name:     reqData.Name,
		Content:  reqData.Content,
		Visible:  reqData.Visible,
		ParentID: &params.ParentID,
		Order:    lastOrder + 1,
	}

	err = tiptap.FindAndSaveRelations(transaction, userData.ID, reqData.Content, &chapter, "ContentFiles")
	if err != nil {
		return err
	}
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

	chapterService := services.NewChapterService(chapterRepo)
	chapter, err = chapterService.GetChapterByID(transaction, params.CourseID, chapter.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, ChapterInsertResponse{
		Data: dtos.ChapterDTO{}.From(chapter),
	})

	return nil
}
