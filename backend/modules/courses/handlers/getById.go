package handlers

import (
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/modules/courses/dtos"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Course data
type CourseGetByIdResponse struct {
	Data dtos.CourseDTO `json:"data"`
}

// @Summary Get course by id
// @Tags Courses
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the requested course"
// @Success 200 {object} CourseGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId} [get]
func CourseGetByID(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
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

	courseService := services.CourseService{}
	course, err := courseService.GetCourseByID(initializers.DB, params.CourseID, userData.ID, userRole, nil, true, nil)
	if err != nil {
		return err
	}

	// TODO validate from here

	c.JSON(200, CourseGetByIdResponse{
		Data: dtos.CourseDTO{}.From(course),
	})

	return nil
}
