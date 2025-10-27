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

// @Description Updated list of tutors
type ListTutorResponse struct {
	Tutors []dtos.ClassUserDTO `json:"tutors"`
}

// @Summary List tutors of class
// @Tags Classes
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} ListTutorResponse "Successful operation"
// @Failure 400 {object} common.ErrorResponse "Invalid resource or patch"
// @Failure 403 {object} common.ErrorResponse "Permission or atuhentication errors"
// @Failure 422 {object} common.ErrorResponse "Data validation errors"
// @Failure 500 {object} common.ErrorResponse "Fatal failure"
// @Router /api/v2/courses/ [post] // TODO
func ListTutors(c *gin.Context, userData authdtos.LoggedUserDTO, userRole enums.CourseUserRoleEnum) *common.ErrorResponse {
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
	if err := auth.GetClaimCourseRole(userData, params.CourseID, userRole); err != nil {
		return err
	}

	classService := services.NewClassService(repositories.NewClassRepository())
	class, err := classService.GetClassByID(initializers.DB, params.CourseID, params.ClassID, userData.ID, userRole, nil, false, nil)
	if err != nil {
		return err
	}

	// Convert to DTOs
	dtoList := make([]dtos.ClassUserDTO, len(class.Tutors))
	for i, c := range class.Tutors {
		dtoList[i] = dtos.ClassUserDTO{}.From(c.User)
	}

	c.JSON(200, ListTutorResponse{
		Tutors: dtoList,
	})
	return nil
}
