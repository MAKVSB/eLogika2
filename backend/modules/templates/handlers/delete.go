package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Newly created template
type TemplateDeleteResponse struct {
	Success bool `json:"success"`
}

// @Summary Get template by id
// @Tags Templates
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param templateId path int true "ID of the requested template"
// @Success 200 {object} TemplateDeleteResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/templates/{templateId} [delete]
func TemplateDelete(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
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

	transaction := initializers.DB.Begin()

	templateService := services.TemplateService{}
	template, err := templateService.GetTemplateByID(transaction, params.CourseID, params.TemplateID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		transaction.Rollback()
		return err
	}

	var courseItemsUsingTemplate []*models.CourseItem
	if err := transaction.
		InnerJoins("TestDetail", initializers.DB.Where("TestDetail.test_template_id = ?", template.ID)).
		Find(&courseItemsUsingTemplate).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Failed to delete instance",
		}
	}

	// If used by another resources, cannot delete
	if len(courseItemsUsingTemplate) != 0 {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Template is used by other resources",
			// TODO Resources
		}
	}

	if err := transaction.
		Delete(&template).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Failed to delete template",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Failed to commit changes",
		}
	}

	c.JSON(200, TemplateDeleteResponse{
		Success: true,
	})
	return nil
}
