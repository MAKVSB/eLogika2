package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/categories/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Request to update category
type CategoryUpdateRequest struct {
	Name      string         `json:"name" binding:"required"`    // Name of the category
	Steps     []dtos.StepDTO `json:"steps"`                      // Steps
	ChapterID uint           `json:"chapterId"`                  // Chapter data
	Version   uint           `json:"version" binding:"required"` // Version signature to prevent concurrency problems
}

// @Description Newly created category
type CategoryUpdateResponse struct {
	Data dtos.CategoryDTO `json:"data"`
}

// @Summary Modify category
// @Tags Categories
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param categoryId path int true "ID of the updated category"
// @Param body body CategoryUpdateRequest true "New data for category"
// @Success 200 {object} CategoryUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/categories/{categoryId} [put]
func Update(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID   uint `uri:"courseId" binding:"required"`
			CategoryID uint `uri:"categoryId" binding:"required"`
		},
		CategoryUpdateRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin or garant
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	categoryRepo := repositories.CategoryRepository{}
	category, err := categoryRepo.GetCategoryByID(initializers.DB, params.CourseID, params.CategoryID, &reqData.Version)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	// Update only selected values
	category.Version = reqData.Version + 1
	category.Name = reqData.Name
	if category.ChapterID != reqData.ChapterID {
		category.ChapterID = reqData.ChapterID
		// Update all connected questions.
		if err := transaction.Model(models.Question{}).
			Where("category_id = ?", category.ID).
			Update("chapter_id", reqData.ChapterID).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to update connected questions",
			}
		}
	}

	if err := transaction.Save(&category).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update category",
		}
	}

	// sync steps
	for _, step := range reqData.Steps {
		if step.Deleted {
			if step.ID != 0 {
				if err := transaction.Delete(&models.Step{}, step.ID).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    500,
						Message: "Failed to delete answer",
					}
				}
			}
		} else {
			var stepData models.Step
			if step.ID == 0 {
				stepData = models.Step{
					ID:         0,
					Name:       step.Name,
					Difficulty: step.Difficulty,
					CategoryID: category.ID,
				}
			} else {
				if err := transaction.First(&stepData, step.ID).Error; err != nil {
					transaction.Rollback()
					return &common.ErrorResponse{
						Code:    404,
						Message: "Failed to load answer",
					}
				}
				stepData.Name = step.Name
				stepData.Difficulty = step.Difficulty
			}
			if err := transaction.Save(&stepData).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to create or update answer",
					Details: step,
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

	category, err = categoryRepo.GetCategoryByID(initializers.DB, params.CourseID, params.CategoryID, nil)
	if err != nil {
		return err
	}

	c.JSON(200, CategoryUpdateResponse{
		Data: dtos.CategoryDTO{}.From(category),
	})

	return nil
}
