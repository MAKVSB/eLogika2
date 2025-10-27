package handlers

import (
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/activities/dtos"
	"elogika.vsb.cz/backend/modules/activities/helpers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type ListAvailableActivitiesResponse struct {
	Instances []dtos.ActivityInstanceStudentListItemDTO `json:"instances"`
	Items     []dtos.StudentActivityDTO                 `json:"items"`
}

// @Summary List activities that user is allowed to submit
// @Tags Activities
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Success 200 {object} ListAvailableActivitiesResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/activities/available [get]
func ListAvailable(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
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

	if userRole != enums.CourseUserRoleStudent {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	var terms []models.Term
	if err := initializers.DB.
		Model(&models.Term{}).
		Where("terms.course_id = ?", params.CourseID).
		InnerJoins("CourseItem").
		Joins("JOIN user_terms ON user_terms.term_id = terms.id AND user_terms.user_id = ? AND user_terms.deleted_at is NULL", userData.ID).
		Where("active_from < ?", time.Now()).
		Where("active_to > ?", time.Now()).
		Find(&terms).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch terms",
			Details: err.Error(),
		}
	}

	for _, t := range terms {
		t.CourseItem.RecursiveLoadChildren(nil)
	}

	activeTests := make([]dtos.StudentActivityDTO, 0)

	for _, t := range terms {
		if t.CourseItem.Type == enums.CourseItemTypeGroup {
			for _, ci := range t.CourseItem.Children {
				if ci.Type == enums.CourseItemTypeActivity {

					triesLeft, err := helpers.GetActivityAttemptsLeft(ci.ID, t.ID, userData.ID)
					if err != nil {
						return err
					}

					activeTests = append(activeTests, dtos.StudentActivityDTO{
						TermID:         t.ID,
						TermName:       t.Name,
						CourseItemId:   ci.ID,
						CourseItemName: ci.Name,
						TriesLeft:      triesLeft,
						ActiveFrom:     t.ActiveFrom,
						ActiveTo:       t.ActiveTo,
						CanStart:       triesLeft != 0,
					})
				}
			}
		} else if t.CourseItem.Type == enums.CourseItemTypeActivity {
			triesLeft, err := helpers.GetActivityAttemptsLeft(t.CourseItemID, t.ID, userData.ID)
			if err != nil {
				return err
			}

			activeTests = append(activeTests, dtos.StudentActivityDTO{
				TermID:         t.ID,
				TermName:       t.Name,
				CourseItemId:   t.CourseItemID,
				CourseItemName: t.CourseItem.Name,
				TriesLeft:      triesLeft,
				ActiveFrom:     t.ActiveFrom,
				ActiveTo:       t.ActiveTo,
				CanStart:       triesLeft != 0,
			})
		}
	}

	var activeInstances []models.ActivityInstance
	if err := initializers.DB.
		Model(&models.ActivityInstance{}).
		InnerJoins("CourseItem").
		InnerJoins("Term").
		Where("CourseItem.course_id = ?", params.CourseID).
		Where("Term.active_from < ? and Term.active_to > ?", time.Now(), time.Now()).
		Find(&activeInstances).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch terms",
			Details: err.Error(),
		}
	}

	activeInstancesDtos := make([]dtos.ActivityInstanceStudentListItemDTO, 0)
	for _, ai := range activeInstances {
		activeInstancesDtos = append(activeInstancesDtos, dtos.ActivityInstanceStudentListItemDTO{}.From(&ai))
	}

	c.JSON(200, ListAvailableActivitiesResponse{
		Instances: activeInstancesDtos,
		Items:     activeTests,
	})
	return nil
}
