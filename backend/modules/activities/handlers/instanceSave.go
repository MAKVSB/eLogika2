package handlers

import (
	"encoding/json"
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
	"gorm.io/gorm"
)

type ActivityInstanceSaveRequest struct {
	Points  *float64         `json:"points"`
	Content *json.RawMessage `json:"content"`
}

type ActivityInstanceSaveResponse struct {
	InstanceData dtos.ActivityInstanceDTO `json:"instanceData" binding:"required"`
}

// @Summary Saves test instance
// @Tags Activities
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param instanceId path int true "ID of the instance"
// @Param body body ActivityInstanceSaveRequest true "Ability to filter results"
// @Success 200 {object} ActivityInstanceSaveResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/activities/{instanceId} [put]
func ActivityInstanceSave(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID   uint `uri:"courseId" binding:"required"`
			InstanceID uint `uri:"instanceId" binding:"required"`
		},
		ActivityInstanceSaveRequest,
	](c)
	if err != nil {
		return err
	}

	transaction := initializers.DB.Begin()

	// TODO validate from here
	var activityInstance *models.ActivityInstance
	if err := transaction.
		Preload("Result").
		Preload("CourseItem").
		Preload("CourseItem.ActivityDetail").
		InnerJoins("Term").
		First(&activityInstance, params.InstanceID).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    404,
			Message: "Activity instance does not exists",
			Details: err.Error(),
		}
	}

	editable := false

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		transaction.Rollback()
		return err
	}
	if userRole == enums.CourseUserRoleAdmin || userRole == enums.CourseUserRoleGarant || userRole == enums.CourseUserRoleTutor {
		editable = false
		err = elevatedSave(transaction, activityInstance, reqData, userRole, params.CourseID, userData.ID)
		if err != nil {
			transaction.Rollback()
			return err
		}
	} else if userRole == enums.CourseUserRoleStudent {
		editable = true
		err = studentSave(transaction, activityInstance, reqData, userData.ID)
		if err != nil {
			transaction.Rollback()
			return err
		}
	} else {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	if err := transaction.Commit().Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to commit changes",
			Details: err.Error(),
		}
	}

	// Fetch updated data
	if err := initializers.DB.
		InnerJoins("Term").
		InnerJoins("CourseItem").
		InnerJoins("CourseItem.ActivityDetail").
		InnerJoins("Result").
		First(&activityInstance, params.InstanceID).Error; err != nil {
		return &common.ErrorResponse{
			Code:    404,
			Message: "Failed to fetch updated activity instance",
			Details: err.Error(),
		}
	}

	c.JSON(200, ActivityInstanceSaveResponse{
		InstanceData: dtos.ActivityInstanceDTO{}.From(activityInstance, editable),
	})

	return nil
}

func elevatedSave(dbRef *gorm.DB, activityInstance *models.ActivityInstance, reqData *ActivityInstanceSaveRequest, userRole enums.CourseUserRoleEnum, courseId uint, userId uint) *common.ErrorResponse {
	courseItemServ := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
	courseItem, err := courseItemServ.GetCourseItemByID(dbRef, courseId, activityInstance.CourseItemID, userId, userRole, nil, false, nil)
	if err != nil {
		return err
	}
	if !courseItem.Editable {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permission for this item",
		}
	}

	if reqData.Points != nil {
		activityInstance.Result.Points = *reqData.Points
		activityInstance.Result.UpdatedByID = &userId
		if err := dbRef.Save(&activityInstance.Result).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to result for activity",
				Details: err.Error(),
			}
		}
	} else {
		return &common.ErrorResponse{
			Code:    422,
			Message: "Points not specified and is required",
		}
	}
	return nil
}

func studentSave(dbRef *gorm.DB, activityInstance *models.ActivityInstance, reqData *ActivityInstanceSaveRequest, userId uint) *common.ErrorResponse {
	if activityInstance.ParticipantID != userId {
		return &common.ErrorResponse{
			Code:    403,
			Message: "User cannot modify this instance",
		}
	}

	if activityInstance.Term.ActiveTo.Compare(time.Now()) < 0 {
		return &common.ErrorResponse{
			Code:    409,
			Message: "Time expired",
		}
	}

	// User is student
	if reqData.Content != nil {
		activityInstance.Content = *reqData.Content

		// Sync content files
		var files []models.File
		if err := dbRef.Where("id IN ?", utils.GetFilesInsideContent(*reqData.Content)).Find(&files).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to load files",
				Details: err.Error(),
			}
		}

		activityInstance.ContentFiles = files

		if err := dbRef.Model(&activityInstance).Association("ContentFiles").Replace(&activityInstance.ContentFiles); err != nil {

			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to update files",
				Details: err.Error(),
			}
		}

		if err := dbRef.Save(&activityInstance).Error; err != nil {

			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to result for activity",
				Details: err.Error(),
			}
		}
	} else {
		return &common.ErrorResponse{
			Code:    422,
			Message: "Content not specified and is required",
		}
	}
	return nil
}
