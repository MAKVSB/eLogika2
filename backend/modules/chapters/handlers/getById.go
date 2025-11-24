package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/chapters/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Newly created chapter
type ChapterGetByIdResponse struct {
	Data dtos.ChapterDTO `json:"data"`
}

// @Summary Get chapter by id
// @Tags Chapters
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param chapterId path int true "ID of the edited chapter"
// @Success 200 {object} ChapterUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/chapters/{chapterId} [get]
func ChapterGetByID(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID  uint `uri:"courseId" binding:"required"`
			ChapterID uint `uri:"chapterId" binding:"required"`
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

	chapterService := services.NewChapterService(repositories.NewChapterRepository())
	chapter, err := chapterService.GetChapterByID(initializers.DB, params.CourseID, params.ChapterID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, ChapterGetByIdResponse{
		Data: dtos.ChapterDTO{}.From(chapter),
	})
	return nil
}
