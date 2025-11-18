package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_item_terms/dtos"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TermsListRecursiveResponse struct {
	Items []dtos.TermDTO `json:"items"`
}

type TermsListRecursiveRequest struct {
	common.SearchRequest
}

// @Summary List all available terms of course item recursive
// @Tags Terms
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TermsListRecursiveRequest true "Ability to filter results"
// @Success 200 {object} TermsListRecursiveResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId}/terms/recursive [get]
func ListByItemRecursive(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemId uint `uri:"courseItemId" binding:"required"`
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

	courseItemIDS := []uint{params.CourseItemId}

	var courseItem models.CourseItem
	if err := initializers.DB.Select("parent_id").First(&courseItem, params.CourseItemId).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course item terms",
			Details: err.Error(),
		}
	}

	if courseItem.ParentID != nil {
		courseItemIDS = append(courseItemIDS, *courseItem.ParentID)
	}

	var terms []models.Term
	query := initializers.DB.
		Where("course_item_id IN ?", courseItemIDS)

	if err := query.Find(&terms).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course item terms",
			Details: err.Error(),
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.TermDTO, len(terms))
	for i, q := range terms {
		dtoList[i] = dtos.TermDTO{}.From(&q)
	}

	c.JSON(200, TermsListRecursiveResponse{
		Items: dtoList,
	})

	return nil
}
