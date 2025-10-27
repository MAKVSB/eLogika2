package handlers

import (
	"time"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/tests/dtos"
	"elogika.vsb.cz/backend/modules/tests/helpers"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type ListAvailableTestsResponse struct {
	Instances []dtos.TestInstanceStudentListItemDTO `json:"instances"`
	Items     []dtos.StudentTestDTO                 `json:"items"`
}

// @Summary List tests that user is allowed to write
// @Tags Tests
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Success 200 {object} ListAvailableTestsResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/tests/available [get]
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
		Where("course_id = ?", params.CourseID).
		Preload("CourseItem").
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

	activeTests := make([]dtos.StudentTestDTO, 0)

	for _, t := range terms {
		if t.CourseItem.Type == enums.CourseItemTypeGroup {
			for _, ci := range t.CourseItem.Children {
				if ci.Type == enums.CourseItemTypeTest {

					triesLeft, err := helpers.GetTestAttemptsLeft(ci.ID, t.ID, userData.ID)
					if err != nil {
						return err
					}

					activeTests = append(activeTests, dtos.StudentTestDTO{
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
		} else if t.CourseItem.Type == enums.CourseItemTypeTest {
			triesLeft, err := helpers.GetTestAttemptsLeft(t.CourseItemID, t.ID, userData.ID)
			if err != nil {
				return err
			}

			activeTests = append(activeTests, dtos.StudentTestDTO{
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

	var activeInstances []models.TestInstance
	if err := initializers.DB.
		Model(&models.TestInstance{}).
		InnerJoins("CourseItem").
		InnerJoins("Term").
		Where("CourseItem.course_id = ?", params.CourseID).
		Where("form = ?", enums.TestInstanceFormOnline).
		Where("(state = ? and active_from < ? and Term.active_to > ?) or (state = ? and ends_at> ?)", enums.TestInstanceStateReady, time.Now(), time.Now(), enums.TestInstanceStateActive, time.Now()).
		Find(&activeInstances).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch terms",
			Details: err.Error(),
		}
	}

	activeInstancesDtos := make([]dtos.TestInstanceStudentListItemDTO, 0)
	for _, ai := range activeInstances {
		activeInstancesDtos = append(activeInstancesDtos, dtos.TestInstanceStudentListItemDTO{}.From(&ai))
	}

	c.JSON(200, ListAvailableTestsResponse{
		Instances: activeInstancesDtos,
		Items:     activeTests,
	})
	return nil
}
