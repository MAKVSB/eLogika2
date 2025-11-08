package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/print/helpers"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// @Description Request to insert new question
type PrintTestRequest struct {
	CourseItemID         uint  `json:"courseItemId" binding:"required"`
	PrintAnswerSheets    bool  `json:"printAnswerSheets"`
	SeparateAnswerSheets bool  `json:"separateAnswerSheets"`
	TestID               *uint `json:"testId"`
	InstanceID           *uint `json:"instanceId"`
	TermID               *uint `json:"termId"`
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
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	// Check if tutor/garant can view/modify courseItem
	courseItemService := services_course_item.CourseItemService{}
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, reqData.CourseItemID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	var printData []*models.Test
	query := initializers.DB.
		Where("tests.course_item_id = ?", reqData.CourseItemID).
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.
				Unscoped().
				Joins("Question", initializers.DB.Unscoped()).
				Order("\"order\" ASC")
		}).
		Preload("Questions.Answers", func(db *gorm.DB) *gorm.DB {
			return db.
				Unscoped().
				Joins("Answer", initializers.DB.Unscoped()).
				Order("\"order\" ASC")
		}).
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

	filepath, err2 := helpers.PrintTests(printData, courseItem, reqData.PrintAnswerSheets, reqData.SeparateAnswerSheets)
	if err2 != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to generate PDF file",
			Details: err2.Error(),
		}
	}
	c.FileAttachment(filepath, uuid.NewString())

	return nil
}
