package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/categories/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type CategoryListResponse struct {
	Items      []dtos.CategoryListItemDTO `json:"items"`
	ItemsCount int64                      `json:"itemsCount"`
}

type CategoryListRequest struct {
	common.SearchRequest
}

// @Summary List all available categories in course
// @Tags Categories
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body CategoryListRequest true "Ability to filter results"
// @Success 200 {object} CategoryListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/categories [get]
func List(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID  uint `uri:"courseId" binding:"required"`
			ChapterId uint `uri:"chapterId"`
		},
		any,
	](c, "search")
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin, or garant
	if userRole != enums.CourseUserRoleAdmin && userRole != enums.CourseUserRoleGarant && userRole != enums.CourseUserRoleTutor {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	categoryRepo := repositories.CategoryRepository{}
	categories, categoryCount, err := categoryRepo.ListCategories(initializers.DB, params.CourseID, nil, false, searchParams)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.CategoryListItemDTO, len(categories))
	for i, q := range categories {
		dtoList[i] = dtos.CategoryListItemDTO{}.From(q)
	}

	c.JSON(200, CategoryListResponse{
		Items:      dtoList,
		ItemsCount: categoryCount,
	})
	return nil
}
