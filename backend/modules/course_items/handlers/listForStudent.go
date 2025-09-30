package handlers

import (
	"fmt"

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
	"gorm.io/gorm"
)

type StudentCourseItemListResponse struct {
	Items       []dtos.StudentCourseItemDTO       `json:"items"`
	TotalPoints float64                           `json:"totalPoints"`
	TotalPassed bool                              `json:"totalPassed"`
	Results     []dtos.StudentCourseItemResultDTO `json:"results"`
}

// @Summary List all available course items in course
// @Tags CourseItems
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body QuestionListRequest true "Ability to filter results"
// @Success 200 {object} QuestionListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/students [get]
func ListForStudent(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
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

	// Check role validity
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}
	// If not admin, garant, or tutor
	if userRole != enums.CourseUserRoleStudent {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	// TODO validate from here

	var results []*models.CourseItemResult
	if err := initializers.DB.
		Where("student_id = ?", userData.ID).
		InnerJoins("Term").
		InnerJoins("CourseItem", initializers.DB.Where("CourseItem.course_id = ?", params.CourseID)).
		Preload("CourseItem.Parent").
		Order("course_item_results.course_item_id, Term.active_from, course_item_results.created_at DESC").
		Find(&results).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch student restuls",
			Details: err.Error(),
		}
	}

	var courseUser *models.CourseUser
	if err := initializers.DB.
		Select("StudyForm").
		Where("course_id = ?", params.CourseID).
		Where("user_id = ?", userData.ID).
		First(&courseUser).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch user data",
			Details: err.Error(),
		}
	}

	var rootCourseItems []models.CourseItem
	if err := initializers.DB.
		Where("course_id = ?", params.CourseID).
		Where("parent_id is NULL").
		Where("study_form = ?", courseUser.StudyForm).
		Joins("ActivityDetail").
		Joins("GroupDetail").
		Joins("TestDetail").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("ActivityDetail").
				Joins("GroupDetail").
				Joins("TestDetail")
		}).
		Find(&rootCourseItems).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course items",
			Details: err.Error(),
		}
	}

	// Convert to DTOs
	totalPoints := float64(0)
	totalPassed := true
	dtoList := make([]dtos.StudentCourseItemDTO, len(rootCourseItems))
	for i, ci := range rootCourseItems {
		innerDto, innerPoints, innerPassed, innerMandatory, _ := CalculateItemResult(&ci, userData.ID, &results, false)
		dtoList[i] = innerDto
		totalPoints += innerPoints

		if innerMandatory {
			if !innerPassed {
				totalPassed = false
			}
		}
	}

	dtoList2 := make([]dtos.StudentCourseItemResultDTO, len(results))
	for i, res := range results {
		dtoList2[i] = dtos.StudentCourseItemResultDTO{}.From(res)
	}

	c.JSON(200, StudentCourseItemListResponse{
		Items:       dtoList,
		Results:     dtoList2,
		TotalPoints: totalPoints,
		TotalPassed: totalPassed,
	})

	return nil
}

// Returns dto, points to count to parent, passed, optional
func CalculateItemResult(ci *models.CourseItem, studentId uint, allResults *[]*models.CourseItemResult, returnResults bool) (dtos.StudentCourseItemDTO, float64, bool, bool, []*dtos.CourseItemResultDTO) {
	dto := dtos.StudentCourseItemDTO{}.From(ci)
	passed := true
	points := float64(0)
	passedOptionalCount := 0
	var resultDtos []*dtos.CourseItemResultDTO

	switch ci.Type {
	case enums.CourseItemTypeActivity, enums.CourseItemTypeTest:

		var selectedResult *models.CourseItemResult
		if returnResults {
			itemResults := services_course_item.FindInResults(allResults, ci.ID, studentId, nil, false)

			dto.Results = make([]*dtos.CourseItemResultDTO, len(itemResults))
			for _, itemResult := range itemResults {
				irDto := dtos.CourseItemResultDTO{}.From(itemResult)
				resultDtos = append(resultDtos, &irDto)
			}

			selectedResults := services_course_item.FindInResults(&itemResults, ci.ID, studentId, nil, true)
			if len(selectedResults) != 0 {
				selectedResult = selectedResults[0]
			}
		} else {
			selectedResults := services_course_item.FindInResults(allResults, ci.ID, studentId, nil, true)
			if len(selectedResults) != 0 {
				selectedResult = selectedResults[0]
			}
		}

		if selectedResult == nil {
			if ci.Mandatory {
				passed = false
			}
		} else {
			points = selectedResult.Points
			if ci.Mandatory {
				if points < float64(ci.PointsMin) {
					passed = false
				}
			}
		}
	case enums.CourseItemTypeGroup:
		for _, ciChildren := range ci.Children {
			innerDto, innerPoints, innerPassed, innerMandatory, innerResults := CalculateItemResult(ciChildren, studentId, allResults, returnResults)

			if returnResults && innerResults != nil {
				resultDtos = append(resultDtos, innerResults...)
			}

			if innerMandatory {
				if !innerPassed {
					passed = false
				}
			} else {
				if innerPassed {
					passedOptionalCount += 1
				}
			}

			points += innerPoints
			dto.Childs = append(dto.Childs, &innerDto)
		}

		if ci.GroupDetail.Choice {
			if passedOptionalCount < int(ci.GroupDetail.ChooseMin) {
				passed = false
			}
		}

		if ci.Mandatory {
			if points < float64(ci.PointsMin) {
				passed = false
			}
		}
	default:
		panic(fmt.Sprintf("unexpected enums.CourseItemTypeEnum: %#v", ci.Type))
	}

	dto.Passed = passed
	dto.Points = points
	dto.Results = resultDtos

	if !ci.AllowNegative {
		points = max(points, 0)
	}
	return dto, points, passed, ci.Mandatory, resultDtos
}
