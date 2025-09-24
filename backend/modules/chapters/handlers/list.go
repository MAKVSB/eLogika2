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
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type ChapterListResponse struct {
	Items      []dtos.ChapterListItemDTO `json:"items"`
	ItemsCount int64                     `json:"itemsCount"`
}

type ChapterListRequest struct {
	common.SearchRequest
}

// @Summary List all available chapters in course
// @Tags Chapters
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body ChapterListRequest true "Ability to filter results"
// @Success 200 {object} ChapterListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/chapters [get]
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

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}

	chapterRepo := repositories.NewChapterRepository()

	var chapters []*models.Chapter
	var chapterCount int64

	// If not admin, garant, or tutor
	if userRole == enums.CourseUserRoleAdmin || userRole == enums.CourseUserRoleGarant || userRole == enums.CourseUserRoleTutor {
		chapters, chapterCount, err = chapterRepo.ListChapters(initializers.DB, params.CourseID, nil, false, searchParams)
		if err != nil {
			return err
		}
	} else if userRole == enums.CourseUserRoleStudent {
		chapters, chapterCount, err = chapterRepo.ListChaptersStudent(initializers.DB, params.CourseID, nil, false, searchParams)
		if err != nil {
			return err
		}
	} else {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.ChapterListItemDTO, len(chapters))
	for i, q := range chapters {
		dtoList[i] = dtos.ChapterListItemDTO{}.From(q)
	}

	c.JSON(200, ChapterListResponse{
		Items:      dtoList,
		ItemsCount: chapterCount,
	})
	return nil
}
