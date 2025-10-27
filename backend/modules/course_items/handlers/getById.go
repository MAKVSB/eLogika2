package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_items/dtos"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Newly created courseItem
type CourseItemGetByIdResponse struct {
	Data dtos.CourseItemDTO `json:"data"`
}

// @Summary Get courseItem by id
// @Tags CourseItems
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param courseItemId path int true "ID of the requested item"
// @Success 200 {object} CourseItemGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId} [get]
func GetByID(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
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

	c.JSON(200, CourseItemGetByIdResponse{
		Data: dtos.CourseItemDTO{}.From(courseItem),
	})
	return nil
}
