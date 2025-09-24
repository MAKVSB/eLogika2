package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_items/dtos"
	"elogika.vsb.cz/backend/repositories"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type CourseItemListResponse struct {
	Items      []dtos.CourseItemDTO `json:"items"`
	ItemsCount int64                `json:"itemsCount"`
}

type CourseItemListRequest struct {
	common.SearchRequest
}

// @Summary List all available course items in course
// @Tags CourseItems
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body CourseItemListRequest true "Ability to filter results"
// @Success 200 {object} CourseItemListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items [get]
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
	courseItemServ := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
	courseItems, courseItemCount, err := courseItemServ.ListCourseItems(initializers.DB, params.CourseID, userData.ID, userRole, nil, false, searchParams)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.CourseItemDTO, len(courseItems))
	for i, q := range courseItems {
		dtoList[i] = dtos.CourseItemDTO{}.From(q)
	}

	c.JSON(200, CourseItemListResponse{
		Items:      dtoList,
		ItemsCount: courseItemCount,
	})
	return nil
}
