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

type TestInstanceListResponse struct {
	Items      []dtos.TestInstanceListItemDTO `json:"items"`
	ItemsCount int64                          `json:"itemsCount"`
}

type CategoryListRequest struct {
	common.SearchRequest
}

// @Summary List all available test instances of test
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body CategoryListRequest true "Ability to filter results"
// @Success 200 {object} CategoryListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/{testId} [get]
func ListInstance(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID     uint  `uri:"courseId" binding:"required"`
			CourseItemId uint  `uri:"courseItemId" binding:"required"`
			TermId       *uint `uri:"termId"`
			TestID       uint  `uri:"testId" binding:"required"`
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
	testInstances, testInstancesCount, err := testService.ListTestInstances(initializers.DB, params.CourseItemId, params.TermId, params.TestID, userData.ID, userRole, nil, false, searchParams)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.TestInstanceListItemDTO, len(testInstances))
	for i, q := range testInstances {
		dtoList[i] = dtos.TestInstanceListItemDTO{}.From(q)
	}

	c.JSON(200, TestInstanceListResponse{
		Items:      dtoList,
		ItemsCount: testInstancesCount,
	})

	return nil
}
