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

// @Description Newly created chapter
type ChapterMoveResponse struct {
	Childs []dtos.ChapterListItemDTO `json:"childs"`
}

// @Summary Move chapter order
// @Description Changes the ordering in relation to its parent
// @Tags Chapters
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param chapterId path int true "ID of the parent chapter"
// @Param direction path enums.MoveDirectionEnum true "Direction of the move"
// @Success 200 {object} ChapterMoveResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/chapters/{chapterId}/move/{direction} [patch]
func ChapterMove(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID  uint                    `uri:"courseId" binding:"required"`
			ChapterID uint                    `uri:"chapterId" binding:"required"`
			Direction enums.MoveDirectionEnum `uri:"direction" binding:"required"`
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
	// If not admin, garant, or tutor
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant && userRole != enums.CourseUserRoleTutor {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	transaction := initializers.DB.Begin()

	chapterRepo := repositories.NewChapterRepository()
	movingChapter, err := chapterRepo.GetChapterByID(transaction, params.CourseID, params.ChapterID, false, nil)
	if err != nil {
		transaction.Rollback()
		return err
	}
	swapWithChapter, err := chapterRepo.GetNextInOrder(transaction, params.CourseID, *movingChapter.ParentID, params.Direction, movingChapter.Order)
	if err != nil {
		transaction.Rollback()
		return err
	}

	swapWithChapterOrder := swapWithChapter.Order
	movingChapterOrder := movingChapter.Order

	if err := transaction.
		Model(&swapWithChapter).
		Update("order", movingChapterOrder).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update chapter (1)",
		}
	}

	if err := transaction.
		Model(&movingChapter).
		Update("order", swapWithChapterOrder).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update chapter (2)",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	var resultChapters []models.Chapter
	if err := initializers.DB.Where("parent_id", movingChapter.ParentID).Find(&resultChapters).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch chapter order",
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.ChapterListItemDTO, len(resultChapters))
	for i, c := range resultChapters {
		dtoList[i] = dtos.ChapterListItemDTO{}.From(&c)
	}

	c.JSON(200, ChapterMoveResponse{
		Childs: dtoList,
	})
	return nil
}
