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
	"github.com/gin-gonic/gin"
)

// @Description Newly created courseItem
type CourseItemListResultsResponse struct {
	Data []dtos.CourseItemResultsDTO `json:"data"`
}

// @Summary Get courseItem by id
// @Tags CourseItems
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param courseItemId path int true "ID of the requested item"
// @Success 200 {object} CourseItemListResultsResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId}/results [get]
func ListResults(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
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

	courseItemService := services_course_item.CourseItemService{}
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	var courseItemIDs []uint
	courseItemIDs = append(courseItemIDs, courseItem.ID)
	for _, courseItemChildren := range courseItem.Children {
		courseItemIDs = append(courseItemIDs, courseItemChildren.ID)
	}

	var participatingUsers []uint
	if err := initializers.DB.
		Model(&models.CourseItemResult{}).
		Select("Distinct student_id").
		Where("course_item_id in ?", courseItemIDs).
		Pluck("student_id", &participatingUsers).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course items",
			Details: err.Error(),
		}
	}

	var users []models.User
	if err := initializers.DB.
		Select("users.id, username, first_name, family_name, email").
		Where("users.id in ?", participatingUsers).
		Order("family_name").
		Find(&users).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch course items",
			Details: err.Error(),
		}
	}

	var results []*models.CourseItemResult
	if err := initializers.DB.
		InnerJoins("Term").
		Preload("CourseItem").
		Preload("CourseItem.Parent").
		Preload("TestInstance").
		Preload("ActivityInstance").
		Order("course_item_results.course_item_id, Term.active_from, course_item_results.created_at DESC").
		Where("course_item_results.course_item_id in ?", courseItemIDs).
		Find(&results).Error; err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch student restuls",
			Details: err.Error(),
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.CourseItemResultsDTO, len(users))
	for i, u := range users {
		dtoList[i] = dtos.CourseItemResultsDTO{}.From(&u)

		_, innerPoints, innerPassed, _, innerResults := CalculateItemResult(courseItem, u.ID, &results, true, courseItem.Type != enums.CourseItemTypeGroup)
		dtoList[i].Points = innerPoints
		dtoList[i].Passed = innerPassed
		dtoList[i].Results = innerResults
	}

	c.JSON(200, CourseItemListResultsResponse{
		Data: dtoList,
	})

	return nil
}
