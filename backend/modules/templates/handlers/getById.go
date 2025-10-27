package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/templates/dtos"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Newly created template
type TemplateGetByIdResponse struct {
	Data dtos.TemplateDTO `json:"data"`
}

// @Summary Get template by id
// @Tags Templates
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param templateId path int true "ID of the requested template"
// @Success 200 {object} TemplateGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/templates/{templateId} [get]
func TemplateGetByID(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID   uint `uri:"courseId" binding:"required"`
			TemplateID uint `uri:"templateId" binding:"required"`
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

	templateService := services.TemplateService{}
	template, err := templateService.GetTemplateByID(initializers.DB, params.CourseID, params.TemplateID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, TemplateGetByIdResponse{
		Data: dtos.TemplateDTO{}.From(template),
	})
	return nil
}
