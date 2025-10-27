package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_item_terms/dtos"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Get Term by id
type TermsGetByIdResponse struct {
	Data dtos.TermDTO `json:"data"`
}

// @Summary Get term by id
// @Tags Terms
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param courseItemId path int true "ID of the corresponding course item"
// @Param termId path int true "ID of the requested item"
// @Success 200 {object} TermsGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId}/terms/{termId} [get]
func GetByID(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			TermID       uint `uri:"termId" binding:"required"`
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

	termService := services.TermService{}
	term, err := termService.GetTermByID(initializers.DB, params.CourseID, params.CourseItemID, params.TermID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, TermsGetByIdResponse{
		Data: dtos.TermDTO{}.From(term),
	})
	return nil
}
