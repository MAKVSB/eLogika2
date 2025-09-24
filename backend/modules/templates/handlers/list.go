package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/templates/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TemplateListResponse struct {
	Items      []dtos.TemplateListItemDTO `json:"items"`
	ItemsCount int64                      `json:"itemsCount"`
}

type TemplateListRequest struct {
	common.SearchRequest
}

// @Summary List all available templates in course
// @Tags Templates
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TemplateListRequest true "Ability to filter results"
// @Success 200 {object} TemplateListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/templates [get]
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
	templateServ := services.NewTemplateService(repositories.NewTemplateRepository())
	templates, templateCount, err := templateServ.ListTemplates(initializers.DB, params.CourseID, userData.ID, userRole, nil, false, searchParams)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.TemplateListItemDTO, len(templates))
	for i, q := range templates {
		dtoList[i] = dtos.TemplateListItemDTO{}.From(q)
	}

	c.JSON(200, TemplateListResponse{
		Items:      dtoList,
		ItemsCount: templateCount,
	})

	return nil
}
