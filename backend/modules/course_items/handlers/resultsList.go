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
// @Router /api/v2/courses/{courseId}/items/{courseItemId} [get]
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
	if err := auth.GetClaimCourseRole(userData.Courses, params.CourseID, userRole); err != nil {
		return err
	}

	courseItemService := services_course_item.CourseItemService{}
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	// Here data
	switch courseItem.Type {
	case enums.CourseItemTypeActivity, enums.CourseItemTypeTest:
		var participatingUsers []uint
		if err := initializers.DB.
			Model(&models.CourseItemResult{}).
			Select("Distinct student_id").
			Where("course_item_id = ?", courseItem.ID).
			Pluck("student_id", &participatingUsers).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch course items",
				Details: err.Error(),
			}
		}

		var users []models.User
		if err := initializers.DB.
			Select("id, username, first_name, family_name, email").
			Where("id in ?", participatingUsers).
			Preload("Results", func(db *gorm.DB) *gorm.DB {
				return db.Where("course_item_id", courseItem.ID).Order("selected DESC")
			}).
			Find(&users).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch course items",
				Details: err.Error(),
			}
		}

		// Convert to DTOs
		dtoList := make([]dtos.CourseItemResultsDTO, len(users))
		for i, u := range users {
			dtoList[i] = dtos.CourseItemResultsDTO{}.From(&u)
			if len(u.Results) != 0 {
				dtoList[i].Points = u.Results[0].Points
				dtoList[i].Final = u.Results[0].Final
				if u.Results[0].Points > float64(courseItem.PointsMin) || !courseItem.Mandatory {
					dtoList[i].Passed = true
				}
			}
			if !courseItem.Mandatory {
				dtoList[i].Passed = true
			}
		}

		c.JSON(200, CourseItemListResultsResponse{
			Data: dtoList,
		})
	case enums.CourseItemTypeGroup:
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

		var rootCourseItem []models.CourseItem
		if err := initializers.DB.
			Where("course_id = ?", params.CourseID).
			Joins("ActivityDetail").
			Joins("GroupDetail").
			Joins("TestDetail").
			Preload("Result", func(db *gorm.DB) *gorm.DB {
				return db.Where("student_id = ?", userData.ID).Where("selected = ?", true)
			}).
			Preload("Results", func(db *gorm.DB) *gorm.DB {
				return db.Where("student_id in ?", courseItemIDs)
			}).
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.
					Joins("ActivityDetail").
					Joins("GroupDetail").
					Joins("TestDetail").
					Preload("Result", func(db *gorm.DB) *gorm.DB {
						return db.Where("student_id = ?", userData.ID).Where("selected = ?", true)
					})
			}).
			Find(&rootCourseItem, params.CourseItemID).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch course items",
				Details: err.Error(),
			}
		}

		var users []models.User
		if err := initializers.DB.
			Select("id, username, first_name, family_name, email").
			Where("id in ?", participatingUsers).
			Preload("Results", func(db *gorm.DB) *gorm.DB {
				return db.Where("course_item_id in ?", courseItemIDs).Order("course_item_id, selected DESC")
			}).
			Find(&users).Error; err != nil {
			return &common.ErrorResponse{
				Code:    500,
				Message: "Failed to fetch course items",
				Details: err.Error(),
			}
		}

		innerDto, innerPoints, innerPassed, innerMandatory := CalculateItemResult(courseItem)

		fmt.Println(innerDto)
		fmt.Println(innerPoints)
		fmt.Println(innerPassed)
		fmt.Println(innerMandatory)

		// Convert to DTOs
		dtoList := make([]dtos.CourseItemResultsDTO, len(users))
		for i, u := range users {
			dtoList[i] = dtos.CourseItemResultsDTO{}.From(&u)
			if len(u.Results) != 0 {
				dtoList[i].Points = u.Results[0].Points
				dtoList[i].Final = u.Results[0].Final
				if u.Results[0].Points > float64(courseItem.PointsMin) || !courseItem.Mandatory {
					dtoList[i].Passed = true
				}
			}
			if !courseItem.Mandatory {
				dtoList[i].Passed = true
			}
		}

		c.JSON(200, CourseItemListResultsResponse{
			Data: dtoList,
		})

	default:
		panic(fmt.Sprintf("unexpected enums.CourseItemTypeEnum: %#v", courseItem.Type))
	}

	return nil
}
