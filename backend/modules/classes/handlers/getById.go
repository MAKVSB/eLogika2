package handlers

import (
	"elogika.vsb.cz/backend/auth"
	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/models"
	authdtos "elogika.vsb.cz/backend/modules/auth/dtos"
	"elogika.vsb.cz/backend/modules/classes/dtos"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/modules/common/enums"
	"elogika.vsb.cz/backend/repositories"
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

	classRepo := repositories.NewClassRepository()
	var class *models.Class
	// If not admin, garant or tutor
	if userRole == enums.CourseUserRoleAdmin {
		class, err = classRepo.GetClassByIDAdmin(initializers.DB, params.CourseID, params.ClassID, userData.ID, true, nil)
	} else if userRole == enums.CourseUserRoleGarant {
		class, err = classRepo.GetClassByIDGarant(initializers.DB, params.CourseID, params.ClassID, userData.ID, true, nil)
	} else if userRole == enums.CourseUserRoleTutor {
		class, err = classRepo.GetClassByIDTutor(initializers.DB, params.CourseID, params.ClassID, userData.ID, true, nil)
	} else {
		return &common.ErrorResponse{
			Code:    403,
			Message: "Not enough permissions",
		}
	}
	if err != nil {
		return err
	}

	c.JSON(200, ClassGetByIdResponse{
		Data: dtos.ClassDTO{}.From(class),
	})
	return nil
}
