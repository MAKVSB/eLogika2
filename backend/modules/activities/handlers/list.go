package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/modules/activities/dtos"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type ActivityListResponse struct {
	Items      []dtos.ActivityListItemDTO `json:"items"`
	ItemsCount int64                      `json:"itemsCount"`
}

type ActivityListRequest struct {
	common.SearchRequest
}

// @Summary List all available activities of specific course item
// @Tags Activities
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body ActivityListRequest true "Ability to filter results"
// @Success 200 {object} ActivityListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/activities [get]
func List(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemId uint `uri:"courseItemId" binding:"required"`
			TermId       uint `uri:"termId"`
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
	activityService := services.NewActivityService(repositories.NewActivityRepository())
	activities, activitiesCount, err := activityService.ListActivityInstances(initializers.DB, params.CourseItemId, &params.TermId, userData.ID, userRole, nil, false, searchParams)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.ActivityListItemDTO, len(activities))
	for i, q := range activities {
		dtoList[i] = dtos.ActivityListItemDTO{}.From(q)
	}

	c.JSON(200, ActivityListResponse{
		Items:      dtoList,
		ItemsCount: activitiesCount,
	})
	return nil
}
