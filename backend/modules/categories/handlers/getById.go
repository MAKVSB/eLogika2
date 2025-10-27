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

// @Description Category data
type CategoryGetByIdResponse struct {
	Data dtos.CategoryDTO `json:"data"`
}

// @Summary Get category by id
// @Tags Categories
// @Security ApiKeyAuth
// @Produce  json
// @Param categoryId path int true "ID of the requested category"
// @Success 200 {object} CategoryGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/categories/{categoryId} [get]
func GetByID(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID   uint `uri:"courseId" binding:"required"`
			CategoryID uint `uri:"categoryId" binding:"required"`
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

	categoryRepo := repositories.CategoryRepository{}
	category, err := categoryRepo.GetCategoryByID(initializers.DB, params.CourseID, params.CategoryID, nil)
	if err != nil {
		return err
	}

	c.JSON(200, CategoryGetByIdResponse{
		Data: dtos.CategoryDTO{}.From(category),
	})
	return nil
}
