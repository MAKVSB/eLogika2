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

// @Description Request to insert new category
type CategoryInsertRequest struct {
	Name      string         `json:"name" binding:"required"`      // Name of the category
	Steps     []dtos.StepDTO `json:"steps"`                        // All steps for this category
	ChapterID uint           `json:"chapterId" binding:"required"` // Chapter data
}

// @Description Newly created category
type CategoryInsertResponse struct {
	Data dtos.CategoryDTO `json:"data"`
}

// @Summary Create new category
// @Tags Categories
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body CategoryInsertRequest true "New data for category"
// @Success 200 {object} CategoryInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/categories [post]
func Insert(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		CategoryInsertRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin or garant
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	transaction := initializers.DB.Begin()

	category := &models.Category{
		ID:        0,
		Version:   1,
		CourseID:  params.CourseID,
		Name:      reqData.Name,
		ChapterID: reqData.ChapterID,
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
		if !step.Deleted {
			stepData := models.Step{
				ID:         0,
				CategoryID: category.ID,
				Name:       step.Name,
				Difficulty: step.Difficulty,
			}
			if err := transaction.Save(&stepData).Error; err != nil {
				transaction.Rollback()
				return &common.ErrorResponse{
					Code:    500,
					Message: "Failed to create step",
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

	categoryRepo := repositories.CategoryRepository{}
	category, err = categoryRepo.GetCategoryByID(initializers.DB, params.CourseID, category.ID, nil)
	if err != nil {
		return err
	}

	c.JSON(200, CategoryInsertResponse{
		Data: dtos.CategoryDTO{}.From(category),
	})
	return nil
}
