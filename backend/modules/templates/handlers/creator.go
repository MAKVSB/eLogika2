package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/templates/dtos"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TemplateCreatorResponse struct {
	Chapters []dtos.TemplateCreatorDTO `json:"chapters"`
}

// @Summary Get chapters, categories, steps and other data required to create or modify test template
// @Tags Templates
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Success 200 {object} QuestionListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/templates/creator [get]
func Creator(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
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
	//

	//

	//

	m := models.Chapter{}

	var chapters []*models.Chapter
	query := initializers.DB.
		Model(&m).
		Where("course_id = ?", params.CourseID).
		Preload("Categories").
		Preload("Categories.Steps")

	if err := query.Find(&chapters).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch questions",
			Details: err.Error(),
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.TemplateCreatorDTO, len(chapters))
	for i, q := range chapters {
		dtoList[i] = dtos.TemplateCreatorDTO{}.From(q)
	}

	c.JSON(200, TemplateCreatorResponse{
		Chapters: dtoList,
	})
	return nil
}
