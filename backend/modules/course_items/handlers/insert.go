package handlers

import (
	"encoding/json"

	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_items/dtos"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type ActivityDetailCourseItemInsertRequest struct {
	Description    json.RawMessage `json:"description" ts_type:"JSONContent"`    // Assignemnt of activity
	ExpectedResult json.RawMessage `json:"expectedResult" ts_type:"JSONContent"` // Expected result of the activity
}

type GroupDetailCourseItemInsertRequest struct {
	Choice    bool `json:"choice"`    // Is students can pick from child items
	ChooseMin uint `json:"chooseMin"` // Student must pick minimum of N items
	ChooseMax uint `json:"chooseMax"` // Student can pick maximum of N items
}

type TestDetailCourseItemInsertRequest struct {
	TestType        enums.QuestionTypeEnum `json:"testType"`                           // Further filters questions used for this test
	TimeLimit       uint                   `json:"timeLimit"`                          // Time limit for this test
	ShowResults     bool                   `json:"showResults"`                        // Show results to student after finishing test
	ShowTest        bool                   `json:"showTest"`                           // Show the exact test (questions)
	ShowCorrectness bool                   `json:"showCorrectness"`                    // Show what answers were correct
	AllowOffline    bool                   `json:"allowOffline"`                       // Allow offline answer sending
	IsPaper         bool                   `json:"isPaper"`                            // Test is written physically on paper
	IPRanges        string                 `json:"ipRanges"`                           // Allowed ip ranges to write a test
	TestTemplateID  uint                   `json:"testTemplateId" validate:"required"` // Id of selected test template
}

// @Description Request to insert new course item
type CourseItemInsertRequest struct {
	Name              string                      `json:"name"`                         // Name
	Type              enums.CourseItemTypeEnum    `json:"type"`                         // Type of item
	PointsMin         uint                        `json:"pointsMin"`                    // Minimum number of points to pass
	PointsMax         uint                        `json:"pointsMax"`                    // Maximum number of points
	Mandatory         bool                        `json:"mandatory"`                    // Passing is mandatory
	StudyForm         enums.StudyFormEnum         `json:"studyForm"`                    // Study form (allows only students of the same type to join)
	MaxAttempts       uint                        `json:"maxAttempts"`                  // Maximum attempts to sign up to this item
	AllowNegative     bool                        `json:"allowNegative"`                // Allow passing negative points outside of this item
	ParentID          *uint                       `json:"parentId" validate:"optional"` // Id of parent group object
	EvaluateByAttempt enums.EvaluateByAttemptEnum `json:"evaluateByAttempt"`            // Evaluation mode
	// TODO ManagedBy

	ActivityDetail *ActivityDetailCourseItemInsertRequest `json:"activityDetail" validate:"optional"` // Additional data for type ACTIVITY (homework, project, ...)
	TestDetail     *TestDetailCourseItemInsertRequest     `json:"testDetail" validate:"optional"`     // Additional data for type TEST
	GroupDetail    *GroupDetailCourseItemInsertRequest    `json:"groupDetail" validate:"optional"`    // additional data for type GROUP
}

// @Description Newly created course item
type CourseItemInsertResponse struct {
	Data dtos.CourseItemDTO `json:"data"`
}

