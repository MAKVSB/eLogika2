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

type TermsListResponse struct {
	Items      []dtos.TermDTO `json:"items"`
	ItemsCount int64          `json:"itemsCount"`
}

type TermsListRequest struct {
	common.SearchRequest
}

// @Summary List all available terms of course item
// @Tags Terms
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TermsListRequest true "Ability to filter results"
// @Success 200 {object} TermsListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId}/terms [get]
func ListByItem(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
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

	termService := services.TermService{}
	terms, termsCount, err := termService.ListTerms(initializers.DB, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, false, searchParams)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.TermDTO, len(terms))
	for i, q := range terms {
		dtoList[i] = dtos.TermDTO{}.From(q)
	}

	c.JSON(200, TermsListResponse{
		Items:      dtoList,
		ItemsCount: termsCount,
	})
	return nil
}
