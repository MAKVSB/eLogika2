package handlers

import (
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/activities/dtos"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type ActivityInstanceGetResponse struct {
	InstanceData dtos.ActivityInstanceDTO `json:"instanceData"`
}

// @Summary Starts activity instance for user
// @Tags Activities
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param instanceId path int true "ID of the corresponding test instance"
// @Success 200 {object} ActivityInstanceGetResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/activities/{instanceId} [get]
func ActivityInstanceGet(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID   uint `uri:"courseId" binding:"required"`
			InstanceID uint `uri:"instanceId" binding:"required"`
		},
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}

	var activityInstance *models.ActivityInstance
	if err := initializers.DB.
		InnerJoins("Term").
		InnerJoins("CourseItem").
		InnerJoins("CourseItem.ActivityDetail").
		InnerJoins("Result").
		First(&activityInstance, params.InstanceID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    404,
			Message: "Activity instance does not exists",
			Details: err.Error(),
		}
	}

	editable := true

	// Check if tutor/garant can view/modify courseItem
	if userRole == enums.CourseUserRoleAdmin || userRole == enums.CourseUserRoleGarant || userRole == enums.CourseUserRoleTutor {
		courseItemServ := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
		courseItem, err := courseItemServ.GetCourseItemByID(initializers.DB, params.CourseID, activityInstance.CourseItemID, userData.ID, userRole, nil, false, nil)
		if err != nil {
			return err
		}
		if !courseItem.Editable {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else if userRole == enums.CourseUserRoleStudent {
		if activityInstance.ParticipantID != userData.ID {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permissions",
			}
		}

		if activityInstance.Term.ActiveTo.Compare(time.Now()) < 0 {
			editable = false
		}
	} else {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	c.JSON(200, ActivityInstanceGetResponse{
		InstanceData: dtos.ActivityInstanceDTO{}.From(activityInstance, editable),
	})
	return nil
}
