package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type TestInstancePrepareRequest struct {
	TermID       uint `json:"termId" uri:"termId" binding:"required"`             // TODO check that user has permission for this term
	CourseItemID uint `json:"courseItemId" uri:"courseItemId" binding:"required"` // TODO check that user has permission for this term
}

type TestInstancePrepareResponse struct {
	InstanceID uint `json:"instanceId"`
}

// @Summary Starts activity instance for user
// @Tags Activities
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TestGeneratorRequest true "Ability to filter results"
// @Success 200 {object} TestInstancePrepareResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/activities/start [get]
func ActivityInstancePrepare(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		TestInstancePrepareRequest,
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

	transaction := initializers.DB.Begin()

	// Check if student has joined this term
	var term models.Term
	if err := transaction.
		Joins("JOIN user_terms ON user_terms.term_id = terms.id AND user_terms.user_id = ?", userData.ID).
		First(&term, reqData.TermID).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch term",
			Details: err.Error(),
		}
	}

	courseItemQuery := `
		WITH course_tree AS (
			-- start with the course_item linked to the term
			SELECT ci.*
			FROM course_items ci
			INNER JOIN terms t ON ci.id = t.course_item_id
			WHERE t.id = ?  -- replace with the term id
			UNION ALL
			-- recursively select children
			SELECT ci.*
			FROM course_items ci
			INNER JOIN course_tree ct ON ci.parent_id = ct.id
		)
		SELECT *
		FROM course_tree WHERE id = ?;
	`

	var courseItem models.CourseItem
	if err := transaction.
		Raw(courseItemQuery, reqData.TermID, reqData.CourseItemID).
		Preload("TestDetail").
		First(&courseItem).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course item",
			Details: err.Error(),
		}
	}

	var activityInstance *models.ActivityInstance
	if err := transaction.
		Where("term_id = ?", reqData.TermID).
		Where("course_item_id = ?", reqData.CourseItemID).
		Find(&activityInstance).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to check for running instances",
			Details: err.Error(),
		}
	}

	if activityInstance == nil || activityInstance.ID == 0 {
		activityInstance = &models.ActivityInstance{
			ParticipantID: userData.ID,
			TermID:        term.ID,
			CourseItemID:  courseItem.ID,
		}

		// Add result object
		activityInstance.Result = &models.CourseItemResult{
			Version:      0,
			CourseItemID: activityInstance.CourseItemID,
			TermID:       activityInstance.TermID,
			StudentID:    activityInstance.ParticipantID,
		}

		if err := transaction.Save(&activityInstance).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to create test instance",
				Details: err.Error(),
			}
		}

		if err := transaction.Commit().Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to commit changes",
			}
		}
	}

	c.JSON(200, TestInstancePrepareResponse{
		InstanceID: activityInstance.ID,
	})

	return nil
}
