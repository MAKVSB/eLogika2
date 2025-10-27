package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
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
	// Check if tutor/garant can view/modify courseItem
	if userRole == enums.CourseUserRoleAdmin {
	} else if userRole == enums.CourseUserRoleGarant {
		var courseItem models.CourseItem
		if err := initializers.DB.
			Where("managed_by = ?", enums.CourseUserRoleGarant).
			Find(&courseItem, params.CourseItemId).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else if userRole == enums.CourseUserRoleTutor {
		var courseItem models.CourseItem
		if err := initializers.DB.
			Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userData.ID).
			Find(&courseItem, params.CourseItemId).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else {
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
