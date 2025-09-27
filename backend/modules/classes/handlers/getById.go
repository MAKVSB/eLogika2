package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/classes/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
	"elogika.vsb.cz/backend/services"
	"elogika.vsb.cz/backend/utils"
	"github.com/gin-gonic/gin"
)

// @Description Course data
type ClassGetByIdResponse struct {
	Data dtos.ClassDTO `json:"data"`
}

// @Summary Get class by id
// @Tags Classes
// @Security ApiKeyAuth
// @Produce  json
// @Param courseId path int true "ID of the requested course"
// @Success 200 {object} ClassGetByIdResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/{courseId} [get]
func ClassGetByID(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
	// Load request data
	err, params, _ := utils.GetRequestData[
		struct {
			CourseID uint `uri:"courseId" binding:"required"`
			ClassID  uint `uri:"classId" binding:"required"`
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

	classService := services.NewClassService(repositories.NewClassRepository())
	class, err := classService.GetClassByID(initializers.DB, params.CourseID, params.ClassID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}

	c.JSON(200, ClassGetByIdResponse{
		Data: dtos.ClassDTO{}.From(class),
	})
	return nil
}
