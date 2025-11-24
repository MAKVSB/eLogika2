package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TestListResponse struct {
	Items      []dtos.TestListItemDTO `json:"items"`
	ItemsCount int64                  `json:"itemsCount"`
}

type TestListRequest struct {
	common.SearchRequest
}

// @Summary List all available tests of specific course item
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TestListRequest true "Ability to filter results"
// @Success 200 {object} TestListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests [get]
func List(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID     uint  `uri:"courseId" binding:"required"`
			CourseItemId uint  `uri:"courseItemId" binding:"required"`
			TermId       *uint `uri:"termId"`
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

	courseItemService := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, params.CourseItemId, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}
	if !courseItem.Editable {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	// Load data
	testService := services.NewTestService(repositories.NewTestRepository())
	tests, testsCount, err := testService.ListTests(initializers.DB, params.CourseID, params.CourseItemId, params.TermId, userData.ID, userRole, nil, false, searchParams)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.TestListItemDTO, len(tests))
	for i, q := range tests {
		dtoList[i] = dtos.TestListItemDTO{}.From(q)
	}

	c.JSON(200, TestListResponse{
		Items:      dtoList,
		ItemsCount: testsCount,
	})

	return nil
}
