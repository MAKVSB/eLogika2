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
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TestInstanceGetTelemetryResponse struct {
	Items      []dtos.TestInstanceEventDTO `json:"items"`
	ItemsCount int64                       `json:"itemsCount"`
}

// @Summary Gets events associated with test instance
// @Tags Tests
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param instanceId path int true "ID of the corresponding test instance"
// @Success 200 {object} TestInstanceGetTelemetryResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/{courseItemId}/instance/{instanceId}/telemetry [GET]
func TestInstanceGetTelemetry(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			InstanceID   uint `uri:"instanceId" binding:"required"`
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
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}
	if !courseItem.Editable {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	var testInstanceEvents []models.TestInstanceEvent
	query := initializers.DB.
		Model(&models.TestInstanceEvent{}).
		Preload("User").
		Where("test_instance_id = ?", params.InstanceID)

	query, err = models.TestInstanceEvent{}.ApplyFilters(query, searchParams.ColumnFilters, models.TestInstanceEvent{}, map[string]interface{}{}, "")
	if err != nil {
		return err
	}
	query = models.TestInstanceEvent{}.ApplySorting(query, searchParams.Sorting, "id ASC")
	totalCount := models.TestInstanceEvent{}.GetCount(query) // Gets count before pagination
	query = models.TestInstanceEvent{}.ApplyPagination(query, searchParams.Pagination)

	if err := query.
		Find(&testInstanceEvents).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.TestInstanceEventDTO, len(testInstanceEvents))
	for i, e := range testInstanceEvents {
		dtoList[i] = dtos.TestInstanceEventDTO{}.From(&e)
	}

	c.JSON(200, TestInstanceGetTelemetryResponse{
		Items:      dtoList,
		ItemsCount: totalCount,
	})

	return nil
}
