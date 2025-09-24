package handlers

import (
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/courses/dtos"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

type UserCourseListResponse struct {
	Items []dtos.LoggedUserCourseDTO2 `json:"items"`
}

// @Summary List courses available to user
// @Tags Courses
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param body body CourseListRequest true "Ability to filter results"
// @Success 200 {object} UserCourseListResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/user [get]
func ListUserCourses(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, _, _ := utils.GetRequestData[
		any,
		any,
	](c)
	if err != nil {
		return err
	}

	// TODO validate from here

	courseService := services.NewCourseService(&repositories.CourseRepository{})
	courses, err := courseService.GetUserCourses(userData.ID, &userData.Type)
	if err != nil {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}

	// Convert to DTOs
	dtoList := make([]dtos.LoggedUserCourseDTO2, len(courses))
	for i, c := range courses {
		dtoList[i] = dtos.LoggedUserCourseDTO2{}.From(c)
	}

	c.JSON(200, UserCourseListResponse{
		Items: dtoList,
	})
	return nil
}