// @Summary Create new course item
// @Tags CourseItems
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body CourseItemInsertRequest true "New data for question"
// @Success 200 {object} CourseItemInsertResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items [post]
func Insert(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, reqData := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
		},
		CourseItemInsertRequest,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin, garant, or tutor
	if userRole != enums.CourseUserRoleGarant && userRole != enums.CourseUserRoleTutor {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	if reqData.Type == enums.CourseItemTypeGroup && reqData.ParentID != nil {
		return &common.ErrorResponse{
			Code:    422,
			Message: "Validation failed",
			Details: "Group cannot have a parent object",
		}
	}

	if reqData.Type == enums.CourseItemTypeTest && reqData.TestDetail.IPRanges != "" {
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

	courseItem := &models.CourseItem{
		Version:           1,
		CourseID:          params.CourseID,
		Name:              reqData.Name,
		Type:              reqData.Type,
		PointsMin:         reqData.PointsMin,
		PointsMax:         reqData.PointsMax,
		Mandatory:         reqData.Mandatory,
		ManagedBy:         userRole,
		CreatedById:       userData.ID,
		StudyForm:         reqData.StudyForm,
		MaxAttempts:       reqData.MaxAttempts,
		AllowNegative:     reqData.AllowNegative,
		EvaluateByAttempt: reqData.EvaluateByAttempt,
	}
	if reqData.ParentID != nil && *reqData.ParentID != uint(0) {
		courseItem.ParentID = reqData.ParentID
	}

	transaction := initializers.DB.Begin()

	switch reqData.Type {
	case enums.CourseItemTypeActivity:
		innerCourseItem := models.CourseItemActivity{
			Description:    reqData.ActivityDetail.Description,
			ExpectedResult: reqData.ActivityDetail.ExpectedResult,
		}
		if err := transaction.Save(&innerCourseItem).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to insert inner object",
			}
		}

		// Sync content files
		var files1 []models.File
		if err := transaction.Where("id IN ?", utils.GetFilesInsideContent(reqData.ActivityDetail.Description)).Find(&files1).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to load files",
				Details: err.Error(),
			}
		}

		var files2 []models.File
		if err := transaction.Where("id IN ?", utils.GetFilesInsideContent(reqData.ActivityDetail.ExpectedResult)).Find(&files2).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to load files",
				Details: err.Error(),
			}
		}

		innerCourseItem.ContentFiles = files1
		innerCourseItem.ContentFiles = append(courseItem.ActivityDetail.ContentFiles, files2...)

		if err := transaction.Model(&innerCourseItem).Association("ContentFiles").Replace(&innerCourseItem.ContentFiles); err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to update files",
				Details: err.Error(),
			}
		}

		courseItem.ActivityDetail = &innerCourseItem
		courseItem.ActivityDetailID = &innerCourseItem.ID
	case enums.CourseItemTypeGroup:
		innerCourseItem := models.CourseItemGroup{
			Choice:    reqData.GroupDetail.Choice,
			ChooseMin: reqData.GroupDetail.ChooseMin,
			ChooseMax: reqData.GroupDetail.ChooseMax,
		}
		if err := transaction.Save(&innerCourseItem).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to insert inner object",
				Details: err.Error(),
			}
		}

		courseItem.GroupDetail = &innerCourseItem
		courseItem.GroupDetailID = &innerCourseItem.ID
	case enums.CourseItemTypeTest:
		innerCourseItem := models.CourseItemTest{
			TestType:        reqData.TestDetail.TestType,
			TimeLimit:       reqData.TestDetail.TimeLimit,
			ShowResults:     reqData.TestDetail.ShowResults,
			ShowTest:        reqData.TestDetail.ShowTest,
			ShowCorrectness: reqData.TestDetail.ShowCorrectness,
			AllowOffline:    reqData.TestDetail.AllowOffline,
			IsPaper:         reqData.TestDetail.IsPaper,
			IPRanges:        reqData.TestDetail.IPRanges,
			TestTemplateID:  reqData.TestDetail.TestTemplateID,
		}
		if err := transaction.Save(&innerCourseItem).Error; err != nil {
			transaction.Rollback()
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to insert inner object",
			}
		}

		courseItem.TestDetail = &innerCourseItem
		courseItem.TestDetailID = &innerCourseItem.ID
	}

	if err := transaction.Save(&courseItem).Error; err != nil {
		transaction.Rollback()
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to insert course item",
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

	courseItemService := services_course_item.CourseItemService{}
	courseItem, err = courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, courseItem.ID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	c.JSON(200, CourseItemInsertResponse{
		Data: dtos.CourseItemDTO{}.From(courseItem),
	})

	return nil
}
