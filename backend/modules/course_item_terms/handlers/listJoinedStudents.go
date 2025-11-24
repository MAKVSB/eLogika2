package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/course_item_terms/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	services_course_item "elogika.vsb.cz/backend/services/courseItem"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type ListJoinedStudentsResponse struct {
	Items      []dtos.JoinedStudentDTO `json:"items"`
	ItemsCount int64                   `json:"itemsCount"`
}

// @Summary List all available terms of course item
// @Tags Terms
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param courseId path int true "ID of the corresponding course"
// @Param body body TermsListRequest true "Ability to filter results"
// @Success 200 {object} TermsListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId}/items/{courseItemId}/terms/{termId}/students [get]
func ListJoinedStudents(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _, searchParams := utils.GetRequestDataWithSearch[
		struct {
			CourseID     uint `uri:"courseId" binding:"required"`
			CourseItemID uint `uri:"courseItemId" binding:"required"`
			TermID       uint `uri:"termId" binding:"required"`
		},
		any,
	](c, "search")
	if err != nil {
		return err
	}

	// TODO validate from here

	// Check role validity
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	// Check if tutor/garant can view/modify courseItem
	courseItemService := services_course_item.NewCourseItemService(repositories.NewCourseItemRepository())
	courseItem, err := courseItemService.GetCourseItemByID(initializers.DB, params.CourseID, params.CourseItemID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}
	if !courseItem.Editable {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	skipUsersWithInstance := c.Query("skipUsersWithInstance") != ""
	showHistory := c.Query("showHistory") != ""

	termService := services.TermService{}
	joinedStudents, joinedStudentsCount, err := termService.ListJoinedStudents(
		initializers.DB,
		params.TermID,
		userData.ID,
		userRole,
		nil,
		false,
		searchParams,
		skipUsersWithInstance,
		showHistory,
	)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.JoinedStudentDTO, len(joinedStudents))
	for i, js := range joinedStudents {
		dtoList[i] = dtos.JoinedStudentDTO{}.From(js)
	}

	c.JSON(200, ListJoinedStudentsResponse{
		Items:      dtoList,
		ItemsCount: joinedStudentsCount,
	})
	return nil
}
