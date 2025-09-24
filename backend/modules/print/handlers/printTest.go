package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/print/helpers"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// @Description Request to insert new question
type PrintTestRequest struct {
	CourseItemID      uint  `json:"courseItemId" binding:"required"`
	PrintAnswerSheets bool  `json:"printInstances"`
	TestID            *uint `json:"testId"`
	InstanceID        *uint `json:"instanceId"`
	TermID            *uint `json:"termId"`
}

// @Summary Print tests
// @Tags Print
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body PrintTestRequest true "New data for question"
// @Success 200 {file} file "PDF file of tests"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/print/tests [post]
func PrintTest(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		PrintTestRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}
	var courseItem *models.CourseItem
	// Check if tutor/garant can view/modify courseItem
	if userRole == enums.CourseUserRoleAdmin {
		if err := initializers.DB.
			Find(&courseItem, reqData.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else if userRole == enums.CourseUserRoleGarant {
		if err := initializers.DB.
			Where("managed_by = ?", enums.CourseUserRoleGarant).
			Find(&courseItem, reqData.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    403,
				Message: "Not enough permission for this item",
			}
		}
	} else if userRole == enums.CourseUserRoleTutor {
		if err := initializers.DB.
			Where("managed_by = ? AND created_by_id = ?", enums.CourseUserRoleTutor, userData.ID).
			Find(&courseItem, reqData.CourseItemID).Error; err != nil {
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

	var printData []*models.Test
	query := initializers.DB.
		Where("tests.course_item_id = ?", reqData.CourseItemID).
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("\"order\" ASC")
		}).
		Preload("Questions.Question").
		Preload("Questions.Answers").
		Preload("Questions.Answers.Answer").
		Preload("Course")

	if reqData.TermID != nil {
		query = query.InnerJoins("Term", initializers.DB.Where("Term.id = ?", reqData.TermID))
	} else {
		query = query.InnerJoins("Term")
	}

	if reqData.TestID != nil {
		query = query.Where("tests.id = ?", reqData.TestID)
	}

	if reqData.InstanceID != nil {
		query = query.Preload("Instances", func(db *gorm.DB) *gorm.DB {
			return db.Joins("Participant").Where("test_instances.id = ?", reqData.InstanceID)
		})
	} else {
		query = query.Preload("Instances", func(db *gorm.DB) *gorm.DB {
			return db.Joins("Participant")
		})
	}

	if err := query.Find(&printData).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch test",
			Details: err.Error(),
		}
	}

	filepath := helpers.PrintTests(printData, courseItem, reqData.PrintAnswerSheets)
	c.FileAttachment(filepath, uuid.NewString())

	return nil
}
