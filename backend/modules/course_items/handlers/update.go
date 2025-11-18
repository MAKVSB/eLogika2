package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_items/dtos"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"elogika.vsb.cz/backend/utils/tiptap"
	"github.com/gin-gonic/gin"
)

type ActivityDetailCourseItemUpdateRequest struct {
	Description    *models.TipTapContent `json:"description" ts_type:"JSONContent"`    // Assignemnt of activity
	ExpectedResult *models.TipTapContent `json:"expectedResult" ts_type:"JSONContent"` // Expected result of the activity
}

type GroupDetailCourseItemUpdateRequest struct {
	Choice    bool `json:"choice"`    // Is students can pick from child items
	ChooseMin uint `json:"chooseMin"` // Student must pick minimum of N items
	ChooseMax uint `json:"chooseMax"` // Student can pick maximum of N items
}

type TestDetailCourseItemUpdateRequest struct {
	TestType        enums.QuestionTypeEnum `json:"testType"`        // Further filters questions used for this test
	TimeLimit       uint                   `json:"timeLimit"`       // Time limit for this test
	ShowResults     bool                   `json:"showResults"`     // Show results to student after finishing test
	ShowTest        bool                   `json:"showTest"`        // Show the exact test (questions)
	ShowCorrectness bool                   `json:"showCorrectness"` // Show what answers were correct
	AllowOffline    bool                   `json:"allowOffline"`    // Allow offline answer sending
	IsPaper         bool                   `json:"isPaper"`         // Test is written physically on paper
	IPRanges        string                 `json:"ipRanges"`        // Allowed ip ranges to write a test
	TestTemplateID  uint                   `json:"testTemplateId"`  // Id of selected test template
}

// @Description Request to update course item
type CourseItemUpdateRequest struct {
	Name              string                      `json:"name"`              // Name
	PointsMin         uint                        `json:"pointsMin"`         // Minimum number of points to pass
	PointsMax         uint                        `json:"pointsMax"`         // Maximum number of points
	Mandatory         bool                        `json:"mandatory"`         // Passing is mandatory
	StudyForm         enums.StudyFormEnum         `json:"studyForm"`         // Study form (allows only students of the same type to join)
	MaxAttempts       uint                        `json:"maxAttempts"`       // Maximum attempts to sign up to this item
	AllowNegative     bool                        `json:"allowNegative"`     // Allow passing negative points outside of this item
	EvaluateByAttempt enums.EvaluateByAttemptEnum `json:"evaluateByAttempt"` // Evaluation mode
	IncludeInResults  bool                        `json:"includeInResults"`  // Include in overall results
	// TODO ManagedBy

	ActivityDetail *ActivityDetailCourseItemUpdateRequest `json:"activityDetail"` // Additional data for type ACTIVITY (homework, project, ...)
	TestDetail     *TestDetailCourseItemUpdateRequest     `json:"testDetail"`     // Additional data for type TEST
	GroupDetail    *GroupDetailCourseItemUpdateRequest    `json:"groupDetail"`    // additional data for type GROUP

	Version uint `json:"version"` // Version signature to prevent concurrency problems
}

// @Description Newly updated course item
type CourseItemUpdateResponse struct {
	Data dtos.CourseItemDTO `json:"data"`
}

// @Summary Update course item
// @Tags CourseItems
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body CourseItemUpdateRequest true "New data for question"
// @Success 200 {object} CourseItemUpdateResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items [put]
func Update(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
		},
		CourseItemUpdateRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin, garant, or tutor
	if userRole != enums.CourseUserRoleGarant && userRole != enums.CourseUserRoleTutor {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	courseItemService := services_course_item.CourseItemService{}
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, true, &reqData.Version)
	if err != nil {
		return err
	}

	if !courseItem.Editable {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	if reqData.TestDetail != nil && reqData.TestDetail.IPRanges != "" {
		if !utils.IsValidIPCondition(reqData.TestDetail.IPRanges) {
			return &common.ErrorResponse{
				Code:    422,
				Message: "Validation failed",
				Details: "IPRange string is not valid",
				FormErrors: common.ErrorObject{
					"testDetail": common.ErrorObject{
						"ipRanges": "IPRange string is not valid string",
					},
				},
			}
		}
	}

	courseItem.Version = courseItem.Version + 1
	courseItem.Name = reqData.Name
	courseItem.PointsMin = reqData.PointsMin
	courseItem.PointsMax = reqData.PointsMax
	courseItem.Mandatory = reqData.Mandatory
	courseItem.MaxAttempts = reqData.MaxAttempts
	courseItem.AllowNegative = reqData.AllowNegative
	courseItem.EvaluateByAttempt = reqData.EvaluateByAttempt
	courseItem.IncludeInResults = reqData.IncludeInResults

	transaction := initializers.DB.Begin()

	switch courseItem.Type {
	case enums.CourseItemTypeActivity:
		courseItem.ActivityDetail.Description = reqData.ActivityDetail.Description
		err = tiptap.FindAndSaveRelations(transaction, userData.ID, reqData.ActivityDetail.Description, &courseItem.ActivityDetail, "DescriptionFiles")
		if err != nil {
			return err
		}
		courseItem.ActivityDetail.ExpectedResult = reqData.ActivityDetail.ExpectedResult
		err = tiptap.FindAndSaveRelations(transaction, userData.ID, reqData.ActivityDetail.ExpectedResult, &courseItem.ActivityDetail, "ExpectedResultFiles")
		if err != nil {
			return err
		}

		if err := transaction.Save(&courseItem.ActivityDetail).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to update inner object",
			}
		}
	case enums.CourseItemTypeGroup:
		courseItem.GroupDetail.Choice = reqData.GroupDetail.Choice
		courseItem.GroupDetail.ChooseMin = reqData.GroupDetail.ChooseMin
		courseItem.GroupDetail.ChooseMax = reqData.GroupDetail.ChooseMax
		if err := transaction.Save(&courseItem.GroupDetail).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to update inner object",
			}
		}
	case enums.CourseItemTypeTest:
		courseItem.TestDetail.TestType = reqData.TestDetail.TestType
		courseItem.TestDetail.TimeLimit = reqData.TestDetail.TimeLimit
		courseItem.TestDetail.ShowResults = reqData.TestDetail.ShowResults
		courseItem.TestDetail.ShowTest = reqData.TestDetail.ShowTest
		courseItem.TestDetail.ShowCorrectness = reqData.TestDetail.ShowCorrectness
		courseItem.TestDetail.AllowOffline = reqData.TestDetail.AllowOffline
		courseItem.TestDetail.IsPaper = reqData.TestDetail.IsPaper
		courseItem.TestDetail.IPRanges = reqData.TestDetail.IPRanges
		courseItem.TestDetail.TestTemplateID = reqData.TestDetail.TestTemplateID
		if err := transaction.Save(&courseItem.TestDetail).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to update inner object",
			}
		}
	}

	if err := transaction.Save(&courseItem).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to update course item",
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

	courseItem, err = courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, CourseItemUpdateResponse{
		Data: dtos.CourseItemDTO{}.From(courseItem),
	})

	return nil
}
